package poolagent

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/services/scanner"
	"google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
)

// Constants
const (
	DefaultBufferSize = 2000
)

// Agent receives blocks and transactions, and produces results.
type Agent struct {
	config config.AgentConfig

	txRequests    chan *protocol.EvaluateTxRequest // never closed - deallocated when agent is discarded
	txResults     chan<- *scanner.TxResult
	blockRequests chan *protocol.EvaluateBlockRequest // never closed - deallocated when agent is discarded
	blockResults  chan<- *scanner.BlockResult

	errCounter *errorCounter
	msgClient  clients.MessageClient

	client    clients.AgentClient
	ready     chan struct{}
	readyOnce sync.Once
	closed    chan struct{}
	closeOnce sync.Once
}

// New creates a new agent.
func New(agentCfg config.AgentConfig, msgClient clients.MessageClient, txResults chan<- *scanner.TxResult, blockResults chan<- *scanner.BlockResult) *Agent {
	return &Agent{
		config:        agentCfg,
		txRequests:    make(chan *protocol.EvaluateTxRequest, DefaultBufferSize),
		txResults:     txResults,
		blockRequests: make(chan *protocol.EvaluateBlockRequest, DefaultBufferSize),
		blockResults:  blockResults,
		errCounter:    NewErrorCounter(3, isCriticalErr),
		msgClient:     msgClient,
		ready:         make(chan struct{}),
		closed:        make(chan struct{}),
	}
}

func isCriticalErr(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, codes.DeadlineExceeded.String()) ||
		strings.Contains(errStr, codes.Unavailable.String())
}

// LogStatus logs the status of the agent.
func (agent *Agent) LogStatus() {
	log.WithFields(log.Fields{
		"agent":         agent.config.ID,
		"buffer-blocks": len(agent.blockRequests),
		"buffer-txs":    len(agent.txRequests),
		"ready":         agent.IsReady(),
		"closed":        agent.IsClosed(),
	}).Debug("agent status")
}

// Config returns the agent config.
func (agent *Agent) Config() config.AgentConfig {
	return agent.config
}

// TxRequestCh returns the transaction request channel safely.
func (agent *Agent) TxRequestCh() chan<- *protocol.EvaluateTxRequest {
	return agent.txRequests
}

// BlockRequestCh returns the block request channel safely.
func (agent *Agent) BlockRequestCh() chan<- *protocol.EvaluateBlockRequest {
	return agent.blockRequests
}

// Close implements io.Closer.
func (agent *Agent) Close() error {
	agent.closeOnce.Do(func() {
		close(agent.closed) // never close this anywhere else
		agent.client.Close()
	})
	return nil
}

// SetReady sets the agent ready.
func (agent *Agent) SetReady() {
	agent.readyOnce.Do(func() {
		close(agent.ready) // never close this anywhere else
	})
}

// Ready returns the ready channel.
func (agent *Agent) Ready() <-chan struct{} {
	return agent.ready
}

// Closed returns the closed channel.
func (agent *Agent) Closed() <-chan struct{} {
	return agent.closed
}

// IsReady tells if the agent is ready.
func (agent *Agent) IsReady() bool {
	return isChanClosed(agent.ready)
}

// IsClosed tells if the agent is closed.
func (agent *Agent) IsClosed() bool {
	return isChanClosed(agent.closed)
}

func isChanClosed(ch chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

// SetClient sets the agent client for sending the requests.
func (agent *Agent) SetClient(agentClient clients.AgentClient) {
	agent.client = agentClient
}

// StartProcessing launches the goroutines to concurrently process incoming requests
// from request channels.
func (agent *Agent) StartProcessing() {
	go agent.processTransactions()
	go agent.processBlocks()
}

func (agent *Agent) processTransactions() {
	lg := log.WithFields(log.Fields{
		"agent":     agent.config.ID,
		"component": "agent",
		"evaluate":  "transaction",
	})
	for request := range agent.txRequests {
		startTime := time.Now()
		if agent.IsClosed() {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
		resp, err := agent.client.EvaluateTx(ctx, request)
		cancel()
		if err == nil {
			var duration time.Duration
			resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
			lg.WithField("duration", duration).Debugf("request successful")
			resp.Metadata["imageHash"] = agent.config.ImageHash()
			agent.txResults <- &scanner.TxResult{
				AgentConfig: agent.config,
				Request:     request,
				Response:    resp,
			}
			lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
			continue
		}
		lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
		if agent.errCounter.TooManyErrs(err) {
			lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
			agent.Close()
			agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agent.config})
			return
		}
	}
}

func (agent *Agent) processBlocks() {
	lg := log.WithFields(log.Fields{
		"agent":     agent.config.ID,
		"component": "agent",
		"evaluate":  "block",
	})
	for request := range agent.blockRequests {
		startTime := time.Now()
		if agent.IsClosed() {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
		resp, err := agent.client.EvaluateBlock(ctx, request)
		cancel()
		if err == nil {
			var duration time.Duration
			resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
			lg.WithField("duration", duration).Debugf("request successful")
			resp.Metadata["imageHash"] = agent.config.ImageHash()
			agent.blockResults <- &scanner.BlockResult{
				AgentConfig: agent.config,
				Request:     request,
				Response:    resp,
			}
			lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
			continue
		}
		lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
		if agent.errCounter.TooManyErrs(err) {
			lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
			agent.Close()
			agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agent.config})
			return
		}
	}
}

func calculateResponseTime(startTime *time.Time) (timestamp string, latencyMs uint32, duration time.Duration) {
	now := time.Now().UTC()
	return now.Format(time.RFC3339), uint32(duration.Milliseconds()), duration
}

// ShouldProcessBlock tells if the agent should process block.
func (agent *Agent) ShouldProcessBlock(blockNumber string) bool {
	n, _ := strconv.ParseUint(blockNumber, 10, 64)
	var isAtLeastStartBlock bool
	if agent.config.StartBlock != nil {
		isAtLeastStartBlock = *agent.config.StartBlock >= n
	} else {
		isAtLeastStartBlock = true
	}

	var isAtMostStopBlock bool
	if agent.config.StopBlock != nil {
		isAtMostStopBlock = *agent.config.StopBlock <= n
	} else {
		isAtMostStopBlock = true
	}

	return isAtLeastStartBlock && isAtMostStopBlock
}
