package store

import (
	"testing"
)

func TestWebUpdaterStore_GetLatestVersions(t *testing.T) {
	wu := &WebUpdaterStore{
		url: "https://api.defender.openzeppelin.com/autotasks/746f38e9-1c51-4ff1-8753-f03fe99931fc/runs/webhook/62ea5767-415e-412d-aa34-ff31ed60b640/3CpfVP7ndomeenF8L6h4oU",
		ipfs: &ipfsClient{
			gatewayURL: "https://ipfs.forta.network",
		},
	}

	v, err := wu.GetLatestReference()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(v)

}
