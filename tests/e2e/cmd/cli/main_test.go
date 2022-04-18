package cli_test

import (
	"log"
	"testing"

	"github.com/forta-network/forta-node/cmd"
)

func TestCLI(t *testing.T) {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
