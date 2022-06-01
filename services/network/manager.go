package network

import (
	"context"
	"net"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
)

// BotManager manages bot networking.
type BotManager interface {
	SetBotAdminRules(containerName string) error
}

// botManager keeps track of networking rules and forwards to administrators.
type botManager struct {
	ctx          context.Context
	dockerClient clients.DockerClient

	defaultGateway *net.IP
	allSubnets     []*net.IPNet
}

// NewBotManager creates a new bot manager.
func NewBotManager(
	ctx context.Context, dockerClient clients.DockerClient, defaultGateway *net.IP, allSubnets []*net.IPNet,
) *botManager {
	return &botManager{
		ctx:            ctx,
		dockerClient:   dockerClient,
		defaultGateway: defaultGateway,
		allSubnets:     allSubnets,
	}
}

// SetBotAdminRules sets the default rules for any bot.
func (bm *botManager) SetBotAdminRules(containerName string) error {
	proxyIpAddress, err := bm.dockerClient.GetContainerIPAddress(
		bm.ctx,
		config.DockerJSONRPCProxyContainerName,
		config.DockerBotNetworkName,
	)
	if err != nil {
		return err
	}
	scannerIpAddress, err := bm.dockerClient.GetContainerIPAddress(
		bm.ctx,
		config.DockerScannerContainerName,
		config.DockerBotNetworkName,
	)
	if err != nil {
		return err
	}

	ruleCmds := [][]string{
		// clear all first
		{"-F"},

		// allow making JSON-RPC requests to the proxy container
		{"-A", "OUTPUT", "-p", "tcp", "--dport", "8545", "-d", proxyIpAddress, "-j", "ACCEPT"},
		// allow responses to them
		{"-A", "INPUT", "-s", proxyIpAddress, "-j", "ACCEPT"},

		// allow gRPC requests from the scanner container
		{"-A", "INPUT", "-s", scannerIpAddress, "-j", "ACCEPT"},
		// allow only responding to those requests
		{"-A", "OUTPUT", "-d", scannerIpAddress, "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT"},
	}
	// finally, restrict access to all subnets by default
	for _, subnet := range bm.allSubnets {
		ruleCmds = append(
			ruleCmds,
			[]string{"-A", "OUTPUT", "-d", subnet.String(), "-j", "DROP"},
		)
	}

	return NewUnixSockClient(containerName).IPTables(ruleCmds)
}
