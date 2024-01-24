package docker

type ContainerResources struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
	} `json:"cpu_stats"`
	MemoryStats struct {
		Usage int `json:"usage"`
	} `json:"memory_stats"`
	NetworkStats map[string]struct {
		// Bytes received
		RxBytes uint64 `json:"rx_bytes"`
		// Bytes sent
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"networks"`
}
