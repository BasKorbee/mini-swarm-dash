package main

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

var (
	cpuCache   = map[string]cpuSample{}
	cpuCacheMu sync.Mutex
)

// calcCPUPercent computes CPU usage % by diffing the current sample against
// the previous one stored in cpuCache. Returns 0 on the first call for a
// given container (no prior sample to diff against).
func calcCPUPercent(id string, s *container.StatsResponse) float64 {
	cur := cpuSample{
		totalUsage:  s.CPUStats.CPUUsage.TotalUsage,
		systemUsage: s.CPUStats.SystemUsage,
		onlineCPUs:  s.CPUStats.OnlineCPUs,
	}
	cpuCacheMu.Lock()
	prev, ok := cpuCache[id]
	cpuCache[id] = cur
	cpuCacheMu.Unlock()

	if !ok {
		return 0
	}
	cpuDelta := float64(cur.totalUsage) - float64(prev.totalUsage)
	sysDelta := float64(cur.systemUsage) - float64(prev.systemUsage)
	if sysDelta == 0 {
		return 0
	}
	return (cpuDelta / sysDelta) * float64(cur.onlineCPUs) * 100
}

// calcMemUsage returns actual memory consumption in bytes.
// We subtract the page cache ("file" in cgroups v2) from raw usage because
// the kernel can evict it at any time — it isn't truly "used" by the process.
func calcMemUsage(s *container.StatsResponse) int64 {
	cache := s.MemoryStats.Stats["file"]
	usage := int64(s.MemoryStats.Usage) - int64(cache)
	if usage < 0 {
		return 0
	}
	return usage
}

// getLocalStats collects CPU and memory stats for every running swarm container
// on this node by calling the local Docker socket. It is called by /api/local-stats
// on each node instance, and by the manager when it fans out to peer nodes.
func getLocalStats(cli *client.Client) (*NodeStats, error) {
	ctx := context.Background()

	info, err := cli.Info(ctx, client.InfoOptions{})
	if err != nil {
		return nil, err
	}

	// Only collect stats for swarm-managed containers (have the swarm task label)
	f := make(client.Filters).
		Add("status", "running").
		Add("label", "com.docker.swarm.task")

	result, err := cli.ContainerList(ctx, client.ContainerListOptions{Filters: f})
	if err != nil {
		return nil, err
	}
	logger.Debug("listed local swarm containers", "count", len(result.Items))

	stats := &NodeStats{
		NodeID:      info.Info.Swarm.NodeID,
		Hostname:    info.Info.Name,
		TempCelsius: readNodeTemp(),
	}

	type indexedContainer struct {
		i  int
		cs ContainerStats
		ok bool
	}
	ch := make(chan indexedContainer, len(result.Items))

	for i, c := range result.Items {
		go func(i int, c container.Summary) {
			logger.Debug("fetching container stats", "id", c.ID[:12], "service", c.Labels["com.docker.swarm.service.name"])
			resp, err := cli.ContainerStats(ctx, c.ID, client.ContainerStatsOptions{})
			if err != nil {
				logger.Debug("container stats error", "id", c.ID[:12], "err", err)
				ch <- indexedContainer{i: i}
				return
			}
			var s container.StatsResponse
			err = json.NewDecoder(resp.Body).Decode(&s)
			resp.Body.Close()
			if err != nil {
				logger.Debug("container stats decode error", "id", c.ID[:12], "err", err)
				ch <- indexedContainer{i: i}
				return
			}

			name := ""
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}
			cs := ContainerStats{
				ID:          c.ID,
				Name:        name,
				ServiceName: c.Labels["com.docker.swarm.service.name"],
				CPUPercent:  calcCPUPercent(c.ID, &s),
				MemUsage:    calcMemUsage(&s),
				MemLimit:    int64(s.MemoryStats.Limit),
			}
			logger.Debug("container stats collected", "name", cs.Name, "cpu_pct", cs.CPUPercent, "mem_usage", cs.MemUsage, "mem_limit", cs.MemLimit)
			ch <- indexedContainer{i: i, cs: cs, ok: true}
		}(i, c)
	}

	ordered := make([]ContainerStats, len(result.Items))
	for range result.Items {
		r := <-ch
		if r.ok {
			ordered[r.i] = r.cs
		}
	}
	for _, cs := range ordered {
		if cs.ID != "" {
			stats.Containers = append(stats.Containers, cs)
		}
	}
	return stats, nil
}
