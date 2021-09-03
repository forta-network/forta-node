package scanner

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/domain"
	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/store"
	"github.com/forta-network/forta-node/utils"
)

// BlockAnalyzerService reads TX info, calls agents, and emits results
type BlockAnalyzerService struct {
	queryNode protocol.QueryNodeClient
	cfg       BlockAnalyzerServiceConfig
	ctx       context.Context
}

type BlockAnalyzerServiceConfig struct {
	BlockChannel <-chan *domain.BlockEvent
	AlertSender  clients.AlertSender
	AgentPool    AgentPool
}

func (t *BlockAnalyzerService) calculateAlertID(result *BlockResult, f *protocol.Finding) string {
	idStr := strings.Join([]string{
		result.Request.Event.Network.ChainId,
		result.Request.Event.BlockHash,
		f.AlertId,
		f.Severity.String(),
		result.AgentConfig.ID}, "")
	return crypto.Keccak256Hash([]byte(idStr)).Hex()
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
	return &protocol.Alert{
		Id:        alertID,
		Finding:   f,
		Timestamp: ts.Format(store.AlertTimeFormat),
		Type:      protocol.AlertType_BLOCK,
		Agent:     result.AgentConfig.ToAgentInfo(),
		Tags: map[string]string{
			"agentImage":  result.AgentConfig.Image,
			"agentId":     result.AgentConfig.ID,
			"blockHash":   result.Request.Event.BlockHash,
			"blockNumber": blockNumber.String(),
			"chainId":     chainId.String(),
		},
	}, nil
}

func (t *BlockAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	//TODO: change this protocol when we know more about query-node delivery
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
				continue
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
