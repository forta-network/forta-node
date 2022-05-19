package hostnet

import (
	"fmt"

	"github.com/forta-network/forta-node/services/network"
)

func Run() {
	host, err := network.DetectHostNetworking()
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"%s %s %s",
		host.DefaultInterface.Name, host.DefaultSubnet, host.DefaultGateway,
	)
}
