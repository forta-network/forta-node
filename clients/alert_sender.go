package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/security"
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

type alertSender struct {
	ctx     context.Context
	cfg     AlertSenderConfig
	pClient protocol.PublisherNodeClient
}

type AlertSenderConfig struct {
	Key               *keystore.Key
	PublisherNodeAddr string
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

func NewAlertSender(ctx context.Context, cfg AlertSenderConfig) (*alertSender, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:8770", cfg.PublisherNodeAddr), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Errorf("could not reach %s within timeout", cfg.PublisherNodeAddr)
		return nil, err
	}
	pc := protocol.NewPublisherNodeClient(conn)
	return &alertSender{
		ctx:     ctx,
		cfg:     cfg,
		pClient: pc,
	}, nil
}
