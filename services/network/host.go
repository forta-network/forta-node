package network

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

// Host contains host networking configuration.
type Host struct {
	DefaultInterfaceName string
	DefaultSubnet        string
	DefaultGateway       string

	Docker0Subnet string
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

	docker0Subnet, err := getDocker0Subnet(ifaces)
	if err != nil {
		return nil, err
	}

	return &Host{
		DefaultInterfaceName: defaultInterface.Name,
		DefaultGateway:       defaultGatewayAddr.String(),
		DefaultSubnet:        defaultSubnet.String(),

		Docker0Subnet: docker0Subnet.String(),
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
				return nil, nil, fmt.Errorf("failed to parse interface address: %v", err)
			}
			if ifaceSubnet.Contains(defaultGatewayAddr) {
				return &iface, ifaceSubnet, nil
			}
		}

	}
	return nil, nil, errors.New("could not find the default interface")
}

func getDocker0Subnet(ifaces []net.Interface) (*net.IPNet, error) {
	for _, iface := range ifaces {
		if iface.Name != "docker0" {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to get interface addresses: %v", err)
		}
		for _, ifaceAddr := range addrs {
			ifaceIpAddr, ifaceSubnet, err := net.ParseCIDR(ifaceAddr.String())
			if err != nil {
				return nil, fmt.Errorf("failed to parse interface address: %v", err)
			}
			if ifaceIpAddr.To4() == nil {
				continue
			}
			return ifaceSubnet, nil
		}
	}
	return nil, errors.New("could not find docker0 interface")
}

// OutputHostNetworking outputs host networking info.
func OutputHostNetworking(host *Host) {
	fmt.Printf(
		"%s %s %s %s",
		host.DefaultInterfaceName, host.DefaultSubnet, host.DefaultGateway, host.Docker0Subnet,
	)
}

// ReadHostNetworking outputs host networking info.
func ReadHostNetworking(output string) *Host {
	parts := strings.Split(output, " ")
	return &Host{
		DefaultInterfaceName: parts[0],
		DefaultSubnet:        parts[1],
		DefaultGateway:       parts[2],

		Docker0Subnet: parts[3],
	}
}
