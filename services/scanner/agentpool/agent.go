package agentpool

import (
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/protocol"
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Constants
const (
	DefaultBufferSize = 100
)

// Agent receives blocks and transactions, and produces results.
type Agent struct {
	config config.AgentConfig

	evalTxCh     chan *protocol.EvaluateTxRequest
	txResults    chan<- *TxResult
	evalBlockCh  chan *protocol.EvaluateBlockRequest
	blockResults chan<- *BlockResult

	client protocol.AgentClient
	conn   *grpc.ClientConn
	ready  bool
}

// NewAgent creates a new agent.
func NewAgent(agentCfg config.AgentConfig, txResults chan<- *TxResult, blockResults chan<- *BlockResult) *Agent {
	return &Agent{
		config:       agentCfg,
		evalTxCh:     make(chan *protocol.EvaluateTxRequest, DefaultBufferSize),
		txResults:    txResults,
		evalBlockCh:  make(chan *protocol.EvaluateBlockRequest, DefaultBufferSize),
		blockResults: blockResults,
	}
}

// Config returns the agent config.
func (agent *Agent) Config() config.AgentConfig {
	return agent.config
}

// Close implements io.Closer.
func (agent *Agent) Close() error {
	close(agent.evalTxCh)
	close(agent.evalBlockCh)
	if agent.conn != nil {
		return agent.conn.Close()
	}
	return nil
}

func (agent *Agent) connect() error {
	cfg := agent.config
	var (
		conn *grpc.ClientConn
		err  error
	)
	for i := 0; i < 10; i++ {
		conn, err = grpc.Dial(fmt.Sprintf("%s:%s", cfg.ContainerName(), cfg.GrpcPort()), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
		if err == nil {
			break
		}
		err = fmt.Errorf("failed to connect to agent '%s': %v", cfg.ContainerName(), err)
		log.Debug(err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Panic(err)
	}
	agent.client = protocol.NewAgentClient(conn)
	agent.conn = conn
	log.Debugf("connected to agent: %s", cfg.ContainerName())
	return nil
}

func (agent *Agent) processTransactions() {
	for request := range agent.evalTxCh {
		processingState.waitIfPaused()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := agent.client.EvaluateTx(ctx, request)
		cancel()
		if err != nil {
			log.Error("error invoking agent", err)
			continue
		}
		resp.Metadata["imageHash"] = agent.config.ImageHash
		agent.txResults <- &TxResult{
			AgentConfig: agent.config,
			Request:     request,
			Response:    resp,
		}
	}
}

func (agent *Agent) processBlocks() {
	for request := range agent.evalBlockCh {
		processingState.waitIfPaused()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := agent.client.EvaluateBlock(ctx, request)
		cancel()
		if err != nil {
			log.Error("error invoking agent", err)
			continue
		}
		resp.Metadata["imageHash"] = agent.config.ImageHash
		agent.blockResults <- &BlockResult{
			AgentConfig: agent.config,
			Request:     request,
			Response:    resp,
		}
	}
}

func (agent *Agent) shouldProcessBlock(blockNumber string) bool {
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
