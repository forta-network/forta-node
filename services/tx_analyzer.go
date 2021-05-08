package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"OpenZeppelin/zephyr-node/feeds"
	"OpenZeppelin/zephyr-node/protocol"
)

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	queryNode protocol.QueryNodeClient
	cfg       TxAnalyzerServiceConfig
	ctx       context.Context
	agents    []protocol.AgentClient
}

type TxAnalyzerServiceConfig struct {
	TxChannel      <-chan *feeds.TransactionEvent
	AgentAddresses []string
	QueryNodeAddr  string
}

// newAgentStream creates a agent transaction handler (sends and receives request)
func newAgentStream(ctx context.Context, agent protocol.AgentClient, input <-chan *protocol.EvaluateRequest, output chan<- *protocol.EvaluateResponse) func() error {
	return func() error {
		for request := range input {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			resp, err := agent.Evaluate(ctx, request)
			cancel()
			if err != nil {
				log.Error("error invoking agent", err)
				continue
			}
			output <- resp

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
	output := make(chan *protocol.EvaluateResponse, 100)
	grp.Go(func() error {
		for resp := range output {
			ts := time.Now()
			for _, f := range resp.Findings {
				alert := &protocol.Alert{
					Id:        "test",
					Finding:   f,
					Timestamp: ts.String(),
					Metadata:  nil,
				}
				_, err := t.queryNode.Notify(ctx, &protocol.NotifyRequest{
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
	var clients []protocol.AgentClient
	for _, addr := range cfg.AgentAddresses {
		conn, err := grpc.Dial(fmt.Sprintf("%s:50051", addr), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect to %s, %v", addr, err)
		}
		clients = append(clients, protocol.NewAgentClient(conn))
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:8770", cfg.QueryNodeAddr), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	qn := protocol.NewQueryNodeClient(conn)

	return &TxAnalyzerService{
		cfg:       cfg,
		ctx:       ctx,
		agents:    clients,
		queryNode: qn,
	}, nil
}
