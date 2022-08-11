package scanner

import (
	"context"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/protocol/alerthash"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/metrics"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	ctx context.Context
	cfg TxAnalyzerServiceConfig

	lastInputActivity  health.TimeTracker
	lastOutputActivity health.TimeTracker
}

type TxAnalyzerServiceConfig struct {
	TxChannel   <-chan *domain.TransactionEvent
	AlertSender clients.AlertSender
	AgentPool   AgentPool
	MsgClient   clients.MessageClient
}

func (t *TxAnalyzerService) publishMetrics(result *TxResult) {
	m := metrics.GetTxMetrics(result.AgentConfig, result.Response, result.Timestamps)
	t.cfg.MsgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: m})
}

func (t *TxAnalyzerService) findingToAlert(result *TxResult, ts time.Time, f *protocol.Finding) (*protocol.Alert, error) {
	alertID := alerthash.ForTransactionAlert(
		&alerthash.Inputs{
			Transaction: result.Request.Event,
			Finding:     f,
			BotInfo: alerthash.BotInfo{
				BotImage: result.AgentConfig.Image,
				BotID:    result.AgentConfig.ID,
			},
		},
	)
	blockNumber, err := utils.HexToBigInt(result.Request.Event.Block.BlockNumber)
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
	if !f.Private && !result.Response.Private {
		alertType = protocol.AlertType_TRANSACTION
		tags["txHash"] = result.Request.Event.Transaction.Hash
		tags["blockHash"] = result.Request.Event.Block.BlockHash
		tags["blockNumber"] = blockNumber.String()
	}

	return &protocol.Alert{
		Id:         alertID,
		Finding:    f,
		Timestamp:  ts.Format(utils.AlertTimeFormat),
		Type:       alertType,
		Agent:      result.AgentConfig.ToAgentInfo(),
		Tags:       tags,
		Timestamps: result.Timestamps.ToMessage(),
	}, nil
}

func (t *TxAnalyzerService) Start() error {
	go func() {
		for result := range t.cfg.AgentPool.TxResults() {
			ts := time.Now().UTC()

			rt := &clients.AgentRoundTrip{
				AgentConfig:    result.AgentConfig,
				EvalTxRequest:  result.Request,
				EvalTxResponse: result.Response,
			}

			if len(result.Response.Findings) == 0 {
				if err := t.cfg.AlertSender.NotifyWithoutAlert(
					rt, result.Timestamps,
				); err != nil {
					log.WithError(err).Panic("failed to notify without alert")
				}
			}

			// TODO: validate finding returned is well-formed
			for _, f := range result.Response.Findings {
				alert, err := t.findingToAlert(result, ts, f)
				if err != nil {
					log.WithError(err).Error("failed to transform finding to alert")
					continue
				}
				if err := t.cfg.AlertSender.SignAlertAndNotify(
					rt, alert, result.Request.Event.Network.ChainId, result.Request.Event.Block.BlockNumber, result.Timestamps,
				); err != nil {
					log.WithError(err).Panic("failed to sign alert and notify")
				}
			}
			t.publishMetrics(result)

			t.lastOutputActivity.Set()
		}
	}()

	// Gear 1: loops over transactions and distributes to all agents
	go func() {
		// for each transaction
		for tx := range t.cfg.TxChannel {
			// convert to message
			msg, err := tx.ToMessage()
			if err != nil {
				log.WithError(err).Error("error converting tx event to message (skipping)")
				continue
			}

			// create a request
			requestId := uuid.Must(uuid.NewUUID())
			request := &protocol.EvaluateTxRequest{RequestId: requestId.String(), Event: msg}

			// forward to the pool
			t.cfg.AgentPool.SendEvaluateTxRequest(request)

			t.lastInputActivity.Set()
		}
	}()

	return nil
}

func (t *TxAnalyzerService) Stop() error {
	return nil
}

func (t *TxAnalyzerService) Name() string {
	return "tx-analyzer"
}

// Health implements the health.Reporter interface.
func (t *TxAnalyzerService) Health() health.Reports {
	return health.Reports{
		t.lastInputActivity.GetReport("event.input.time"),
		t.lastOutputActivity.GetReport("event.output.time"),
	}
}

func NewTxAnalyzerService(ctx context.Context, cfg TxAnalyzerServiceConfig) (*TxAnalyzerService, error) {
	return &TxAnalyzerService{
		cfg: cfg,
		ctx: ctx,
	}, nil
}
