package main

import (
	json_rpc "github.com/forta-protocol/forta-node/cmd/json-rpc"
	"github.com/forta-protocol/forta-node/cmd/publisher"
	"github.com/forta-protocol/forta-node/cmd/scanner"
	"github.com/spf13/cobra"
)

var (
	cmdFortaNode = &cobra.Command{
		Use: "forta-node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}

	cmdScanner = &cobra.Command{
		Use: "scanner",
		RunE: func(cmd *cobra.Command, args []string) error {
			scanner.Run()
			return nil
		},
	}

	cmdPublisher = &cobra.Command{
		Use: "publisher",
		RunE: func(cmd *cobra.Command, args []string) error {
			publisher.Run()
			return nil
		},
	}

	cmdJsonRpc = &cobra.Command{
		Use: "json-rpc",
		RunE: func(cmd *cobra.Command, args []string) error {
			json_rpc.Run()
			return nil
		},
	}
)

func init() {
	cmdFortaNode.AddCommand(cmdScanner)
	cmdFortaNode.AddCommand(cmdPublisher)
	cmdFortaNode.AddCommand(cmdJsonRpc)
}

func main() {
	cmdFortaNode.Execute()
}
