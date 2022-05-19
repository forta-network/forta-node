package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/forta-network/forta-node/config"
)

type botAdminUnixSockClient struct {
	sockURL string
}

// NewUnixSockClient returns a bot admin implementation as a client.
func NewUnixSockClient(containerName string) BotAdmin {
	return &botAdminUnixSockClient{
		sockURL: sockURL(containerName),
	}
}

func (ba *botAdminUnixSockClient) IPTables(ruleCmds [][]string) error {
	b, _ := json.Marshal(ruleCmds)
	resp, err := http.Post(ba.sockURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("request to bot admin failed: %v", err)
	}
	resp.Body.Close()
	return nil
}

func sockURL(containerName string) string {
	return fmt.Sprintf("http://%s", sockPath(containerName))
}

func sockPath(containerName string) string {
	return path.Join(config.DefaultContainerFortaDirPath, ".botadmin", fmt.Sprintf("%s.sock", containerName))
}
