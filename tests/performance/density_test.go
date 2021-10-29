package performance

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func generateAgents(count int) []*config.AgentConfig {
	var agts []*config.AgentConfig
	for i := 0; i < count; i++ {
		agts = append(agts, &config.AgentConfig{
			ID:       agentId(i),
			Image:    "disco.forta.network/bafybeibwzulzj5ua46w5gjwulivrvjbp24blio4tz4zlyzgu4pp6o7qpjy@sha256:a423779dfc43e3588579f5aa703d074413c734cb24495334776e01749f63dda9",
			Manifest: "QmReurJ6XsKQNkWxw7DaSTTnZcmZia2P9J7ptUQo8DT3Mk",
		})
	}
	return agts
}

func agentId(index int) string {
	id := crypto.Keccak256Hash([]byte(fmt.Sprintf("agent-%d", index)))
	return id.Hex()
}

type TestConfig struct {
	agentCount int
	start      int64
	end        int64
	rate       int64
}

type TestContext struct {
	t         *testing.T
	cfg       *TestConfig
	msgClient clients.MessageClient

	ready []config.AgentConfig
}

func (tc *TestContext) runBlocks() {
	url := fmt.Sprintf("http://localhost:8989/start?start=%d&end=%d&rate=%d",
		tc.cfg.start,
		tc.cfg.end,
		tc.cfg.rate,
	)
	resp, err := http.Get(url)
	assert.NoError(tc.t, err)
	assert.Equal(tc.t, 200, resp.StatusCode)
}

func (tc *TestContext) handleReady(payload messaging.AgentPayload) error {
	tc.ready = append(tc.ready, payload...)
	return nil
}

func (tc *TestContext) waitForReady(duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("agent start failed: loaded %d of %d agents", len(tc.ready), tc.cfg.agentCount)
		default:
			if len(tc.ready) == tc.cfg.agentCount {
				return nil
			}
		}
	}
}

func (tc *TestContext) Setup() error {
	agts := generateAgents(tc.cfg.agentCount)
	tc.msgClient.Subscribe(messaging.SubjectAgentReady, messaging.AgentsHandler(tc.handleReady))
	tc.runAgents(agts)
	return tc.waitForReady(5 * time.Minute)
}

func (tc *TestContext) runAgents(agts []*config.AgentConfig) {
	tc.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, agts)
}

func NewTestContext(t *testing.T, cfg *TestConfig) *TestContext {
	return &TestContext{
		t:         t,
		cfg:       cfg,
		msgClient: messaging.NewClient("perf-test", "localhost:4222"),
		ready:     nil,
	}
}

func TestPerformance(t *testing.T) {
	startDate := time.Now()
	tctx := NewTestContext(t, &TestConfig{
		agentCount: 3,
		start:      13513743,
		end:        13513753,
		rate:       15000,
	})
	assert.NoError(t, tctx.Setup())
	tctx.runBlocks()
	t.Logf("startDate: %d000", startDate.Unix())
	// start agents
	// are they started?
	// trigger start
}
