package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/domain"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"
)

// BlockAnalyzerService reads TX info, calls agents, and emits results
type BlockAnalyzerService struct {
	queryNode protocol.QueryNodeClient
	cfg       BlockAnalyzerServiceConfig
	ctx       context.Context
	agents    []AnalyzerAgent
}

type BlockAnalyzerServiceConfig struct {
	BlockChannel <-chan *domain.BlockEvent
	AlertSender  clients.AlertSender
	AgentConfigs []config.AgentConfig
}

type evalBlockResp struct {
	request  *protocol.EvaluateBlockRequest
	agent    AnalyzerAgent
	response *protocol.EvaluateBlockResponse
}

// newAgentStream creates a agent transaction handler (sends and receives request)
func newAgentBlockStream(ctx context.Context, agent AnalyzerAgent, input <-chan *protocol.EvaluateBlockRequest, output chan<- *evalBlockResp) func() error {
	return func() error {
		for request := range input {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			resp, err := agent.client.EvaluateBlock(ctx, request)
			cancel()
			if err != nil {
				log.Error("error invoking agent", err)
				continue
			}
			resp.Metadata["imageHash"] = agent.config.ImageHash
			output <- &evalBlockResp{
				agent:    agent,
				response: resp,
				request:  request,
			}

			m := jsonpb.Marshaler{}
			resStr, err := m.MarshalToString(resp)
			if err != nil {
				log.Error("error marshaling response", err)
				continue
			}
			log.Debugf(resStr)
		}
		return nil
	}
}

func (t *BlockAnalyzerService) findingToAlert(resp *evalBlockResp, ts time.Time, f *protocol.Finding) (*protocol.Alert, error) {
	b, err := proto.Marshal(f)
	if err != nil {
		return nil, err
	}
	//TODO: come up with explicit way of producing ID
	alertID := base58.Encode(sha3.New256().Sum(b))
	return &protocol.Alert{
		Id:        alertID,
		Finding:   f,
		Timestamp: ts.Format(store.AlertTimeFormat),
		Type:      protocol.AlertType_BLOCK,
		Agent: &protocol.AgentInfo{
			Name:      resp.agent.config.Name,
			Image:     resp.agent.config.Image,
			ImageHash: resp.agent.config.ImageHash,
		},
		Tags: map[string]string{
			"blockHash":   resp.request.Event.BlockHash,
			"blockNumber": resp.request.Event.BlockNumber,
		},
	}, nil
}

func (t *BlockAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	//TODO: change this protocol when we know more about query-node delivery
	// Gear 3: receive result from agent
	output := make(chan *evalBlockResp, 100)
	grp.Go(func() error {
		for resp := range output {
			ts := time.Now().UTC()
			for _, f := range resp.response.Findings {
				alert, err := t.findingToAlert(resp, ts, f)
				if err != nil {
					return err
				}
				if err := t.cfg.AlertSender.SignAndNotify(alert); err != nil {
					return err
				}
			}
		}
		return nil
	})

	// Gear 2: set of agents pulling from their own individual channels
	var agentChannels []chan *protocol.EvaluateBlockRequest
	for _, agt := range t.agents {
		agent := agt
		input := make(chan *protocol.EvaluateBlockRequest, 100)
		agentChannels = append(agentChannels, input)
		grp.Go(newAgentBlockStream(ctx, agent, input, output))
	}

	// Gear 1: loops over transactions and distributes to all agents
	grp.Go(func() error {
		defer func() {
			for _, agtCh := range agentChannels {
				close(agtCh)
			}
			close(output)
		}()

		// for each block
		for block := range t.cfg.BlockChannel {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// convert to message
			blockEvt, err := block.ToMessage()
			if err != nil {
				log.Error("error converting block event to message (skipping)", err)
				continue
			}

			// create a request
			requestId := uuid.Must(uuid.NewUUID())
			request := &protocol.EvaluateBlockRequest{RequestId: requestId.String(), Event: blockEvt}

			// forward to each agent channel
			for _, agtCh := range agentChannels {
				agtCh <- request
			}
		}
		return nil
	})

	return grp.Wait()
}

func (t *BlockAnalyzerService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *BlockAnalyzerService) Name() string {
	return "BlockAnalyzerService"
}

func NewBlockAnalyzerService(ctx context.Context, cfg BlockAnalyzerServiceConfig) (*BlockAnalyzerService, error) {
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

	return &BlockAnalyzerService{
		cfg:    cfg,
		ctx:    ctx,
		agents: agents,
	}, nil
}
