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

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	cfg    TxAnalyzerServiceConfig
	ctx    context.Context
	agents []AnalyzerAgent
}

type TxAnalyzerServiceConfig struct {
	TxChannel    <-chan *domain.TransactionEvent
	AlertSender  clients.AlertSender
	AgentConfigs []config.AgentConfig
}

type evalTxResp struct {
	request  *protocol.EvaluateTxRequest
	agent    AnalyzerAgent
	response *protocol.EvaluateTxResponse
}

// newAgentStream creates a agent transaction handler (sends and receives request)
func newAgentTxStream(ctx context.Context, agent AnalyzerAgent, input <-chan *protocol.EvaluateTxRequest, output chan<- *evalTxResp) func() error {
	return func() error {
		for request := range input {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			resp, err := agent.client.EvaluateTx(ctx, request)
			cancel()
			if err != nil {
				log.Error("error invoking agent", err)
				continue
			}
			resp.Metadata["imageHash"] = agent.config.ImageHash
			output <- &evalTxResp{
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

func (t *TxAnalyzerService) calculateAlertID(resp *evalTxResp, f *protocol.Finding) (string, error) {
	findingBytes, err := proto.Marshal(f)
	if err != nil {
		return "", err
	}
	idStr := fmt.Sprintf("%s%s%s", resp.request.Event.Network.ChainId, resp.request.Event.Transaction.Hash, string(findingBytes))
	return base58.Encode(sha3.New256().Sum([]byte(idStr))), nil
}

func (t *TxAnalyzerService) findingToAlert(resp *evalTxResp, ts time.Time, f *protocol.Finding) (*protocol.Alert, error) {
	alertID, err := t.calculateAlertID(resp, f)
	if err != nil {
		return nil, err
	}
	return &protocol.Alert{
		Id:        alertID,
		Finding:   f,
		Timestamp: ts.Format(store.AlertTimeFormat),
		Type:      protocol.AlertType_TRANSACTION,
		Agent: &protocol.AgentInfo{
			Name:      resp.agent.config.Name,
			Image:     resp.agent.config.Image,
			ImageHash: resp.agent.config.ImageHash,
		},
		Tags: map[string]string{
			"chainId":     resp.request.Event.Network.ChainId,
			"blockHash":   resp.request.Event.Receipt.BlockHash,
			"blockNumber": resp.request.Event.Receipt.BlockNumber,
			"txHash":      resp.request.Event.Receipt.TransactionHash,
		},
	}, nil
}

func (t *TxAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	//TODO: change this protocol when we know more about query-node delivery
	// Gear 3: receive result from agent
	output := make(chan *evalTxResp, 100)
	grp.Go(func() error {
		for resp := range output {
			ts := time.Now().UTC()
			//TODO: validate finding returned is well-formed
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
	var agentChannels []chan *protocol.EvaluateTxRequest
	for _, agt := range t.agents {
		agent := agt
		input := make(chan *protocol.EvaluateTxRequest, 100)
		agentChannels = append(agentChannels, input)
		grp.Go(newAgentTxStream(ctx, agent, input, output))
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
				log.Error("error converting tx event to message (skipping)", err)
				continue
			}

			// create a request
			requestId := uuid.Must(uuid.NewUUID())
			request := &protocol.EvaluateTxRequest{RequestId: requestId.String(), Event: msg}

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
	return &TxAnalyzerService{
		cfg:    cfg,
		ctx:    ctx,
		agents: agents,
	}, nil
}
