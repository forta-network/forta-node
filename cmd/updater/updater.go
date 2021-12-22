package updater

import (
	"context"
	"strings"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/updater"
	"github.com/forta-protocol/forta-node/store"
)

//TODO: this is throwaway once we use a contract
const devUrl = "https://api.defender.openzeppelin.com/autotasks/746f38e9-1c51-4ff1-8753-f03fe99931fc/runs/webhook/62ea5767-415e-412d-aa34-ff31ed60b640/3CpfVP7ndomeenF8L6h4oU"
const prodUrl = "https://api.defender.openzeppelin.com/autotasks/aaea6377-66be-4836-854f-db47820f746b/runs/webhook/62ea5767-415e-412d-aa34-ff31ed60b640/G62StUaANRRLpTXToEiR1n"

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	ipfs := store.NewIPFSClient(cfg.Registry.IPFS.GatewayURL)

	//TODO: this is throwaway once we use a contract
	url := prodUrl
	if strings.Contains(cfg.Registry.IPFS.GatewayURL, "-dev") {
		url = devUrl
	}
	ws := store.NewWebUpdaterStore(url)
	return []services.Service{
		updater.NewUpdaterService(ctx, ws, ipfs, config.DefaultUpdaterPort),
	}, nil
}

func Run() {
	services.ContainerMain("updater", initServices)
}
