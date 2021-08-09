package query

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"
)

const (
	DefaultAlertBroadcastInterval = time.Second * 30
	DefaultAlertBroadcastSize     = 50
)

// AlertListener allows retrieval of alerts from the database
type AlertListener struct {
	protocol.UnimplementedQueryNodeServer
	ctx       context.Context
	store     store.AlertStore
	cfg       AlertListenerConfig
	msgClient clients.MessageClient
}

type AlertListenerConfig struct {
	Port int
}

func (al *AlertListener) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	log.Infof("alert: %s", req.SignedAlert.Alert.Id)
	req.SignedAlert.Alert.Published = false // enforce
	if err := al.store.AddAlert(req.SignedAlert); err != nil {
		return nil, err
	}
	return &protocol.NotifyResponse{}, nil
}

func (al *AlertListener) broadcastAlerts(ctx context.Context) {
	ticker := time.NewTicker(DefaultAlertBroadcastInterval)
	for {
		results, err := al.store.QueryAlerts(&store.AlertQueryRequest{
			Limit: DefaultAlertBroadcastSize,
			Criteria: []*store.FilterCriterion{
				{
					Operator: store.Equals,
					Field:    "published",
					Values:   []string{"false"},
				},
			},
		})
		if err != nil {
			log.Errorf("failed to query unpublished alerts: %v", err)
		} else {
			al.msgClient.Publish(messaging.SubjectAlertsStatusPending, results.Alerts)
		}
		<-ticker.C
	}
}

func (al *AlertListener) handleAlertsPublished(payload messaging.AlertsPayload) error {
	for _, alert := range payload {
		if err := al.store.AddAlert(alert); err != nil {
			return fmt.Errorf("failed to update published alert: %v", err)
		}
	}
	return nil
}

func (al *AlertListener) registerMessageHandlers() {
	al.msgClient.Subscribe(messaging.SubjectAlertsStatusPublished, messaging.AlertsHandler(al.handleAlertsPublished))
}

func (al *AlertListener) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:8770")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterQueryNodeServer(grpcServer, al)

	al.registerMessageHandlers()
	go al.broadcastAlerts(al.ctx)

	return grpcServer.Serve(lis)
}

func (al *AlertListener) Stop() error {
	log.Infof("Stopping %s", al.Name())
	return nil
}

func (al *AlertListener) Name() string {
	return "AlertListener"
}

func NewAlertListener(ctx context.Context, store store.AlertStore, cfg AlertListenerConfig, msgClient clients.MessageClient) (*AlertListener, error) {
	return &AlertListener{
		ctx:       ctx,
		store:     store,
		cfg:       cfg,
		msgClient: msgClient,
	}, nil
}
