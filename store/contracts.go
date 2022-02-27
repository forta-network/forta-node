package store

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type scannerNodeVersion interface {
	ScannerNodeVersion(opts *bind.CallOpts) (string, error)
}
