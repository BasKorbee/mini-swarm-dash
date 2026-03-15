package main

// cpuSample holds the last raw CPU counters for a container so we can compute
// a delta on the next poll without needing IncludePreviousSample (which blocks ~1s).
type cpuSample struct {
	totalUsage  uint64
	systemUsage uint64
	onlineCPUs  uint32
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
	ID          string `json:"id"`
	Hostname    string `json:"hostname"`
	State       string `json:"state"`
	Addr        string `json:"addr"`
	Role        string `json:"role"`
	IsLeader    bool   `json:"is_leader"`
	NanoCPUs    int64  `json:"nano_cpus"`
	MemoryBytes int64  `json:"memory_bytes"`
}

// ServiceSummary holds the subset of swarm.Service fields used by the frontend.
type ServiceSummary struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	IsGlobal bool    `json:"is_global"`
	Replicas *uint64 `json:"replicas"`
}

// DashboardNode is the enriched node sent to the frontend
type DashboardNode struct {
	Details     NodeDetails      `json:"details"`
	TempCelsius *float64         `json:"temp_celsius"`
	Containers  []ContainerStats `json:"containers"`
}
