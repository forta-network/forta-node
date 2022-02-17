package store

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services/registry/regtypes"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
)

type IPFSClient interface {
	GetAgentFile(cid string) (*regtypes.AgentFile, error)
	GetReleaseManifest(cid string) (*config.ReleaseManifest, error)
}

type ipfsClient struct {
	gatewayURL string
}

func NewIPFSClient(gatewayUrl string) *ipfsClient {
	return &ipfsClient{
		gatewayURL: gatewayUrl,
	}
}

func (client *ipfsClient) getJson(cid string, target interface{}) error {
	resp, err := http.Get(fmt.Sprintf("%s/ipfs/%s", client.gatewayURL, cid))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body of metadata: %v", err)
	}
	if err := json.Unmarshal(b, &target); err != nil {
		log.WithError(err).WithField("metadata", string(b)).Error("could not decode ipfs data")
		return err
	}
	return nil
}

func (client *ipfsClient) GetReleaseManifest(cid string) (*config.ReleaseManifest, error) {
	var rm config.ReleaseManifest
	if err := client.getJson(cid, &rm); err != nil {
		log.WithError(err).Error("could not decode release metadata")
		return nil, fmt.Errorf("failed to decode the release file: %v", err)
	}
	return &rm, nil
}

func (client *ipfsClient) GetAgentFile(cid string) (*regtypes.AgentFile, error) {
	var agentData regtypes.AgentFile
	if err := client.getJson(cid, &agentData); err != nil {
		log.WithError(err).Error("could not decode agent metadata")
		return nil, fmt.Errorf("failed to decode the agent file: %v", err)
	}
	return &agentData, nil
}
