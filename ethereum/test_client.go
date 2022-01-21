package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// TestAPI tests given API URL.
func TestAPI(ctx context.Context, rawurl string) error {
	client, err := ethclient.Dial(rawurl)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	_, err = client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %v", err)
	}
	return nil
}
