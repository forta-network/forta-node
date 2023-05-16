package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAgentResourceLimits(t *testing.T) {
	r := require.New(t)

	limits := GetAgentResourceLimits(ResourcesConfig{})
	r.Equal(CPUsToMicroseconds(0.2), limits.CPUQuota)
	r.Equal(MiBToBytes(10000), limits.Memory)
}

func TestGetAgentResourceLimits_CustomValues(t *testing.T) {
	r := require.New(t)

	limits := GetAgentResourceLimits(ResourcesConfig{
		AgentMaxMemoryMiB: 12,
		AgentMaxCPUs:      0.1,
	})
	r.Equal(CPUsToMicroseconds(0.1), limits.CPUQuota)
	r.Equal(MiBToBytes(12), limits.Memory)
}
