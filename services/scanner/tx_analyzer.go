package scanner

import (
	"context"
	"fmt"
	"time"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/domain"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"

	"github.com/btcsuite/btcutil/base58"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
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

func (t *TxAnalyzerService) calculateAlertID(result *TxResult, f *protocol.Finding) (string, error) {
	findingBytes, err := proto.Marshal(f)
	if err != nil {
		return "", err
	}
	idStr := fmt.Sprintf("%s%s%s", result.Request.Event.Network.ChainId, result.Request.Event.Transaction.Hash, string(findingBytes))
	return base58.Encode(sha3.New256().Sum([]byte(idStr))), nil
}

func (t *TxAnalyzerService) findingToAlert(result *TxResult, ts time.Time, f *protocol.Finding) (*protocol.Alert, error) {
	alertID, err := t.calculateAlertID(result, f)
	if err != nil {
		return nil, err
	}
	return &protocol.Alert{
		Id:        alertID,
		Finding:   f,
		Timestamp: ts.Format(store.AlertTimeFormat),
		Type:      protocol.AlertType_TRANSACTION,
		Agent: &protocol.AgentInfo{
			Name:      result.AgentConfig.Name,
			Image:     result.AgentConfig.Image,
			ImageHash: result.AgentConfig.ImageHash,
		},
		Tags: map[string]string{
			"chainId":     result.Request.Event.Network.ChainId,
			"blockHash":   result.Request.Event.Block.BlockHash,
			"blockNumber": result.Request.Event.Block.BlockNumber,
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
				if err := t.cfg.AlertSender.SignAndNotify(alert); err != nil {
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
