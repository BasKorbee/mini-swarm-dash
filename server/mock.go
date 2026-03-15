//go:build mock

package main

import (
	"encoding/json"
	"net/http"
)

func mockServices() []ServiceSummary {
	replicas := uint64(2)
	return []ServiceSummary{
		{ID: "abc123def456", Name: "web", Image: "nginx:alpine", Replicas: &replicas},
		{ID: "def456abc789", Name: "api", Image: "myapp/api:latest", Replicas: &replicas},
		{ID: "ghi789jkl012", Name: "db", Image: "postgres:16", Replicas: &replicas},
	}
}

func mockNodes() []DashboardNode {
	temp1 := 42.5
	temp2 := 38.0
	temp3 := 55.1
	return []DashboardNode{
		{
			Details: NodeDetails{
				ID: "node1", Hostname: "manager-1", State: "ready", Addr: "10.0.0.1",
				Role: "manager", IsLeader: true, NanoCPUs: 4000000000, MemoryBytes: 8 * 1024 * 1024 * 1024,
			},
			TempCelsius: &temp1,
			Containers: []ContainerStats{
				{ID: "c1a2b3", Name: "monitoring_mini-swarm-dash.7q4ru95wcljxkr3vh2l31ytof.ixfs6nco8paq25r393echzs5s", ServiceName: "monitoring_mini-swarm-dash", CPUPercent: 12.4, MemUsage: 128 * 1024 * 1024, MemLimit: 512 * 1024 * 1024},
				{ID: "c4d5e6", Name: "api.1.abc", ServiceName: "api", CPUPercent: 5.1, MemUsage: 64 * 1024 * 1024, MemLimit: 256 * 1024 * 1024},
			},
		},
		{
			Details: NodeDetails{
				ID: "node4", Hostname: "manager-2", State: "ready", Addr: "10.0.0.2",
				Role: "manager", NanoCPUs: 4000000000, MemoryBytes: 8 * 1024 * 1024 * 1024,
			},
			TempCelsius: &temp2,
			Containers: []ContainerStats{
				{ID: "f7g8h9", Name: "web.2.def", ServiceName: "web", CPUPercent: 8.9, MemUsage: 110 * 1024 * 1024, MemLimit: 512 * 1024 * 1024},
				{ID: "i1j2k3", Name: "db.1.ghi", ServiceName: "db", CPUPercent: 22.3, MemUsage: 400 * 1024 * 1024, MemLimit: 1024 * 1024 * 1024},
			},
		},
		{
			Details: NodeDetails{
				ID: "node2", Hostname: "worker-1", State: "ready", Addr: "10.0.0.2",
				Role: "worker", NanoCPUs: 4000000000, MemoryBytes: 8 * 1024 * 1024 * 1024,
			},
			TempCelsius: &temp2,
			Containers: []ContainerStats{
				{ID: "f7g8h9", Name: "web.2.def", ServiceName: "web", CPUPercent: 8.9, MemUsage: 110 * 1024 * 1024, MemLimit: 512 * 1024 * 1024},
				{ID: "i1j2k3", Name: "db.1.ghi", ServiceName: "db", CPUPercent: 22.3, MemUsage: 400 * 1024 * 1024, MemLimit: 1024 * 1024 * 1024},
			},
		},
		{
			Details: NodeDetails{
				ID: "node3", Hostname: "worker-2", State: "ready", Addr: "10.0.0.3",
				Role: "worker", NanoCPUs: 2000000000, MemoryBytes: 4 * 1024 * 1024 * 1024,
			},
			TempCelsius: &temp3,
			Containers: []ContainerStats{
				{ID: "l4m5n6", Name: "api.2.jkl", ServiceName: "api", CPUPercent: 3.7, MemUsage: 72 * 1024 * 1024, MemLimit: 256 * 1024 * 1024},
			},
		},
	}
}

func mockLocalStats() *NodeStats {
	temp := 42.5
	return &NodeStats{
		NodeID:      "node1",
		Hostname:    "manager-1",
		TempCelsius: &temp,
		Containers: []ContainerStats{
			{ID: "c1a2b3", Name: "web.1.xyz", ServiceName: "web", CPUPercent: 12.4, MemUsage: 128 * 1024 * 1024, MemLimit: 512 * 1024 * 1024},
			{ID: "c4d5e6", Name: "api.1.abc", ServiceName: "api", CPUPercent: 5.1, MemUsage: 64 * 1024 * 1024, MemLimit: 256 * 1024 * 1024},
		},
	}
}

func mockHandleServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockServices())
	}
}

func mockHandleNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockNodes())
	}
}

func mockHandleLocalStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockLocalStats())
	}
}
