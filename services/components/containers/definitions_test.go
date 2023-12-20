package containers

import (
	"testing"

	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/require"
)

func TestContainerEnvVar_ChainID(t *testing.T) {
	r := require.New(t)

	botConfig := config.AgentConfig{
		ChainID: 0,
	}
	containerConfig := NewBotContainerConfig("", botConfig, config.LogConfig{}, config.ResourcesConfig{}, "")
	env := containerConfig.Env
	r.Equal("", env[config.EnvFortaChainID])

	botConfig = config.AgentConfig{
		ChainID: 137,
	}
	containerConfig = NewBotContainerConfig("", botConfig, config.LogConfig{}, config.ResourcesConfig{}, "")
	env = containerConfig.Env
	r.Equal("137", env[config.EnvFortaChainID])
}

func TestContainerEnvVar_Sharding(t *testing.T) {
	r := require.New(t)

	botConfig := config.AgentConfig{
		ShardConfig: nil,
	}
	containerConfig := NewBotContainerConfig("", botConfig, config.LogConfig{}, config.ResourcesConfig{}, "")
	env := containerConfig.Env
	r.Equal("", env[config.EnvFortaShardID])
	r.Equal("", env[config.EnvFortaShardCount])

	botConfig = config.AgentConfig{
		ShardConfig: &config.ShardConfig{
			ShardID: 0,
			Shards:  2,
			Target:  3,
		},
	}
	containerConfig = NewBotContainerConfig("", botConfig, config.LogConfig{}, config.ResourcesConfig{}, "")
	env = containerConfig.Env
	r.Equal("0", env[config.EnvFortaShardID])
	r.Equal("2", env[config.EnvFortaShardCount])
}
