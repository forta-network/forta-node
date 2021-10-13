package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/forta-protocol/forta-node/store"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/forta-protocol/forta-node/config"
	"github.com/spf13/cobra"
)

func handleFortaAgentAdd(cmd *cobra.Command, args []string) error {
	reg, err := store.NewRegistryStore(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize registry")
	}

	agentCfg, err := reg.FindAgentGlobally(args[0])
	if err != nil {
		return fmt.Errorf("failed to load the agent: %v", err)
	}

	cfg.LocalAgents, err = readLocalAgents()
	if err != nil {
		return fmt.Errorf("failed to read the local agents: %v", err)
	}

	for _, localAgent := range cfg.LocalAgents {
		if localAgent.ID != agentCfg.ID {
			continue
		}
		if localAgent.Image == agentCfg.Image {
			cmd.Println("Already added to list - ignored")
			return nil
		}
	}
	// Two cases to add append an agent:
	//  1. Different agent
	//  2. Same agent, different image (i.e. different version)
	cfg.LocalAgents = append(cfg.LocalAgents, agentCfg)

	if err := writeLocalAgents(cfg.LocalAgents); err != nil {
		return fmt.Errorf("failed to write the local agents file: %v", err)
	}

	greenBold("Successfully added agent %s locally! You can see the full list in %s.\n", agentCfg.ID, cfg.LocalAgentsPath)
	fmt.Printf("Image: %s\n", color.New(color.FgYellow).Sprintf(agentCfg.Image))

	return nil
}

// readLocalAgents tries to read the local agents and silently returns an
// empty array if the file is not readable or not found.
func readLocalAgents() ([]*config.AgentConfig, error) {
	var agents []*config.AgentConfig
	b, err := ioutil.ReadFile(cfg.LocalAgentsPath)
	if err == nil {
		if err := json.Unmarshal(b, &agents); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the local agents file: %v", err)
		}
	}
	for _, agent := range agents {
		agent.IsLocal = true
	}
	return agents, nil
}

// writeLocalAgents writes the agents to the local list.
func writeLocalAgents(agents []*config.AgentConfig) error {
	if len(agents) == 0 {
		return nil
	}
	b, _ := json.MarshalIndent(agents, "", "  ")
	return ioutil.WriteFile(cfg.LocalAgentsPath, b, 0644)
}
