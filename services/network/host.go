package network

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// Host contains host networking configuration.
type Host struct {
	DefaultInterface *net.Interface
	DefaultGateway   *net.IP
	DefaultSubnet    *net.IPNet
}

// DetectHostNetworking detects host networking configuration.
func DetectHostNetworking() (*Host, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get host interfaces: %v", err)
	}

	routes, err := netlink.RouteList(nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get routes list: %v", err)
	}
	defaultRoute := routes[0]
	defaultGatewayAddr := defaultRoute.Gw

	defaultInterface, defaultSubnet, err := getDefaultInterface(ifaces, defaultGatewayAddr)
	if err != nil {
		return nil, err
	}

	return &Host{
		DefaultInterface: defaultInterface,
		DefaultGateway:   &defaultGatewayAddr,
		DefaultSubnet:    defaultSubnet,
	}, nil
}

func getDefaultInterface(ifaces []net.Interface, defaultGatewayAddr net.IP) (*net.Interface, *net.IPNet, error) {
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get interface addresses: %v", err)
		}
		for _, ifaceAddr := range addrs {
			_, ifaceSubnet, err := net.ParseCIDR(ifaceAddr.String())
			if err != nil {
				panic(err)
			}
			if ifaceSubnet.Contains(defaultGatewayAddr) {
				return &iface, ifaceSubnet, nil
			}
		}

	}
	return nil, nil, errors.New("could not find the default interface")
}
