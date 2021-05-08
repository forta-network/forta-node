package services

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"OpenZeppelin/zephyr-node/protocol"
	"OpenZeppelin/zephyr-node/store"
)

// AlertListener allows retrieval of alerts from the database
type AlertListener struct {
	protocol.UnimplementedQueryNodeServer
	ctx   context.Context
	store store.AlertStore
	cfg   AlertListenerConfig
}

type AlertListenerConfig struct {
	Port int
}

func (al *AlertListener) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	log.Infof("alert: %s", req.Alert.Id)
	if err := al.store.AddAlert(req.Alert); err != nil {
		return nil, err
	}
	return &protocol.NotifyResponse{}, nil
}

func (al *AlertListener) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:8770")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterQueryNodeServer(grpcServer, al)
	return grpcServer.Serve(lis)
}

func (al *AlertListener) Stop() error {
	log.Infof("Stopping %s", al.Name())
	return nil
}

func (al *AlertListener) Name() string {
	return "AlertListener"
}

func NewAlertListener(ctx context.Context, store store.AlertStore, cfg AlertListenerConfig) (*AlertListener, error) {
	return &AlertListener{
		ctx:   ctx,
		store: store,
		cfg:   cfg,
	}, nil
}
