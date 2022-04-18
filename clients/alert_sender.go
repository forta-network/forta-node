package clients

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/keystore"
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
}

type AlertSender interface {
	SignAlertAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string) error
	NotifyWithoutAlert(rt *AgentRoundTrip, chainID, blockNumber string) error
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

func (a *alertSender) SignAlertAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string) error {
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
	})
	return err
}

func (a *alertSender) NotifyWithoutAlert(rt *AgentRoundTrip, chainID, blockNumber string) error {
	_, err := a.pClient.Notify(a.ctx, &protocol.NotifyRequest{
		EvalBlockRequest:  rt.EvalBlockRequest,
		EvalBlockResponse: rt.EvalBlockResponse,
		EvalTxRequest:     rt.EvalTxRequest,
		EvalTxResponse:    rt.EvalTxResponse,
		AgentInfo:         rt.AgentConfig.ToAgentInfo(),
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
