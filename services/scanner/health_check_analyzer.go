package scanner

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/services/components"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	"github.com/forta-network/forta-node/services/components/metrics"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"

	log "github.com/sirupsen/logrus"
)

// HealthCheckAnalyzerService reads TX info, calls agents, and emits results
type HealthCheckAnalyzerService struct {
	ctx context.Context
	cfg HealthCheckAnalyzerServiceConfig

	lastInputActivity  health.TimeTracker
	lastOutputActivity health.TimeTracker
}

type HealthCheckAnalyzerServiceConfig struct {
	AlertSender clients.AlertSender
	MsgClient   clients.MessageClient
	components.BotProcessing
}

func (t *HealthCheckAnalyzerService) publishMetrics(result *botreq.HealthCheckResult) {
	m := metrics.GetHealthCheckMetrics(result.AgentConfig, result.Response, result.Timestamps, result.InvokeError)
	t.cfg.MsgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: m})
}

func (t *HealthCheckAnalyzerService) Start() error {
	go func() {
		for result := range t.cfg.BotProcessing.Results.HealthCheck {
			rt := &clients.AgentRoundTrip{
				AgentConfig:             result.AgentConfig,
				EvalHealthCheckRequest:  result.Request,
				EvalHealthCheckResponse: result.Response,
			}

			if err := t.cfg.AlertSender.NotifyWithoutAlert(
				rt, result.Timestamps,
			); err != nil {
				log.WithError(err).Panic("failed to notify health check response")
			}

			t.publishMetrics(result)

			t.lastOutputActivity.Set()
		}
	}()

	return nil
}

func (t *HealthCheckAnalyzerService) Stop() error {
	return nil
}

func (t *HealthCheckAnalyzerService) Name() string {
	return "health-check-analyzer"
}

// Health implements the health.Reporter interface.
func (t *HealthCheckAnalyzerService) Health() health.Reports {
	return health.Reports{
		t.lastInputActivity.GetReport("event.input.time"),
		t.lastOutputActivity.GetReport("event.output.time"),
	}
}

func NewHealthCheckAnalyzerService(ctx context.Context, cfg HealthCheckAnalyzerServiceConfig) (*HealthCheckAnalyzerService, error) {
	return &HealthCheckAnalyzerService{
		cfg: cfg,
		ctx: ctx,
	}, nil
}
