package scanner

import (
	"sort"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
)

func truncateFinding(finding *protocol.Finding) (truncated bool) {
	sort.Strings(finding.Addresses)

	// truncate finding addresses
	lenFindingAddrs := len(finding.Addresses)

	if lenFindingAddrs > utils.NumMaxAddressesPerAlert {
		finding.Addresses = finding.Addresses[:utils.NumMaxAddressesPerAlert]
		truncated = true
	}

	return truncated
}

func reduceMapToArr(m map[string]bool) (result []string) {
	for s := range m {
		result = append(result, s)
	}

	return
}
