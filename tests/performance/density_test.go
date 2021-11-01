//+build perf_test

package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services/query"
)

func agentId(index int) string {
	id := crypto.Keccak256Hash([]byte(fmt.Sprintf("agent-%d", index)))
	return id.Hex()
}

type TestConfig struct {
	host           string
	agentCount     int
	start          int64
	end            int64
	rate           int64
	expectedAlerts int64
	image          string
	manifest       string
}

type TestContext struct {
	t         *testing.T
	cfg       *TestConfig
	msgClient clients.MessageClient
	startDate time.Time

	ready []config.AgentConfig
	agts  []*config.AgentConfig
}

func (tc *TestContext) generateAgents(count int) []*config.AgentConfig {
	var agts []*config.AgentConfig
	for i := 0; i < count; i++ {
		agts = append(agts, &config.AgentConfig{
			ID:       agentId(i),
			Image:    tc.cfg.image,
			Manifest: tc.cfg.manifest,
		})
	}
	return agts
}

func (tc *TestContext) runBlocks() error {
	url := fmt.Sprintf("http://%s:8989/start?start=%d&end=%d&rate=%d",
		tc.cfg.host,
		tc.cfg.start,
		tc.cfg.end,
		tc.cfg.rate,
	)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s returned %d", url, resp.StatusCode)
	}
	return err
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
				// sanity sleep
				time.Tick(5 * time.Second)
				return nil
			}
		}
	}
}

func (tc *TestContext) Setup() error {
	tc.startDate = time.Now()
	tc.agts = tc.generateAgents(tc.cfg.agentCount)
	tc.msgClient.Subscribe(messaging.SubjectAgentsStatusAttached, messaging.AgentsHandler(tc.handleReady))
	tc.runAgents()
	return tc.waitForReady(5 * time.Minute)
}

func (tc *TestContext) runAgents() {
	tc.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, tc.agts)
}

func (tc *TestContext) getResults() (*query.AgentReport, error) {
	url := fmt.Sprintf("http://%s:8778/report/agents?startDate=%d000",
		tc.cfg.host,
		tc.startDate.Unix(),
	)
	tc.t.Logf("report: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s returned %d", url, resp.StatusCode)
	}
	var report query.AgentReport
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (tc *TestContext) verifyResults() error {
	// wait enough time for blocks to be processed
	blockCount := tc.cfg.end - tc.cfg.start
	<-time.Tick(time.Duration((blockCount+1)*tc.cfg.rate) * time.Millisecond)
	results, err := tc.getResults()
	if err != nil {
		return err
	}

	for _, agt := range tc.agts {
		if r, ok := results.AlertCounts[agt.ID]; !ok {
			return fmt.Errorf("%s not found", agt.ID)
		} else {
			if r != tc.cfg.expectedAlerts {
				return fmt.Errorf("%s had %d instead of %d alerts", agt.ID, r, tc.cfg.expectedAlerts)
			}
		}
	}
	return nil
}

func (tc *TestContext) cleanUp() error {
	tc.msgClient.Publish(messaging.SubjectAgentsActionStop, tc.agts)
	return nil
}

func NewTestContext(t *testing.T, cfg *TestConfig) *TestContext {
	return &TestContext{
		t:         t,
		cfg:       cfg,
		msgClient: messaging.NewClient("perf-test", fmt.Sprintf("%s:4222", cfg.host)),
		ready:     nil,
	}
}

func TestPerformance(t *testing.T) {
	tctx := NewTestContext(t, &TestConfig{
		host:           "localhost", //"54.90.96.23",
		image:          "disco.forta.network/bafybeibwzulzj5ua46w5gjwulivrvjbp24blio4tz4zlyzgu4pp6o7qpjy@sha256:a423779dfc43e3588579f5aa703d074413c734cb24495334776e01749f63dda9",
		manifest:       "QmReurJ6XsKQNkWxw7DaSTTnZcmZia2P9J7ptUQo8DT3Mk",
		agentCount:     3,
		start:          13513750,
		end:            13513753,
		rate:           15000,
		expectedAlerts: 318,
	})

	assert.NoError(t, tctx.Setup())
	assert.NoError(t, tctx.runBlocks())
	err := tctx.verifyResults()
	tctx.cleanUp()
	assert.NoError(t, err)

}
