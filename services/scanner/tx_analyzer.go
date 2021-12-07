package scanner

import (
	"context"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/metrics"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/domain"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	cfg TxAnalyzerServiceConfig
	ctx context.Context
}

type TxAnalyzerServiceConfig struct {
	TxChannel   <-chan *domain.TransactionEvent
	AlertSender clients.AlertSender
	AgentPool   AgentPool
	MsgClient   clients.MessageClient
}

func (t *TxAnalyzerService) calculateAlertID(result *TxResult, f *protocol.Finding) string {
	idStr := strings.Join([]string{
		result.Request.Event.Network.ChainId,
		result.Request.Event.Transaction.Hash,
		f.Name,
		f.Description,
		f.Protocol,
		f.Type.String(),
		f.AlertId,
		f.Severity.String(),
		result.AgentConfig.Image,
		result.AgentConfig.ID,
		strings.Join(utils.MapKeys(result.Request.Event.Addresses), "")}, "")
	return crypto.Keccak256Hash([]byte(idStr)).Hex()
}

func (t *TxAnalyzerService) publishMetrics(result *TxResult) {
	m := metrics.GetTxMetrics(result.AgentConfig, result.Response)
	t.cfg.MsgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: m})
}

func (t *TxAnalyzerService) findingToAlert(result *TxResult, ts time.Time, f *protocol.Finding) (*protocol.Alert, error) {
	alertID := t.calculateAlertID(result, f)
	blockNumber, err := utils.HexToBigInt(result.Request.Event.Block.BlockNumber)
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
		Timestamp: ts.Format(utils.AlertTimeFormat),
		Type:      protocol.AlertType_TRANSACTION,
		Agent:     result.AgentConfig.ToAgentInfo(),
		Tags: map[string]string{
			"agentImage":  result.AgentConfig.Image,
			"agentId":     result.AgentConfig.ID,
			"chainId":     chainId.String(),
			"blockHash":   result.Request.Event.Block.BlockHash,
			"blockNumber": blockNumber.String(),
			"txHash":      result.Request.Event.Transaction.Hash,
		},
	}, nil
}

func (t *TxAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	grp.Go(func() error {
		for result := range t.cfg.AgentPool.TxResults() {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			ts := time.Now().UTC()

			rt := &clients.AgentRoundTrip{
				AgentConfig:    result.AgentConfig,
				EvalTxRequest:  result.Request,
				EvalTxResponse: result.Response,
			}

			if len(result.Response.Findings) == 0 {
				if err := t.cfg.AlertSender.NotifyWithoutAlert(
					rt, result.Request.Event.Network.ChainId, result.Request.Event.Block.BlockNumber,
				); err != nil {
					log.WithError(err).Error("failed to notify without alert")
					return err
				}
			}

			//TODO: validate finding returned is well-formed
			for _, f := range result.Response.Findings {
				alert, err := t.findingToAlert(result, ts, f)
				if err != nil {
					return err
				}
				if err := t.cfg.AlertSender.SignAlertAndNotify(
					rt, alert, result.Request.Event.Network.ChainId, result.Request.Event.Block.BlockNumber,
				); err != nil {
					log.WithError(err).Error("failed to sign alert and notify")
					return err
				}
			}
			t.publishMetrics(result)
		}
		return nil
	})

	// Gear 1: loops over transactions and distributes to all agents
	grp.Go(func() error {
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

			// forward to the pool
			t.cfg.AgentPool.SendEvaluateTxRequest(request)
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
	return &TxAnalyzerService{
		cfg: cfg,
		ctx: ctx,
	}, nil
}
