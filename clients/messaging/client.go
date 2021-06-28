package messaging

import (
	log "github.com/sirupsen/logrus"

	"github.com/nats-io/nats.go"
)

// Notification and client globals
var (
	BufferSize = 1000

	defaultClient *Client
)

// Start starts the default client.
func Start(natsURL string) {
	defaultClient = NewClient(natsURL)
}

// DefaultClient returns the default client.
func DefaultClient() *Client {
	return defaultClient
}

// Client wraps the NATS client to publish and receive our messages.
type Client struct {
	ec *nats.EncodedConn
}

// NewClient creates and starts a new client.
func NewClient(natsURL string) *Client {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Panicf("failed to connect to nats server: %v", err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Panicf("failed to create the encoded connection: %v", err)
	}
	client := &Client{
		ec: ec,
	}
	return client
}

// Subscribe subscribes the consumer to this client.
func (client *Client) Subscribe(subject string, handler func(interface{}) error) {
	// TODO: Configure redelivery options somehow.
	_, err := client.ec.Subscribe(subject, func(m *nats.Msg) {
		payload, ok := schemaReg(subject)
		if !ok {
			return
		}
		if err := handler(payload); err != nil {
			m.Nak()
			log.Errorf("failed to handle msg: %v", err)
		}
		if err := m.Ack(); err != nil {
			log.Errorf("failed to ack msg: %v", err)
		}
	})
	if err != nil {
		log.Panicf("failed to subscribe to '%s': %v", subject, err)
	}
}

// Publish publishes new messages.
func (client *Client) Publish(subject string, payload interface{}) {
	// TODO: Validate payload type?
	if err := client.ec.Publish(subject, payload); err != nil {
		log.Errorf("failed to publish msg: %v", err)
	}
}

// Subscribe uses the default client to call Subscribe.
func Subscribe(subject string, handler func(interface{}) error) {
	defaultClient.Subscribe(subject, handler)
}

// Publish uses the default client to call Publish.
func Publish(subject string, payload interface{}) {
	defaultClient.Publish(subject, payload)
}
