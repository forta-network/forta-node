package config

import "runtime"

// AgentResourceLimits contain the agent resource limits data.
type AgentResourceLimits struct {
	CPUQuota int64 // in microseconds
	Memory   int64 // in bytes
}

// GetAgentResourceLimits calculates and returns the resource limits by
// taking the configuration into account. Zero values mean no limits.
func GetAgentResourceLimits(resourcesCfg ResourcesConfig) *AgentResourceLimits {
	var limits AgentResourceLimits

	if resourcesCfg.DisableAgentLimits {
		return &limits
	}

	limits.CPUQuota = getDefaultCPUQuotaPerAgent()
	if resourcesCfg.AgentMaxCPUs > 0 {
		limits.CPUQuota = int64(resourcesCfg.AgentMaxCPUs * float64(100000))
	}

	limits.Memory = getDefaultMemoryPerAgent()
	if resourcesCfg.AgentMaxMemoryMiB > 0 {
		limits.Memory = int64(resourcesCfg.AgentMaxMemoryMiB * 104858)
	}

	return &limits
}

// Below calculations are made by taking AWS EC2 t2.medium as a reference.
// The agents we run nowadays use around 60 MiBs of memory and we can raise that to 100 MiB.
// After reserving 1 GiB for running Forta node runner/supervisor and the containers:
// (4000 - 1000)/100 = 30 agents
// 100000 / 30 ≈ 3333 microseconds (quota coefficient)

// getDefaultCPUQuotaPerAgent returns the default CFS microseconds value allowed per agent.
func getDefaultCPUQuotaPerAgent() int64 {
	return int64(runtime.NumCPU() * 3333) // to microseconds
}

// getDefaultMemoryPerAgent returns the constant default memory allowed per agent.
func getDefaultMemoryPerAgent() int64 {
	return 104858000 // 100 MiB in bytes
}
