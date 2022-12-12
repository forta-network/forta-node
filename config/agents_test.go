package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_ContainerName(t *testing.T) {
	cfg := AgentConfig{
		ID:      "0x04f65c638f234548104790d7c692c9273d41f82d784b174ff2fdc3e8e5bf1636",
		Image:   "bafybeibvkqkf7i3c5ouehviwjb2dzbukgqied3cg36axl7gzm23r6ielnu@sha256:de866feeb97cba4cad6343c4137cb48bc798be0136015bec16d97c8ef28852b9",
		ShardId: 1,
	}
	assert.Equal(t, "forta-agent-0x04f65c-de86-1", cfg.ContainerName())
}
