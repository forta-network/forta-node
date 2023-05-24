package lifecycle

import (
	"testing"

	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/require"
)

func TestFindMissingBots(t *testing.T) {
	r := require.New(t)

	list1 := []config.AgentConfig{
		{
			ID: "10",
		},
		{
			ID: "20",
		},
		{
			ID: "30",
		},
		{
			ID: "40",
		},
	}

	list2 := []config.AgentConfig{
		{
			ID: "20",
		},
		{
			ID: "30",
		},
	}

	result := FindMissingBots(list1, list2)
	r.Len(result, 2)
	r.Equal("10", result[0].ID)
	r.Equal("40", result[1].ID)
}

func TestFindExtraBots(t *testing.T) {
	r := require.New(t)

	list1 := []config.AgentConfig{
		{
			ID: "10",
		},
		{
			ID: "20",
		},
		{
			ID: "40",
		},
	}

	list2 := []config.AgentConfig{
		{
			ID: "10",
		},
		{
			ID: "30",
		},
		{
			ID: "50",
		},
	}

	result := FindExtraBots(list1, list2)
	r.Len(result, 2)
	r.Equal("30", result[0].ID)
	r.Equal("50", result[1].ID)
}

func TestUpdatedBots(t *testing.T) {
	r := require.New(t)

	list1 := []config.AgentConfig{
		{
			ID: "10",
		},
		{
			ID: "20",
		},
		{
			ID: "40",
		},
	}

	list2 := []config.AgentConfig{
		{
			ID: "10",
			ShardConfig: &config.ShardConfig{
				ShardID: 1,
			},
		},
		{
			ID: "30",
		},
		{
			ID: "50",
		},
	}

	result := FindUpdatedBots(list1, list2)
	r.Len(result, 1)
	r.Equal("10", result[0].ID)
}

func TestDrop(t *testing.T) {
	r := require.New(t)

	list1 := []config.AgentConfig{
		{
			ID: "10",
		},
		{
			ID: "20",
		},
		{
			ID: "40",
		},
	}

	result := Drop(config.AgentConfig{
		ID: "20",
	}, list1)
	r.Len(result, 2)
	r.Equal("10", result[0].ID)
	r.Equal("40", result[1].ID)
}
