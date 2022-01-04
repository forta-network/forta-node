package ethereum

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mock_ethereum "github.com/forta-protocol/forta-node/ethereum/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var testAddr = common.Address{}

func TestFastBackend(t *testing.T) {
	r := require.New(t)

	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))

	// Given that the API nonce is higher
	backend := &contractBackend{ContractBackend: mockBackend, nonce: 100}
	apiNonce := backend.nonce + 1
	mockBackend.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(apiNonce, nil).Times(2)
	mockBackend.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)

	// When the nonce is requested from the backend
	txNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	// Then it should return the API nonce with no errors
	r.NoError(err)
	r.Equal(apiNonce, txNonce)

	// And new transaction should cause a higher local nonce to be returned
	testTx := types.NewTransaction(txNonce, testAddr, big.NewInt(1), 21000, big.NewInt(30), []byte{})
	r.NoError(backend.SendTransaction(context.Background(), testTx))
	postTxNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	r.NoError(err)
	r.Equal(backend.nonce, postTxNonce)
	r.Greater(postTxNonce, txNonce)
}

func TestLaggingBackend(t *testing.T) {
	r := require.New(t)

	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))

	// Given that the local nonce is higher
	backend := &contractBackend{ContractBackend: mockBackend, nonce: 100}
	apiNonce := backend.nonce - 1
	mockBackend.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(apiNonce, nil).Times(2)
	mockBackend.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)

	// When the nonce is requested from the backend
	txNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	// Then it should return the local nonce with no errors
	r.NoError(err)
	r.Equal(backend.nonce, txNonce)

	// And new transaction should cause a higher local nonce to be returned
	testTx := types.NewTransaction(txNonce, testAddr, big.NewInt(1), 21000, big.NewInt(30), []byte{})
	r.NoError(backend.SendTransaction(context.Background(), testTx))
	postTxNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	r.NoError(err)
	r.Equal(backend.nonce, postTxNonce)
	r.Greater(postTxNonce, txNonce)
}

func TestDriftingLocal(t *testing.T) {
	r := require.New(t)

	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))

	// Given that the local nonce is higher
	apiNonce := uint64(100)
	backend := &contractBackend{ContractBackend: mockBackend, nonce: uint64(apiNonce + maxNonceDrift)}
	mockBackend.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(apiNonce, nil).Times(1)

	// When the nonce is requested from the backend
	txNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	// Then it should reset the local nonce and return the server nonce with no errors
	r.NoError(err)
	r.Equal(apiNonce, backend.nonce)
	r.Equal(apiNonce, txNonce)
}

func TestSuggestGasPrice_Default(t *testing.T) {
	r := require.New(t)

	suggested := big.NewInt(100)

	// Given no max price
	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))
	backend := &contractBackend{ContractBackend: mockBackend, nonce: 1}
	mockBackend.EXPECT().SuggestGasPrice(gomock.Any()).Return(suggested, nil).Times(1)

	// When the SuggestedGasPrice is called
	res, err := backend.SuggestGasPrice(context.Background())

	// Then it should default to the suggested + 10%
	r.NoError(err)
	r.Equal(int64(110), res.Int64())
}

func TestSuggestGasPrice_MaxExceeded(t *testing.T) {
	r := require.New(t)

	maxPrice := big.NewInt(50)
	suggested := big.NewInt(100)

	// Given no max price
	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))
	backend := &contractBackend{ContractBackend: mockBackend, nonce: 1, maxPrice: maxPrice}
	mockBackend.EXPECT().SuggestGasPrice(gomock.Any()).Return(suggested, nil).Times(1)

	// When the SuggestedGasPrice is called
	res, err := backend.SuggestGasPrice(context.Background())

	// Then it should default to maxPrice
	r.NoError(err)
	r.Equal(int64(50), res.Int64())
}

func TestSuggestGasPrice_MaxNotExceeded(t *testing.T) {
	r := require.New(t)

	maxPrice := big.NewInt(150)
	suggested := big.NewInt(100)

	// Given no max price
	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))
	backend := &contractBackend{ContractBackend: mockBackend, nonce: 1, maxPrice: maxPrice}
	mockBackend.EXPECT().SuggestGasPrice(gomock.Any()).Return(suggested, nil).Times(1)

	// When the SuggestedGasPrice is called
	res, err := backend.SuggestGasPrice(context.Background())

	// Then it should default to suggested + 10%
	r.NoError(err)
	r.Equal(int64(110), res.Int64())
}
