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

// MetaAlertStreamService pulls alert info from providers and emits to channel
type MetaAlertStreamService struct {
	cfg         MetaAlertStreamServiceConfig
	ctx         context.Context
	alertOutput chan *domain.AlertEvent
	alertFeed   feeds.AlertFeed
	msgClient   clients.MessageClient

	subscribeChan     chan string
	unSubscribeChan   chan string
	lastAlertActivity health.TimeTracker
}

type MetaAlertStreamServiceConfig struct {
}

func (t *MetaAlertStreamService) registerMessageHandlers() {
	t.msgClient.Subscribe(messaging.SubjectAgentsAlertSubscribe, messaging.SubscriptionHandler(t.handleMessageSubscribe))
	t.msgClient.Subscribe(messaging.SubjectAgentsAlertUnsubscribe, messaging.SubscriptionHandler(t.handleMessageUnsubscribe))
}

func (t *MetaAlertStreamService) ReadOnlyAlertStream() <-chan *domain.AlertEvent {
	return t.alertOutput
}

func (t *MetaAlertStreamService) handleAlert(evt *domain.AlertEvent) error {
	select {
	case <-t.ctx.Done():
		return nil
	default:
	}
	log.Infof("new alert incoming %s", evt.Event.Alert.Source.Bot.Id)
	t.alertOutput <- evt
	t.lastAlertActivity.Set()
	return nil
}

func (t *MetaAlertStreamService) Start() error {
	t.registerMessageHandlers()
	go func() {
		if err := t.alertFeed.ForEachAlert(t.handleAlert); err != nil {
			logger := log.WithError(err)
			if err != context.Canceled {
				logger.Panic("alert feed error")
			}
			logger.Info("alert feed stopped")
		}
	}()
	return nil
}

func (t *MetaAlertStreamService) Stop() error {
	if t.alertOutput != nil {
		// drain and close block channel
		func(c chan *domain.AlertEvent) {
			for {
				select {
				case a := <-c:
					log.WithFields(log.Fields{"alert": a.Event.Alert.Source.Bot.Id}).Info("gracefully draining alert")
				default:
					close(c)
					return
				}
			}
		}(t.alertOutput)
	}
	return nil
}

func (t *MetaAlertStreamService) Name() string {
	return "alert-stream"
}

// Health implements health.Reporter interface.
func (t *MetaAlertStreamService) Health() health.Reports {
	return health.Reports{
		t.lastAlertActivity.GetReport("event.alert.time"),
	}
}

func (t *MetaAlertStreamService) handleMessageSubscribe(payload messaging.SubscriptionPayload) error {
	for _, cfg := range payload {
		t.alertFeed.AddSubscription(cfg.Dst, cfg.Src)
	}

	return nil
}

func (t *MetaAlertStreamService) handleMessageUnsubscribe(payload messaging.SubscriptionPayload) error {
	for _, cfg := range payload {
		t.alertFeed.RemoveSubscription(cfg.Dst, cfg.Src)
	}

	return nil
}

func NewMetaAlertStreamService(ctx context.Context, alertFeed feeds.AlertFeed, msgClient clients.MessageClient, cfg MetaAlertStreamServiceConfig) (*MetaAlertStreamService, error) {
	alertOutput := make(chan *domain.AlertEvent)

	return &MetaAlertStreamService{
		cfg:         cfg,
		ctx:         ctx,
		msgClient:   msgClient,
		alertOutput: alertOutput,
		alertFeed:   alertFeed,
	}, nil
}
