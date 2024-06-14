package messaging

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/goccy/go-json"
	"github.com/golang/protobuf/proto"
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
	nc, err = nats.Connect(natsURL,
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("Got disconnected! Reason: %s\n", err.Error())
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Printf("Connection closed. Reason: %s", nc.LastError().Error())
		}),
	)

	if err != nil {
		logger.Panic(fmt.Errorf("failed to connect to nats server %w", err))
	}

	logger.Info("successfully connected")
	client := &Client{
		logger: logger,
		nc:     nc,
	}

	return client
}

// AgentsHandler handles agents.* subjects.
type AgentsHandler func(AgentPayload) error
type SubscriptionHandler func(SubscriptionPayload) error
type AgentMetricHandler func(*protocol.AgentMetricList) error
type InspectionResultsHandler func(results *protocol.InspectionResults) error
type ScannerHandler func(ScannerPayload) error

// Subscribe subscribes the consumer to this client.
func (client *Client) Subscribe(subject string, handler interface{}) {
	// TODO: Configure redelivery options somehow.
	logger := client.logger.WithField("subject", subject)
	_, err := client.nc.Subscribe(subject, func(m *nats.Msg) {
		logger.Tracef("received: %s", string(m.Data))

		var err error
		switch h := handler.(type) {
		case AgentsHandler:
			var payload AgentPayload
			err = json.Unmarshal(m.Data, &payload)
			if err != nil {
				break
			}
			err = h(payload)

		case AgentMetricHandler:
			var payload protocol.AgentMetricList
			err = proto.Unmarshal(m.Data, &payload)
			if err != nil {
				break
			}
			err = h(&payload)

		case InspectionResultsHandler:
			var payload protocol.InspectionResults
			err = proto.Unmarshal(m.Data, &payload)
			if err != nil {
				break
			}
			err = h(&payload)

		case ScannerHandler:
			var payload ScannerPayload
			err = json.Unmarshal(m.Data, &payload)
			if err != nil {
				break
			}
			err = h(payload)
		case SubscriptionHandler:
			var payload SubscriptionPayload
			err = json.Unmarshal(m.Data, &payload)
			if err != nil {
				break
			}
			err = h(payload)

		default:
			logger.Panicf("no handler found")
		}

		if err != nil {
			// TODO: Replace nak with whatever is recent.
			// if err := m.Nak() (); err != nil {
			// 	logger.Errorf("failed to send nak: %v", err)
			// }
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
	logger.Tracef("published: %s", string(data))
}

// PublishProto publishes new messages.
func (client *Client) PublishProto(subject string, payload proto.Message) {
	logger := client.logger.WithField("subject", subject)
	data, _ := proto.Marshal(payload)
	if err := client.nc.Publish(subject, data); err != nil {
		logger.Errorf("failed to publish msg: %v", err)
	}
	logger.Tracef("published: %s", string(data))
}

type nopClient struct{}

func NewNopClient() *nopClient {
	return &nopClient{}
}

func (sc *nopClient) Subscribe(subject string, handler interface{}) {

}

func (sc *nopClient) Publish(subject string, payload interface{}) {

}

func (sc *nopClient) PublishProto(subject string, payload proto.Message) {

}
