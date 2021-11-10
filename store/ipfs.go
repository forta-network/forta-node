package store

import (
	"fmt"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	"github.com/forta-protocol/forta-node/services/registry/regtypes"
)

type IPFSClient interface {
	GetAgentFile(cid string) (*regtypes.AgentFile, error)
}

type ipfsClient struct {
	gatewayURL string
}

func (client *ipfsClient) GetAgentFile(cid string) (*regtypes.AgentFile, error) {
	resp, err := http.Get(fmt.Sprintf("%s/ipfs/%s", client.gatewayURL, cid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var agentData regtypes.AgentFile
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of metadata: %v", err)
	}
	if err := json.Unmarshal(b, &agentData); err != nil {
		log.WithField("metadata", string(b)).Error("could not decode metadata")
		return nil, fmt.Errorf("failed to decode the agent file: %v", err)
	}
	return &agentData, nil
}
