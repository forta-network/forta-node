package cmd

import (
	"fmt"
	"github.com/forta-network/forta-node/services/registry"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func handleFortaAgentAdd(cmd *cobra.Command, args []string) error {
	reg := registry.New(cfg, common.HexToAddress(cfg.Registry.ContractAddress), nil)
	if err := reg.Init(); err != nil {
		return fmt.Errorf("failed to initialize")
	}
	agentCfg, err := reg.FindAgentGlobally(args[0], parsedArgs.Version)
	if err != nil {
		return fmt.Errorf("failed to load the agent: %v", err)
	}

	agents, err := reg.ReadLocalAgents()
	if err != nil {
		return fmt.Errorf("failed to read the local agents: %v", err)
	}

	for _, localAgent := range agents {
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
	agents = append(agents, &agentCfg)

	if err := reg.WriteLocalAgents(agents); err != nil {
		return fmt.Errorf("failed to write the local agents file: %v", err)
	}

	greenBold("Successfully added agent %s locally! You can see the full list in %s.\n", agentCfg.ID, cfg.LocalAgentsPath)
	fmt.Printf("Image: %s\n", color.New(color.FgYellow).Sprintf(agentCfg.Image))

	return nil
}
