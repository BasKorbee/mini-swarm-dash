package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"
)

// cpuSample holds the last raw CPU counters for a container so we can compute
// a delta on the next poll without needing IncludePreviousSample (which blocks ~1s).
type cpuSample struct {
	totalUsage  uint64
	systemUsage uint64
	onlineCPUs  uint32
}

var (
	cpuCache   = map[string]cpuSample{}
	cpuCacheMu sync.Mutex
)

// readNodeTemp reads the average temperature (°C) across all thermal zones
// from /sys/class/thermal/thermal_zone*/temp. Returns nil if unavailable.
func readNodeTemp() *float64 {
	zones, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil || len(zones) == 0 {
		return nil
	}
	var sum float64
	var count int
	for _, path := range zones {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		millideg, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			continue
		}
		sum += float64(millideg) / 1000.0
		count++
	}
	if count == 0 {
		return nil
	}
	avg := math.Round(sum/float64(count)*10) / 10
	return &avg
}

// NodeStats is what each node-local instance exposes at /api/local-stats
type NodeStats struct {
	NodeID      string           `json:"node_id"`
	Hostname    string           `json:"hostname"`
	TempCelsius *float64         `json:"temp_celsius"`
	Containers  []ContainerStats `json:"containers"`
}

// ContainerStats holds the resource usage snapshot for a single container.
type ContainerStats struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ServiceName string  `json:"service_name"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemUsage    int64   `json:"mem_usage"`
	MemLimit    int64   `json:"mem_limit"`
}

// NodeDetails holds the subset of swarm.Node fields used by the frontend.
type NodeDetails struct {
	ID          string  `json:"id"`
	Hostname    string  `json:"hostname"`
	State       string  `json:"state"`
	Addr        string  `json:"addr"`
	Role        string  `json:"role"`
	IsLeader    bool    `json:"is_leader"`
	NanoCPUs    int64   `json:"nano_cpus"`
	MemoryBytes int64   `json:"memory_bytes"`
}

// ServiceSummary holds the subset of swarm.Service fields used by the frontend.
type ServiceSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	IsGlobal bool   `json:"is_global"`
	Replicas *uint64 `json:"replicas"`
}

// DashboardNode is the enriched node sent to the frontend
type DashboardNode struct {
	Details     NodeDetails      `json:"details"`
	TempCelsius *float64         `json:"temp_celsius"`
	Containers  []ContainerStats `json:"containers"`
}

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
			resp, err := cli.ContainerStats(ctx, c.ID, client.ContainerStatsOptions{})
			if err != nil {
				ch <- indexedContainer{i: i}
				return
			}
			var s container.StatsResponse
			err = json.NewDecoder(resp.Body).Decode(&s)
			resp.Body.Close()
			if err != nil {
				ch <- indexedContainer{i: i}
				return
			}

			name := ""
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}
			ch <- indexedContainer{
				i: i,
				cs: ContainerStats{
					ID:          c.ID,
					Name:        name,
					ServiceName: c.Labels["com.docker.swarm.service.name"],
					CPUPercent:  calcCPUPercent(c.ID, &s),
					MemUsage:    calcMemUsage(&s),
					MemLimit:    int64(s.MemoryStats.Limit),
				},
				ok: true,
			}
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

// gzipFileServer wraps http.FileServer to serve pre-compressed .gz files when
// the client accepts gzip and a .gz sibling exists on disk.
func gzipFileServer(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fs.ServeHTTP(w, r)
			return
		}
		p := filepath.Join(dir, filepath.Clean("/"+r.URL.Path))
		if _, err := os.Stat(p + ".gz"); err == nil {
			r2 := r.Clone(r.Context())
			r2.URL.Path += ".gz"
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			// Strip Content-Type sniffing from the .gz extension
			ext := filepath.Ext(p)
			if ct := mime.TypeByExtension(ext); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
			fs.ServeHTTP(w, r2)
			return
		}
		fs.ServeHTTP(w, r)
	})
}

func main() {
	// For easy local client development use fake data
	if os.Getenv("MOCK") == "true" {
		log.Println("Mock mode enabled — using fake swarm data")
		registerMockHandlers()
		http.Handle("/", gzipFileServer("client/dist"))
		log.Println("Swarm dashboard running on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
		return
	}

	// Connect to the local Docker socket (/var/run/docker.sock must be mounted)
	cli, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal("Error creating Docker client:", err)
	}

	// Returns all services.
	http.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		services, err := cli.ServiceList(context.Background(), client.ServiceListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := make([]ServiceSummary, len(services.Items))
		for i, s := range services.Items {
			result[i] = ServiceSummary{
				ID:       s.ID,
				Name:     s.Spec.Name,
				Image:    s.Spec.TaskTemplate.ContainerSpec.Image,
				IsGlobal: s.Spec.Mode.Global != nil,
				Replicas: func() *uint64 {
					if s.Spec.Mode.Replicated != nil {
						return s.Spec.Mode.Replicated.Replicas
					}
					return nil
				}(),
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Fans out to each node's /api/local-stats via its swarm IP, then returns
	// node metadata (capacity, status) enriched with live container stats.
	// we reach peers directly via node.Status.Addr from the swarm node list.
	http.HandleFunc("/api/nodes", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		nodesResult, err := cli.NodeList(ctx, client.NodeListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		httpClient := &http.Client{Timeout: 5 * time.Second}
		nodes := nodesResult.Items
		result := make([]DashboardNode, len(nodes))

		type indexedStats struct {
			i  int
			ns *NodeStats
		}
		ch := make(chan indexedStats, len(nodes))

		for i, node := range nodes {
			isLeader := node.ManagerStatus != nil && node.ManagerStatus.Leader
			result[i] = DashboardNode{Details: NodeDetails{
				ID:          node.ID,
				Hostname:    node.Description.Hostname,
				State:       string(node.Status.State),
				Addr:        node.Status.Addr,
				Role:        string(node.Spec.Role),
				IsLeader:    isLeader,
				NanoCPUs:    node.Description.Resources.NanoCPUs,
				MemoryBytes: node.Description.Resources.MemoryBytes,
			}}
			if node.Status.State != swarm.NodeStateReady {
				ch <- indexedStats{i: i}
				continue
			}
			go func(i int, node swarm.Node) {
				url := "http://" + node.Status.Addr + ":8080/api/local-stats"
				resp, err := httpClient.Get(url)
				if err != nil {
					log.Printf("could not reach node %s: %v", node.Description.Hostname, err)
					ch <- indexedStats{i: i}
					return
				}
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				var ns NodeStats
				if err := json.Unmarshal(body, &ns); err != nil {
					log.Printf("could not parse stats from node %s: %v", node.Description.Hostname, err)
					ch <- indexedStats{i: i}
					return
				}
				ch <- indexedStats{i: i, ns: &ns}
			}(i, node)
		}

		for range nodes {
			s := <-ch
			if s.ns != nil {
				result[s.i].Containers = s.ns.Containers
				result[s.i].TempCelsius = s.ns.TempCelsius
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Serve the static frontend from the public/ directory.
	http.Handle("/", gzipFileServer("public"))

	// Returns CPU and memory stats for all swarm containers running on this node.
	// Called by peer nodes and directly by the browser during development.
	http.HandleFunc("/api/local-stats", func(w http.ResponseWriter, r *http.Request) {
		stats, err := getLocalStats(cli)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	log.Println("Swarm dashboard running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
