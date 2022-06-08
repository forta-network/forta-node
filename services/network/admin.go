package network

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Defaults
const (
	DefaultTable = "filter"
	DefaultChain = "BOTADMIN"
	OutputChain  = "OUTPUT"
)

// BotAdmin administrates bot networking.
type BotAdmin interface {
	IPTables(ruleCmds [][]string) error
}

// botAdmin executes networking rules.
type botAdmin struct {
}

// IPTables sets iptables rules.
func (ba *botAdmin) IPTables(ruleCmds [][]string) error {
	for _, ruleCmd := range ruleCmds {
		if err := run(ruleCmd...); err != nil {
			return err
		}
	}
	return nil
}

func run(args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("iptables", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command '%+v' failed: %v: %s", args, err, stderr.String())
	}
	return nil
}
