package docker

type ContainerResources struct {
	CPUStats     CPUStats                `json:"cpu_stats"`
	MemoryStats  MemoryStats             `json:"memory_stats"`
	NetworkStats map[string]NetworkStats `json:"networks"`
}

type CPUStats struct {
	CPUUsage CPUUsage `json:"cpu_usage"`
}

type CPUUsage struct {
	TotalUsage uint64 `json:"total_usage"`
}

type MemoryStats struct {
	Usage uint64 `json:"usage"`
}

type NetworkStats struct {
	// Bytes received
	RxBytes uint64 `json:"rx_bytes"`
	// Bytes sent
	TxBytes uint64 `json:"tx_bytes"`
}
