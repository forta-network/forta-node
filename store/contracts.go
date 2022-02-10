package store

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// implemented in contracts/dispatch.go
type dispatch interface {
	ScannerHash(opts *bind.CallOpts, scannerId *big.Int) (struct {
		Length   *big.Int
		Manifest [32]byte
	}, error)

	AgentRefAt(opts *bind.CallOpts, scannerId *big.Int, pos *big.Int) (struct {
		AgentId      *big.Int
		Enabled      bool
		AgentVersion *big.Int
		Metadata     string
		ChainIds     []*big.Int
	}, error)
}

type agentRegistry interface {
	GetAgent(opts *bind.CallOpts, agentId *big.Int) (struct {
		Version  *big.Int
		Metadata string
		ChainIds []*big.Int
	}, error)
}

type scannerNodeVersion interface {
	ScannerNodeVersion(opts *bind.CallOpts) (string, error)
}
