package main

import (
	"os"

	"github.com/forta-network/forta-node/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
