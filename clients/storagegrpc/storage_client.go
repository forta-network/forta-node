package storagegrpc

import (
	"fmt"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Client allows us to communicate with an agent.
type Client struct{}

// NewClient creates a new client.
func NewClient() *Client {
	return &Client{}
}

// Dial dials an agent using the config.
func Dial(serverURL string) (protocol.StorageClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	for i := 0; i < 10; i++ {
		conn, err = grpc.Dial(
			serverURL,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(10*time.Second),
		)
		if err == nil {
			break
		}
		err = fmt.Errorf("failed to connect to storage '%s': %v", serverURL, err)
		log.Debug(err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("connected to storage: %s", serverURL)
	return protocol.NewStorageClient(conn), nil
}
