package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"
)

func handleServices(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := cli.ServiceList(context.Background(), client.ServiceListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := make([]ServiceSummary, len(services.Items))
		for i, s := range services.Items {
			result[i] = toServiceSummary(s)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func toServiceSummary(s swarm.Service) ServiceSummary {
	ss := ServiceSummary{
		ID:       s.ID,
		Name:     s.Spec.Name,
		Image:    s.Spec.TaskTemplate.ContainerSpec.Image,
		IsGlobal: s.Spec.Mode.Global != nil,
	}
	if s.Spec.Mode.Replicated != nil {
		ss.Replicas = s.Spec.Mode.Replicated.Replicas
	}
	return ss
}

func handleNodes(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		nodesResult, err := cli.NodeList(ctx, client.NodeListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nodes := nodesResult.Items
		result := make([]DashboardNode, len(nodes))
		for i, node := range nodes {
			result[i] = DashboardNode{Details: toNodeDetails(node)}
		}

		result = enrichNodesWithStats(nodes, result)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func toNodeDetails(node swarm.Node) NodeDetails {
	return NodeDetails{
		ID:          node.ID,
		Hostname:    node.Description.Hostname,
		State:       string(node.Status.State),
		Addr:        node.Status.Addr,
		Role:        string(node.Spec.Role),
		IsLeader:    node.ManagerStatus != nil && node.ManagerStatus.Leader,
		NanoCPUs:    node.Description.Resources.NanoCPUs,
		MemoryBytes: node.Description.Resources.MemoryBytes,
	}
}

// enrichNodesWithStats fans out to each ready node's /api/local-stats endpoint
// and merges the container stats and temperature into result.
func enrichNodesWithStats(nodes []swarm.Node, result []DashboardNode) []DashboardNode {
	httpClient := &http.Client{Timeout: 5 * time.Second}

	type indexedStats struct {
		i  int
		ns *NodeStats
	}
	ch := make(chan indexedStats, len(nodes))

	for i, node := range nodes {
		if node.Status.State != swarm.NodeStateReady {
			ch <- indexedStats{i: i}
			continue
		}
		go func(i int, node swarm.Node) {
			ns, err := fetchNodeStats(httpClient, node)
			if err != nil {
				logger.Warn("could not reach node", "node", node.Description.Hostname, "err", err)
				ch <- indexedStats{i: i}
				return
			}
			ch <- indexedStats{i: i, ns: ns}
		}(i, node)
	}

	for range nodes {
		s := <-ch
		if s.ns != nil {
			result[s.i].Containers = s.ns.Containers
			result[s.i].TempCelsius = s.ns.TempCelsius
		}
	}
	return result
}

// fetchNodeStats retrieves the NodeStats from a peer node's /api/local-stats endpoint.
func fetchNodeStats(httpClient *http.Client, node swarm.Node) (*NodeStats, error) {
	url := "http://" + node.Status.Addr + ":8080/api/local-stats"
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var ns NodeStats
	if err := json.Unmarshal(body, &ns); err != nil {
		logger.Warn("could not parse stats from node", "node", node.Description.Hostname, "err", err)
		return nil, err
	}
	return &ns, nil
}

func handleLocalStats(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := getLocalStats(cli)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
