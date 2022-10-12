package scanner

import (
	"context"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	log "github.com/sirupsen/logrus"
)

// AlertStreamService pulls TX info from providers and emits to channel
type AlertStreamService struct {
	cfg         AlertStreamServiceConfig
	ctx         context.Context
	alertOutput chan *domain.AlertEvent
	alertFeed   feeds.AlertFeed
	msgClient   clients.MessageClient

	subscribeChan     chan string
	unSubscribeChan   chan string
	lastAlertActivity health.TimeTracker
}

type AlertStreamServiceConfig struct {
}

func (t *AlertStreamService) registerMessageHandlers() {
	t.msgClient.Subscribe(messaging.SubjectAgentsAlertSubscribe, messaging.SubscriptionHandler(t.handleMessageSubscribe))
	t.msgClient.Subscribe(messaging.SubjectAgentsAlertUnsubscribe, messaging.SubscriptionHandler(t.handleMessageUnsubscribe))
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
	t.registerMessageHandlers()
	go func() {
		if err := t.alertFeed.ForEachAlert(t.handleAlert); err != nil {
			logger := log.WithError(err)
			if err != context.Canceled {
				logger.Panic("tx feed error")
			}
			logger.Info("alert feed stopped")
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

func (t *AlertStreamService) handleMessageSubscribe(payload messaging.SubscriptionPayload) error {
	for _, cfg := range payload {
		t.alertFeed.AddSubscription(cfg.Dst, cfg.Src)
	}

	return nil
}

func (t *AlertStreamService) handleMessageUnsubscribe(payload messaging.SubscriptionPayload) error {
	for _, cfg := range payload {
		t.alertFeed.RemoveSubscription(cfg.Dst, cfg.Src)
	}

	return nil
}

func NewAlertStreamService(ctx context.Context, alertFeed feeds.AlertFeed, msgClient clients.MessageClient, cfg AlertStreamServiceConfig) (*AlertStreamService, error) {
	alertOutput := make(chan *domain.AlertEvent)

	return &AlertStreamService{
		cfg:         cfg,
		ctx:         ctx,
		msgClient:   msgClient,
		alertOutput: alertOutput,
		alertFeed:   alertFeed,
	}, nil
}
