package agentlogs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	r := require.New(t)

	agents := Agents{
		{
			ID:   "123",
			Logs: "hello world",
		},
	}
	reader, err := Encode(agents)
	r.NoError(err)

	decodedAgents, err := Decode(reader)
	r.NoError(err)

	r.Equal(agents[0].ID, decodedAgents[0].ID)
	r.Equal(agents[0].Logs, decodedAgents[0].Logs)
}
