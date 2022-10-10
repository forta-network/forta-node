package clients

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

// AgentRoundTrip contains
type AgentRoundTrip struct {
	AgentConfig       config.AgentConfig
	EvalBlockRequest  *protocol.EvaluateBlockRequest
	EvalBlockResponse *protocol.EvaluateBlockResponse
	EvalTxRequest     *protocol.EvaluateTxRequest
	EvalTxResponse    *protocol.EvaluateTxResponse
	EvalAlertRequest  *protocol.EvaluateAlertRequest
	EvalAlertResponse *protocol.EvaluateAlertResponse
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
}

func (a *alertSender) SignAlertAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string, ts *domain.TrackingTimestamps) error {
	alert.Scanner = &protocol.ScannerInfo{
		Address: a.cfg.Key.Address.Hex(),
	}
	signedAlert, err := security.SignAlert(a.cfg.Key, alert)
	if err != nil {
		log.Errorf("could not sign alert (id=%s), skipping", alert.Id)
		return err
	}
	signedAlert.ChainId = chainID
	signedAlert.BlockNumber = blockNumber
	_, err = a.pClient.Notify(a.ctx, &protocol.NotifyRequest{
		SignedAlert:       signedAlert,
		EvalBlockRequest:  rt.EvalBlockRequest,
		EvalBlockResponse: rt.EvalBlockResponse,
		EvalTxRequest:     rt.EvalTxRequest,
		EvalTxResponse:    rt.EvalTxResponse,
		AgentInfo:         rt.AgentConfig.ToAgentInfo(),
		Timestamps:        ts.ToMessage(),
	})
	return err
}

func (a *alertSender) NotifyWithoutAlert(rt *AgentRoundTrip, ts *domain.TrackingTimestamps) error {
	_, err := a.pClient.Notify(a.ctx, &protocol.NotifyRequest{
		EvalBlockRequest:  rt.EvalBlockRequest,
		EvalBlockResponse: rt.EvalBlockResponse,
		EvalTxRequest:     rt.EvalTxRequest,
		EvalTxResponse:    rt.EvalTxResponse,
		AgentInfo:         rt.AgentConfig.ToAgentInfo(),
		Timestamps:        ts.ToMessage(),
	})
	return err
}

func NewAlertSender(ctx context.Context, publisher PublishClient, cfg AlertSenderConfig) (*alertSender, error) {
	return &alertSender{
		ctx:     ctx,
		cfg:     cfg,
		pClient: publisher,
	}, nil
}
