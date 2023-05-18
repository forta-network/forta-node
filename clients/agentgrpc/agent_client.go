package agentgrpc

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const defaultAgentResponseMaxByteCount = 250000 // 250K

// Method is gRPC method type.
type Method string

// Agent gRPC methods
const (
	MethodInitialize    Method = "/network.forta.Agent/Initialize"
	MethodEvaluateTx    Method = "/network.forta.Agent/EvaluateTx"
	MethodEvaluateBlock Method = "/network.forta.Agent/EvaluateBlock"
	MethodEvaluateAlert Method = "/network.forta.Agent/EvaluateAlert"
)

// Client makes the gRPC requests to evaluate block and txs and receive results.
type Client interface {
	Dial(config.AgentConfig) error
	Invoke(ctx context.Context, method Method, in, out interface{}, opts ...grpc.CallOption) error
	protocol.AgentClient
	io.Closer
}

// client allows us to communicate with an agent.
type client struct {
	conn *grpc.ClientConn
	protocol.AgentClient
}

// NewClient creates a new client.
func NewClient() *client {
	return &client{}
}

// Dial dials an agent using the config.
func (client *client) Dial(cfg config.AgentConfig) error {
	var (
		conn *grpc.ClientConn
		err  error
	)
	for i := 0; i < 10; i++ {
		conn, err = grpc.Dial(
			fmt.Sprintf("%s:%s", cfg.ContainerName(), cfg.GrpcPort()),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(10*time.Second),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultAgentResponseMaxByteCount)),
		)
		if err == nil {
			break
		}
		err = fmt.Errorf("failed to connect to agent '%s': %v", cfg.ContainerName(), err)
		log.Debug(err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	client.WithConn(conn)
	log.Debugf("connected to agent: %s", cfg.ContainerName())
	return nil
}

// WithConn sets the client conn.
func (client *client) WithConn(conn *grpc.ClientConn) {
	client.conn = conn
	client.AgentClient = protocol.NewAgentClient(conn)
}

// Invoke is a generalization of client methods.
func (client *client) Invoke(ctx context.Context, method Method, in, out interface{}, opts ...grpc.CallOption) error {
	return client.conn.Invoke(ctx, string(method), in, out, opts...)
}

// Close implements io.Closer.
func (client *client) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}
