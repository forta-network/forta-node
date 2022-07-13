package sla_tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-inspector-bot/inspect"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/sirupsen/logrus"
)

// SLATracker runs continuous inspections.
type SLATracker struct {
	ctx    context.Context
	cancel context.CancelFunc

	msgClient clients.MessageClient

	lastErr health.ErrorTracker

	cfg SLATrackerConfig
}

type SLATrackerConfig struct {
	JSONRpcHost string
	JSONRpcPort string
}

func (c SLATrackerConfig) jsonRpcURL() string {
	return fmt.Sprintf("http://%s:%s", c.JSONRpcHost, c.JSONRpcPort)
}

func (p *SLATracker) Start() error {
	for {
		select {
		case <-p.ctx.Done():
			return nil
		default:
			p.inspectionWorker(p.ctx)
		}
	}
}

func (p *SLATracker) inspectionWorker(ctx context.Context) {
	inspection, err := inspect.RunAllInspections(ctx, p.cfg.jsonRpcURL())
	if err != nil {
		logrus.WithError(err).Warnf("errors during inspections")
	}

	slaChecks := &protocol.SLACheckList{}

	for key, value := range inspection {
		slaChecks.Checks = append(slaChecks.Checks, &protocol.SLACheck{Name: key, Value: value, Timestamp: time.Now().Format(time.RFC3339)})
	}

	p.msgClient.PublishProto(messaging.SubjectMetricSLA, slaChecks)

	return
}

func (p *SLATracker) Stop() error {
	p.cancel()

	return nil
}

func (p *SLATracker) Name() string {
	return "sla-tracker"
}

// Health implements health.Reporter interface.
func (p *SLATracker) Health() health.Reports {
	return health.Reports{
		p.lastErr.GetReport("api"),
	}
}

func NewSLATracker(c context.Context, cfg SLATrackerConfig) (*SLATracker, error) {
	msgClient := messaging.NewClient("sla-tracker", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	ctx, cancel := context.WithCancel(c)
	return &SLATracker{
		ctx:       ctx,
		msgClient: msgClient,
		cancel:    cancel,
		cfg:       cfg,
	}, nil
}
