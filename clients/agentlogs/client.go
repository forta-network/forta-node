package agentlogs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Agent contains agent data.
type Agent struct {
	ID   string `json:"id"`
	Logs string `json:"logs"`
}

// Agents is a type alias of agent slice.
type Agents []*Agent

// Has tells if the list has the same logs for the same agent.
func (agents Agents) Has(agentID, logs string) bool {
	for _, agent := range agents {
		if agent.ID == agentID && agent.Logs == logs {
			return true
		}
	}
	return false
}

// Encode encodes the agent data.
func Encode(agents Agents) (io.Reader, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	defer w.Close()
	return &buf, json.NewEncoder(w).Encode(agents)
}

// Decode decodes the agent data.
func Decode(r io.Reader) (agents Agents, err error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()
	err = json.NewDecoder(gzipReader).Decode(&agents)
	return
}

// Client interacts with the agent logs API.
type Client interface {
	SendLogs(agents Agents, authToken string) error
}

type client struct {
	endpoint string
}

// NewClient creates a new client.
func NewClient(endpoint string) *client {
	return &client{endpoint: endpoint}
}

func (client *client) SendLogs(agents Agents, authToken string) error {
	body, err := Encode(agents)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", client.endpoint, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	if len(authToken) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with code '%d'", resp.StatusCode)
	}
	return nil
}
