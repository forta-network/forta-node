package scanner

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/metrics"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/domain"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

// WARNING, this must be deterministic (any maps must be converted to sorted lists)
func (t *TxAnalyzerService) calculateAlertID(result *TxResult, f *protocol.Finding) string {
	addrs := utils.MapKeys(result.Request.Event.Addresses)
	sort.Strings(addrs)
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
		strings.Join(addrs, "")}, "")
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
					rt, result.Request.Event.Network.ChainId, result.Request.Event.Block.BlockNumber,
				); err != nil {
					log.WithError(err).Panic("failed to notify without alert")
				}
			}

			//TODO: validate finding returned is well-formed
			for _, f := range result.Response.Findings {
				alert, err := t.findingToAlert(result, ts, f)
				if err != nil {
					log.WithError(err).Error("failed to transform finding to alert")
					continue
				}
				if err := t.cfg.AlertSender.SignAlertAndNotify(
					rt, alert, result.Request.Event.Network.ChainId, result.Request.Event.Block.BlockNumber,
				); err != nil {
					log.WithError(err).Panic("failed to sign alert and notify")
				}
			}
			t.publishMetrics(result)
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
		}
	}()

	return nil
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
