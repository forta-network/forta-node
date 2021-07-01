package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

// Notification and client globals
var (
	BufferSize = 1000
)

// Client wraps the NATS client to publish and receive our messages.
type Client struct {
	logger *log.Entry
	nc     *nats.Conn
}

// NewClient creates and starts a new client.
func NewClient(name, natsURL string) *Client {
	logger := log.WithField("name", fmt.Sprintf("%s/messaging", name)).WithField("nats", natsURL)
	logger.Infof("connecting to: %s", natsURL)
	var (
		nc  *nats.Conn
		err error
	)
	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(natsURL)
		if err == nil {
			break
		}
		err = fmt.Errorf("failed to connect to nats server: %v", err)
		logger.Error(err)
	}
	if err != nil {
		logger.Panic(err)
	}
	client := &Client{
		logger: logger,
		nc:     nc,
	}
	return client
}

// AgentsHandler handles agents.* subjects.
type AgentsHandler func(AgentPayload) error

// Subscribe subscribes the consumer to this client.
func (client *Client) Subscribe(subject string, handler interface{}) {
	// TODO: Configure redelivery options somehow.
	logger := client.logger.WithField("subject", subject)
	_, err := client.nc.Subscribe(subject, func(m *nats.Msg) {
		logger.Debugf("received: %s", string(m.Data))

		var err error
		switch h := handler.(type) {
		case AgentsHandler:
			var payload AgentPayload
			err = json.Unmarshal(m.Data, &payload)
			if err != nil {
				break
			}
			err = h(payload)

		default:
			logger.Panicf("no handler found")
		}

		if err != nil {
			m.Nak()
			logger.Errorf("failed to handle msg: %v", err)
		}
	})
	if err != nil {
		logger.Panicf("failed to subscribe: %v", err)
	}
	logger.Info("subscribed")
}

// Publish publishes new messages.
func (client *Client) Publish(subject string, payload interface{}) {
	logger := client.logger.WithField("subject", subject)
	data, _ := json.Marshal(payload)
	if err := client.nc.Publish(subject, data); err != nil {
		logger.Errorf("failed to publish msg: %v", err)
	}
	logger.Debugf("published: %s", string(data))
}
