package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_ContainerName(t *testing.T) {
	cfg := AgentConfig{
		ID:    "0x04f65c638f234548104790d7c692c9273d41f82d784b174ff2fdc3e8e5bf1636",
		Image: "bafybeibvkqkf7i3c5ouehviwjb2dzbukgqied3cg36axl7gzm23r6ielnu@sha256:de866feeb97cba4cad6343c4137cb48bc798be0136015bec16d97c8ef28852b9",
	}
	assert.Equal(t, "forta-agent-0x04f65c-de86", cfg.ContainerName())
}

func TestAgentConfig_ContainerNameSharded(t *testing.T) {
	cfg := AgentConfig{
		ID:    "0x04f65c638f234548104790d7c692c9273d41f82d784b174ff2fdc3e8e5bf1636",
		Image: "bafybeibvkqkf7i3c5ouehviwjb2dzbukgqied3cg36axl7gzm23r6ielnu@sha256:de866feeb97cba4cad6343c4137cb48bc798be0136015bec16d97c8ef28852b9",
		ShardConfig: &ShardConfig{
			Shards:  5,
			ShardID: 3,
			Target:  3,
		},
	}
	assert.Equal(t, "forta-agent-0x04f65c-de86-3", cfg.ContainerName())
}

func TestAgentConfig_Equal(t *testing.T) {
	tests := []struct {
		name string
		ac1  AgentConfig
		ac2  AgentConfig
		want bool
	}{
		{
			name: "Two identical AgentConfigs",
			ac1: AgentConfig{
				ID:           "agent1",
				Image:        "image1",
				Manifest:     "manifest1",
				IsLocal:      true,
				IsStandalone: false,
				StartBlock:   nil,
				StopBlock:    nil,
				Owner:        "owner1",
				ChainID:      123,
				ShardConfig: &ShardConfig{
					ShardID: 1,
					Shards:  2,
					Target:  3,
				},
			},
			ac2: AgentConfig{
				ID:           "agent1",
				Image:        "image1",
				Manifest:     "manifest1",
				IsLocal:      true,
				IsStandalone: false,
				StartBlock:   nil,
				StopBlock:    nil,
				Owner:        "owner1",
				ChainID:      123,
				ShardConfig: &ShardConfig{
					ShardID: 1,
					Shards:  2,
					Target:  3,
				},
			},
			want: true,
		},
		{
			name: "Two AgentConfigs with different ShardConfig fields",
			ac1: AgentConfig{
				ID:           "agent1",
				Image:        "image1",
				Manifest:     "manifest1",
				IsLocal:      true,
				IsStandalone: false,
				StartBlock:   nil,
				StopBlock:    nil,
				Owner:        "owner1",
				ChainID:      123,
				ShardConfig: &ShardConfig{
					ShardID: 1,
					Shards:  2,
					Target:  3,
				},
			},
			ac2: AgentConfig{
				ID:           "agent1",
				Image:        "image1",
				Manifest:     "manifest1",
				IsLocal:      true,
				IsStandalone: false,
				StartBlock:   nil,
				StopBlock:    nil,
				Owner:        "owner1",
				ChainID:      123,
				ShardConfig: &ShardConfig{
					ShardID: 2,
					Shards:  4,
					Target:  5,
				},
			},
			want: false,
		},
		{
			name: "Two AgentConfigs with different Manifest fields",
			ac1: AgentConfig{
				ID:           "agent1",
				Image:        "image1",
				Manifest:     "manifest1",
				IsLocal:      true,
				IsStandalone: false,
				StartBlock:   nil,
				StopBlock:    nil,
				Owner:        "owner1",
				ChainID:      123,
				ShardConfig: &ShardConfig{
					ShardID: 1,
					Shards:  2,
					Target:  3,
				},
			},
			ac2: AgentConfig{
				ID:           "agent1",
				Image:        "image1",
				Manifest:     "manifest2",
				IsLocal:      true,
				IsStandalone: false,
				StartBlock:   nil,
				StopBlock:    nil,
				Owner:        "owner1",
				ChainID:      123,
				ShardConfig: &ShardConfig{
					ShardID: 1,
					Shards:  2,
					Target:  3,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.ac1.Equal(tt.ac2); got != tt.want {
					t.Errorf("Equal() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
