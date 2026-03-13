package main

import (
	"encoding/json"
	"net/http"

	"github.com/moby/moby/api/types/swarm"
)

func mockServices() []swarm.Service {
	replicas := uint64(2)
	return []swarm.Service{
		{
			ID: "abc123def456",
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{Name: "web"},
				Mode: swarm.ServiceMode{
					Replicated: &swarm.ReplicatedService{Replicas: &replicas},
				},
			},
		},
		{
			ID: "def456abc789",
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{Name: "api"},
				Mode: swarm.ServiceMode{
					Replicated: &swarm.ReplicatedService{Replicas: &replicas},
				},
			},
		},
		{
			ID: "ghi789jkl012",
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{Name: "db"},
				Mode: swarm.ServiceMode{
					Replicated: &swarm.ReplicatedService{Replicas: &replicas},
				},
			},
		},
	}
}

func mockNodes() []DashboardNode {
	temp1 := 42.5
	temp2 := 38.0
	temp3 := 55.1
	return []DashboardNode{
		{
			Details: swarm.Node{
				ID: "node1",
				Description: swarm.NodeDescription{
					Hostname:  "manager-1",
					Resources: swarm.Resources{NanoCPUs: 4000000000, MemoryBytes: 8 * 1024 * 1024 * 1024},
				},
				Status:        swarm.NodeStatus{State: swarm.NodeStateReady, Addr: "10.0.0.1"},
				ManagerStatus: &swarm.ManagerStatus{Leader: true},
				Spec:          swarm.NodeSpec{Role: swarm.NodeRoleManager},
			},
			TempCelsius: &temp1,
			Containers: []ContainerStats{
				{ID: "c1a2b3", Name: "monitoring_mini-swarm-dash.7q4ru95wcljxkr3vh2l31ytof.ixfs6nco8paq25r393echzs5s", ServiceName: "monitoring_mini-swarm-dash", CPUPercent: 12.4, MemUsage: 128 * 1024 * 1024, MemLimit: 512 * 1024 * 1024},
				{ID: "c4d5e6", Name: "api.1.abc", ServiceName: "api", CPUPercent: 5.1, MemUsage: 64 * 1024 * 1024, MemLimit: 256 * 1024 * 1024},
			},
		},
		{
			Details: swarm.Node{
				ID: "node2",
				Description: swarm.NodeDescription{
					Hostname:  "worker-1",
					Resources: swarm.Resources{NanoCPUs: 4000000000, MemoryBytes: 8 * 1024 * 1024 * 1024},
				},
				Status: swarm.NodeStatus{State: swarm.NodeStateReady, Addr: "10.0.0.2"},
				Spec:   swarm.NodeSpec{Role: swarm.NodeRoleWorker},
			},
			TempCelsius: &temp2,
			Containers: []ContainerStats{
				{ID: "f7g8h9", Name: "web.2.def", ServiceName: "web", CPUPercent: 8.9, MemUsage: 110 * 1024 * 1024, MemLimit: 512 * 1024 * 1024},
				{ID: "i1j2k3", Name: "db.1.ghi", ServiceName: "db", CPUPercent: 22.3, MemUsage: 400 * 1024 * 1024, MemLimit: 1024 * 1024 * 1024},
			},
		},
		{
			Details: swarm.Node{
				ID: "node3",
				Description: swarm.NodeDescription{
					Hostname:  "worker-2",
					Resources: swarm.Resources{NanoCPUs: 2000000000, MemoryBytes: 4 * 1024 * 1024 * 1024},
				},
				Status: swarm.NodeStatus{State: swarm.NodeStateReady, Addr: "10.0.0.3"},
				Spec:   swarm.NodeSpec{Role: swarm.NodeRoleWorker},
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

func registerMockHandlers() {
	http.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockServices())
	})
	http.HandleFunc("/api/nodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockNodes())
	})
	http.HandleFunc("/api/local-stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockLocalStats())
	})
}
