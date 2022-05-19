package network

import (
	"context"
	"fmt"
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
	proxyIpAddress, err := bm.getContainerIpAddress(config.DockerJSONRPCProxyContainerName, config.DockerServiceNetworkName)
	if err != nil {
		return err
	}
	scannerIpAddress, err := bm.getContainerIpAddress(config.DockerScannerContainerName, config.DockerServiceNetworkName)
	if err != nil {
		return err
	}

	ruleCmds := [][]string{
		// clear all first
		{"-F"},

		// allow making JSON-RPC requests to the proxy container
		{"-A", "OUTPUT", "-d", proxyIpAddress, "-j", "ACCEPT"},

		// allow responding to gRPC requests from the scanner container
		{"-A", "OUTPUT", "-d", scannerIpAddress, "-j", "ACCEPT"},

		// allow internet connectivity
		{"-A", "OUTPUT", "-d", bm.defaultGateway.String(), "-j", "ACCEPT"},
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

func (bm *botManager) getContainerIpAddress(containerName, networkName string) (string, error) {
	container, err := bm.dockerClient.GetContainerByName(bm.ctx, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to get '%s' container: %v", containerName, err)
	}
	network, ok := container.NetworkSettings.Networks[networkName]
	if !ok {
		return "", fmt.Errorf("container '%s' is not on the '%s' network: %v", containerName, networkName, err)
	}
	return network.IPAddress, nil
}
