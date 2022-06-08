package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/forta-network/forta-node/config"
	"github.com/peterbourgon/unixtransport"
)

var unixTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func init() {
	unixtransport.Register(unixTransport)
}

type botAdminUnixSockClient struct {
	client  *http.Client
	sockURL string
}

// NewUnixSockClient returns a bot admin implementation as a client.
func NewUnixSockClient(containerName string) BotAdmin {
	return &botAdminUnixSockClient{
		client:  &http.Client{Transport: unixTransport},
		sockURL: sockURL(containerName),
	}
}

func (ba *botAdminUnixSockClient) IPTables(ruleCmds [][]string) error {
	b, _ := json.Marshal(ruleCmds)
	resp, err := ba.client.Post(ba.sockURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("request to bot admin failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bot admin failed to execute the rules")
	}
	return nil
}

func sockURL(containerName string) string {
	return fmt.Sprintf("http+unix://%s:/", sockPath(containerName))
}

func sockPath(containerName string) string {
	return path.Join(BotAdminSockDir(), containerName)
}

// BotAdminSockDir returns the bot admin dir.
func BotAdminSockDir() string {
	return path.Join(config.DefaultContainerFortaDirPath, ".sock", "botadmin")
}
