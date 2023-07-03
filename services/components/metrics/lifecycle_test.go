package metrics

import (
	"testing"

	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/require"
)

func TestFromBotConfigs(t *testing.T) {
	r := require.New(t)

	metricName := "test-metric-name"
	details := ""
	bot1Sharded := config.AgentConfig{
		ID: "0x1",
		ShardConfig: &config.ShardConfig{
			ShardID: 1,
			Shards:  2,
			Target:  2,
		},
	}
	bot2NotSharded := config.AgentConfig{
		ID: "0x2",
	}

	metrics := fromBotConfigs(metricName, details, []config.AgentConfig{
		bot1Sharded, bot2NotSharded,
	})

	r.Equal(bot1Sharded.ID, metrics[0].AgentId)
	r.NotEmpty(metrics[0].Timestamp)
	r.Equal(metricName, metrics[0].Name)
	r.Equal(int32(1), metrics[0].ShardId)
	r.Equal(float64(1), metrics[0].Value)

	r.Equal(bot2NotSharded.ID, metrics[1].AgentId)
	r.NotEmpty(metrics[1].Timestamp)
	r.Equal(metricName, metrics[1].Name)
	r.Equal(details, metrics[1].Details)
	r.Equal(float64(1), metrics[1].Value)
}
