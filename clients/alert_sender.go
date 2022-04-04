package clients

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-node/config"
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

// PublishGrpcClient wraps the grpc client
type PublishGrpcClient struct {
	c protocol.PublisherNodeClient
}

func (pgc *PublishGrpcClient) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	return pgc.Notify(ctx, req)
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

func NewLocalAlertSender(ctx context.Context, publisher PublishClient, cfg AlertSenderConfig) (*alertSender, error) {
	return &alertSender{
		ctx:     ctx,
		cfg:     cfg,
		pClient: publisher,
	}, nil
}

func NewGRPCAlertSender(ctx context.Context, cfg AlertSenderConfig) (*alertSender, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:8770", cfg.PublisherNodeAddr), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Errorf("could not reach %s within timeout", cfg.PublisherNodeAddr)
		return nil, err
	}
	pc := protocol.NewPublisherNodeClient(conn)
	return &alertSender{
		ctx:     ctx,
		cfg:     cfg,
		pClient: &PublishGrpcClient{c: pc},
	}, nil
}
