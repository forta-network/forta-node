package clients

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"fortify-node/protocol"
	"fortify-node/security"
)

type AlertSender interface {
	SignAndNotify(alert *protocol.Alert) error
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

func (a *alertSender) SignAndNotify(alert *protocol.Alert) error {
	alert.Scanner = &protocol.ScannerInfo{
		Address: a.cfg.Key.Address.Hex(),
	}
	signedAlert, err := security.SignAlert(a.cfg.Key, alert)
	if err != nil {
		log.Errorf("could not sign alert (id=%s), skipping", alert.Id)
		return err
	}
	_, err = a.qClient.Notify(a.ctx, &protocol.NotifyRequest{
		SignedAlert: signedAlert,
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
