package clients

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/security"
)

// AgentRoundTrip contains
type AgentRoundTrip struct {
	EvalBlockRequest  *protocol.EvaluateBlockRequest
	EvalBlockResponse *protocol.EvaluateBlockResponse
	EvalTxRequest     *protocol.EvaluateTxRequest
	EvalTxResponse    *protocol.EvaluateTxResponse
}

type AlertSender interface {
	SignAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string) error
}

type alertSender struct {
	ctx     context.Context
	cfg     AlertSenderConfig
	qClient protocol.QueryNodeClient
}

type AlertSenderConfig struct {
	Key           *keystore.Key
	QueryNodeAddr string
}

func (a *alertSender) SignAndNotify(rt *AgentRoundTrip, alert *protocol.Alert, chainID, blockNumber string) error {
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
	_, err = a.qClient.Notify(a.ctx, &protocol.NotifyRequest{
		SignedAlert:       signedAlert,
		EvalBlockRequest:  rt.EvalBlockRequest,
		EvalBlockResponse: rt.EvalBlockResponse,
		EvalTxRequest:     rt.EvalTxRequest,
		EvalTxResponse:    rt.EvalTxResponse,
	})
	return err
}

func NewAlertSender(ctx context.Context, cfg AlertSenderConfig) (*alertSender, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:8770", cfg.QueryNodeAddr), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	qc := protocol.NewQueryNodeClient(conn)
	return &alertSender{
		ctx:     ctx,
		cfg:     cfg,
		qClient: qc,
	}, nil
}
