package scanner

import (
	"context"
	"strings"
	"time"

	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/metrics"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/domain"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"
)

// BlockAnalyzerService reads TX info, calls agents, and emits results
type BlockAnalyzerService struct {
	publisherNode protocol.PublisherNodeClient
	cfg           BlockAnalyzerServiceConfig
	ctx           context.Context
}

type BlockAnalyzerServiceConfig struct {
	BlockChannel <-chan *domain.BlockEvent
	AlertSender  clients.AlertSender
	AgentPool    AgentPool
	MsgClient    clients.MessageClient
}

// WARNING, this must be deterministic (any maps must be converted to sorted lists)
func (t *BlockAnalyzerService) calculateAlertID(result *BlockResult, f *protocol.Finding) string {
	idStr := strings.Join([]string{
		result.Request.Event.Network.ChainId,
		result.Request.Event.BlockHash,
		f.AlertId,
		f.Name,
		f.Description,
		f.Protocol,
		f.Type.String(),
		f.Severity.String(),
		result.AgentConfig.Image,
		result.AgentConfig.ID}, "")
	return crypto.Keccak256Hash([]byte(idStr)).Hex()
}

func (t *BlockAnalyzerService) publishMetrics(result *BlockResult) {
	m := metrics.GetBlockMetrics(result.AgentConfig, result.Response)
	t.cfg.MsgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: m})
}

func (t *BlockAnalyzerService) findingToAlert(result *BlockResult, ts time.Time, f *protocol.Finding) (*protocol.Alert, error) {
	alertID := t.calculateAlertID(result, f)
	blockNumber, err := utils.HexToBigInt(result.Request.Event.BlockNumber)
	if err != nil {
		return nil, err
	}
	chainId, err := utils.HexToBigInt(result.Request.Event.Network.ChainId)
	if err != nil {
		return nil, err
	}
	tags := map[string]string{
		"agentImage": result.AgentConfig.Image,
		"agentId":    result.AgentConfig.ID,
		"chainId":    chainId.String(),
	}

	alertType := protocol.AlertType_PRIVATE
	if !result.Response.Private {
		alertType = protocol.AlertType_BLOCK
		tags["blockHash"] = result.Request.Event.BlockHash
		tags["blockNumber"] = blockNumber.String()
	}
	return &protocol.Alert{
		Id:        alertID,
		Finding:   f,
		Timestamp: ts.Format(utils.AlertTimeFormat),
		Type:      alertType,
		Agent:     result.AgentConfig.ToAgentInfo(),
		Tags:      tags,
	}, nil
}

func (t *BlockAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	// Gear 2: receive result from agent
	grp.Go(func() error {
		for result := range t.cfg.AgentPool.BlockResults() {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			ts := time.Now().UTC()

			m := jsonpb.Marshaler{}
			resStr, err := m.MarshalToString(result.Response)
			if err != nil {
				log.Error("error marshaling response", err)
				continue
			}
			log.Debugf(resStr)

			rt := &clients.AgentRoundTrip{
				AgentConfig:       result.AgentConfig,
				EvalBlockRequest:  result.Request,
				EvalBlockResponse: result.Response,
			}

			if len(result.Response.Findings) == 0 {
				if err := t.cfg.AlertSender.NotifyWithoutAlert(
					rt, result.Request.Event.Network.ChainId, result.Request.Event.BlockNumber,
				); err != nil {
					log.WithError(err).Error("failed to notify without alert")
					return err
				}
			}

			for _, f := range result.Response.Findings {
				alert, err := t.findingToAlert(result, ts, f)
				if err != nil {
					return err
				}
				if err := t.cfg.AlertSender.SignAlertAndNotify(
					rt, alert, result.Request.Event.Network.ChainId, result.Request.Event.BlockNumber,
				); err != nil {
					log.WithError(err).Error("failed sign alert and notify")
					return err
				}
			}
			t.publishMetrics(result)
		}
		return nil
	})

	// Gear 1: loops over blocks and distributes to all agents
	grp.Go(func() error {
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

			// forward to the pool
			t.cfg.AgentPool.SendEvaluateBlockRequest(request)
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
	return &BlockAnalyzerService{
		cfg: cfg,
		ctx: ctx,
	}, nil
}
