package agentgrpc

import (
	"fmt"
	"time"

	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-node/config"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Client allows us to communicate with an agent.
type Client struct {
	conn *grpc.ClientConn
	protocol.AgentClient
}

// NewClient creates a new client.
func NewClient() *Client {
	return &Client{}
}

// Dial dials an agent using the config.
func (client *Client) Dial(cfg config.AgentConfig) error {
	var (
		conn *grpc.ClientConn
		err  error
	)
	for i := 0; i < 10; i++ {
		conn, err = grpc.Dial(fmt.Sprintf("%s:%s", cfg.ContainerName(), cfg.GrpcPort()), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
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
	client.conn = conn
	client.AgentClient = protocol.NewAgentClient(conn)
	log.Debugf("connected to agent: %s", cfg.ContainerName())
	return nil
}

// Close implements io.Closer.
func (client *Client) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}
