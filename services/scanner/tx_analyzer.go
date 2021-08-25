package scanner

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/domain"
	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/store"
	"github.com/forta-network/forta-node/utils"

	"github.com/golang/protobuf/jsonpb"
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
}

func (t *TxAnalyzerService) calculateAlertID(result *TxResult, f *protocol.Finding) string {
	idStr := strings.Join([]string{
		result.Request.Event.Network.ChainId,
		result.Request.Event.Transaction.Hash,
		f.AlertId,
		f.Severity.String(),
		result.AgentConfig.Image}, "")
	return crypto.Keccak256Hash([]byte(idStr)).Hex()
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
		Timestamp: ts.Format(store.AlertTimeFormat),
		Type:      protocol.AlertType_TRANSACTION,
		Agent: &protocol.AgentInfo{
			Name:      result.AgentConfig.ID,
			Image:     result.AgentConfig.Image,
			ImageHash: result.AgentConfig.ImageHash(),
		},
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

	//TODO: change this protocol when we know more about query-node delivery
	// Gear 2: receive result from agent
	grp.Go(func() error {
		for result := range t.cfg.AgentPool.TxResults() {
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

			//TODO: validate finding returned is well-formed
			for _, f := range result.Response.Findings {
				alert, err := t.findingToAlert(result, ts, f)
				if err != nil {
					return err
				}
				if err := t.cfg.AlertSender.SignAlertAndNotify(
					&clients.AgentRoundTrip{
						EvalTxRequest:  result.Request,
						EvalTxResponse: result.Response,
					},
					alert, result.Request.Event.Network.ChainId, result.Request.Event.Block.BlockNumber,
				); err != nil {
					return err
				}
			}
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
