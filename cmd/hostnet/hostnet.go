package hostnet

import (
	"github.com/forta-network/forta-node/services/network"
)

func Run() {
	host, err := network.DetectHostNetworking()
	if err != nil {
		panic(err)
	}
	network.OutputHostNetworking(host)
}
