package clients

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/store"
	log "github.com/sirupsen/logrus"
)

// AgentRoundTrip contains
type AgentRoundTrip struct {
	AgentConfig             config.AgentConfig
	EvalBlockRequest        *protocol.EvaluateBlockRequest
	EvalBlockResponse       *protocol.EvaluateBlockResponse
	EvalTxRequest           *protocol.EvaluateTxRequest
	EvalTxResponse          *protocol.EvaluateTxResponse
	EvalAlertRequest        *protocol.EvaluateAlertRequest
	EvalAlertResponse       *protocol.EvaluateAlertResponse
	EvalHealthCheckRequest  *protocol.HealthCheckRequest
	EvalHealthCheckResponse *protocol.HealthCheckResponse
}

type AlertSender interface {
	SignAlertAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string, ts *domain.TrackingTimestamps) error
	NotifyWithoutAlert(rt *AgentRoundTrip, ts *domain.TrackingTimestamps) error
}

// PublishClient implements the interface for a notify
type PublishClient interface {
	Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error)
}

type alertSender struct {
	ctx     context.Context
	cfg     AlertSenderConfig
	pClient PublishClient
}

type AlertSenderConfig struct {
	Key *keystore.Key
	DS  store.DeduplicationStore
}

func (a *alertSender) SignAlertAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string, ts *domain.TrackingTimestamps) error {
	logger := log.WithFields(log.Fields{
		"alert": alert.Id,
	})

	// only if configured (for local mode redundancy)
	if a.cfg.DS != nil {
		if isFirst, err := a.cfg.DS.IsFirst(alert.Id); err != nil {
			// treating errors as non-duplicate
			logger.WithError(err).Error("error checking for duplicate (assuming is first)")
		} else if !isFirst {
			logger.Debug("duplicate alert (ignoring)")
			return nil
		}
	}
	alert.Scanner = &protocol.ScannerInfo{
		Address: a.cfg.Key.Address.Hex(),
	}
	signedAlert, err := security.SignAlert(a.cfg.Key, alert)
	if err != nil {
		logger.Errorf("could not sign alert (id=%s), skipping", alert.Id)
		return err
	}
	signedAlert.ChainId = chainID
	signedAlert.BlockNumber = blockNumber
	_, err = a.pClient.Notify(
		a.ctx, &protocol.NotifyRequest{
			SignedAlert:       signedAlert,
			EvalBlockRequest:  rt.EvalBlockRequest,
			EvalBlockResponse: rt.EvalBlockResponse,
			EvalTxRequest:     rt.EvalTxRequest,
			EvalTxResponse:    rt.EvalTxResponse,
			EvalAlertRequest:  rt.EvalAlertRequest,
			EvalAlertResponse: rt.EvalAlertResponse,
			AgentInfo:         rt.AgentConfig.ToAgentInfo(),
			Timestamps:        ts.ToMessage(),
		},
	)
	return err
}

func (a *alertSender) NotifyWithoutAlert(rt *AgentRoundTrip, ts *domain.TrackingTimestamps) error {
	_, err := a.pClient.Notify(
		a.ctx, &protocol.NotifyRequest{
			EvalBlockRequest:  rt.EvalBlockRequest,
			EvalBlockResponse: rt.EvalBlockResponse,
			EvalAlertRequest:  rt.EvalAlertRequest,
			EvalAlertResponse: rt.EvalAlertResponse,
			EvalTxRequest:     rt.EvalTxRequest,
			EvalTxResponse:    rt.EvalTxResponse,
			AgentInfo:         rt.AgentConfig.ToAgentInfo(),
			Timestamps:        ts.ToMessage(),
		},
	)
	return err
}

func NewAlertSender(ctx context.Context, publisher PublishClient, cfg AlertSenderConfig) (AlertSender, error) {
	return &alertSender{
		ctx:     ctx,
		cfg:     cfg,
		pClient: publisher,
	}, nil
}
