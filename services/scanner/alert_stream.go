package scanner

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/feeds"
	log "github.com/sirupsen/logrus"
)

// AlertStreamService pulls TX info from providers and emits to channel
type AlertStreamService struct {
	cfg         AlertStreamServiceConfig
	ctx         context.Context
	alertOutput chan *domain.AlertEvent
	alertFeed   feeds.AlertFeed

	lastAlertActivity health.TimeTracker
}

type AlertStreamServiceConfig struct {
}

func (t *AlertStreamService) ReadOnlyAlertStream() <-chan *domain.AlertEvent {
	return t.alertOutput
}

func (t *AlertStreamService) handleAlert(evt *domain.AlertEvent) error {
	select {
	case <-t.ctx.Done():
		return nil
	default:
	}
	t.alertOutput <- evt
	t.lastAlertActivity.Set()
	return nil
}

func (t *AlertStreamService) Start() error {
	go func() {
		if err := t.alertFeed.ForEachAlert(t.handleAlert); err != nil {
			logger := log.WithError(err)
			if err != context.Canceled {
				logger.Panic("tx feed error")
			}
			logger.Info("tx feed stopped")
		}
	}()
	return nil
}

func (t *AlertStreamService) Stop() error {
	if t.alertOutput != nil {
		// drain and close block channel
		func(c chan *domain.AlertEvent) {
			for {
				select {
				case a := <-c:
					log.WithFields(log.Fields{"tx": a.Alert.Alert.Id}).Info("gracefully draining block")
				default:
					close(c)
					return
				}
			}
		}(t.alertOutput)
	}
	return nil
}

func (t *AlertStreamService) Name() string {
	return "alert-stream"
}

// Health implements health.Reporter interface.
func (t *AlertStreamService) Health() health.Reports {
	return health.Reports{
		t.lastAlertActivity.GetReport("event.alert.time"),
	}
}

func NewAlertStreamService(ctx context.Context, alertFeed feeds.AlertFeed, cfg AlertStreamServiceConfig) (*AlertStreamService, error) {
	alertOutput := make(chan *domain.AlertEvent)

	return &AlertStreamService{
		cfg:         cfg,
		ctx:         ctx,
		alertOutput: alertOutput,
		alertFeed:   alertFeed,
	}, nil
}
