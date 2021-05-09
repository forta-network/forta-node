package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/feeds"
	"OpenZeppelin/fortify-node/protocol"
)

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	queryNode protocol.QueryNodeClient
	cfg       TxAnalyzerServiceConfig
	ctx       context.Context
	agents    []AnalyzerAgent
}

type AnalyzerAgent struct {
	config config.AgentConfig
	client protocol.AgentClient
}

type TxAnalyzerServiceConfig struct {
	TxChannel     <-chan *feeds.TransactionEvent
	AgentConfigs  []config.AgentConfig
	QueryNodeAddr string
}

type responseWrapper struct {
	agent    AnalyzerAgent
	response *protocol.EvaluateResponse
}

// newAgentStream creates a agent transaction handler (sends and receives request)
func newAgentStream(ctx context.Context, agent AnalyzerAgent, input <-chan *protocol.EvaluateRequest, output chan<- *responseWrapper) func() error {
	return func() error {
		for request := range input {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			resp, err := agent.client.Evaluate(ctx, request)
			cancel()
			if err != nil {
				log.Error("error invoking agent", err)
				continue
			}
			resp.Metadata["image"] = agent.config.Image
			output <- &responseWrapper{
				agent:    agent,
				response: resp,
			}

			//TODO: this print logic is just for test purposes
			m := jsonpb.Marshaler{}
			resStr, err := m.MarshalToString(resp)
			if err != nil {
				log.Error("error marshaling response", err)
				continue
			}
			log.Infof(resStr)
		}
		return nil
	}
}

func (t *TxAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	//TODO: change this protocol when we know more about query-node delivery
	// Gear 3: handles responses (sends to query node)
	output := make(chan *responseWrapper, 100)
	grp.Go(func() error {
		for resp := range output {
			ts := time.Now()
			for _, f := range resp.response.Findings {
				b, err := proto.Marshal(f)
				if err != nil {
					return err
				}
				alertID := hexutil.Encode(sha3.New256().Sum(b))
				alert := &protocol.Alert{
					Id:        alertID,
					Finding:   f,
					Timestamp: ts.String(),
					Metadata: map[string]string{
						"agent-name":  resp.agent.config.Name,
						"agent-image": resp.agent.config.Image,
					},
				}
				_, err = t.queryNode.Notify(ctx, &protocol.NotifyRequest{
					Alert: alert,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	// Gear 2: set of agents pulling from their own individual channels
	var agentChannels []chan *protocol.EvaluateRequest
	for _, agt := range t.agents {
		agent := agt
		input := make(chan *protocol.EvaluateRequest, 100)
		agentChannels = append(agentChannels, input)
		grp.Go(newAgentStream(ctx, agent, input, output))
	}

	// Gear 1: loops over transactions and distributes to all agents
	grp.Go(func() error {
		defer func() {
			for _, agtCh := range agentChannels {
				close(agtCh)
			}
			close(output)
		}()

		// for each transaction
		for tx := range t.cfg.TxChannel {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// convert to message
			msg, err := tx.ToMessage()
			if err != nil {
				log.Error("error converting tx event to message", err)
				continue
			}

			// create a request
			requestId := uuid.Must(uuid.NewUUID())
			request := &protocol.EvaluateRequest{RequestId: requestId.String(), Event: msg}

			// forward to each agent channel
			for _, agtCh := range agentChannels {
				agtCh <- request
			}
		}
		return nil
	})

	return grp.Wait()
}

func (t *TxAnalyzerService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxAnalyzerService) Name() string {
	return "TxAnalyzerService"
}

func NewTxAnalyzerService(ctx context.Context, cfg TxAnalyzerServiceConfig) (*TxAnalyzerService, error) {
	var agents []AnalyzerAgent
	for _, agt := range cfg.AgentConfigs {
		conn, err := grpc.Dial(fmt.Sprintf("%s:50051", agt.ContainerName()), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect to %s, %v", agt.ContainerName(), err)
		}
		client := protocol.NewAgentClient(conn)
		agents = append(agents, AnalyzerAgent{
			config: agt,
			client: client,
		})
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:8770", cfg.QueryNodeAddr), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	qn := protocol.NewQueryNodeClient(conn)

	return &TxAnalyzerService{
		cfg:       cfg,
		ctx:       ctx,
		agents:    agents,
		queryNode: qn,
	}, nil
}
