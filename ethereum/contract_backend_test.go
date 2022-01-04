package ethereum

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

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
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 100}
	apiNonce := backend.localNonce + 1
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
	r.Equal(backend.localNonce, postTxNonce)
	r.Greater(postTxNonce, txNonce)
}

func TestLaggingBackend(t *testing.T) {
	r := require.New(t)

	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))

	// Given that the local nonce is higher
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 100}
	apiNonce := backend.localNonce - 1
	mockBackend.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(apiNonce, nil).Times(2)
	mockBackend.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)

	// When the nonce is requested from the backend
	txNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	// Then it should return the local nonce with no errors
	r.NoError(err)
	r.Equal(backend.localNonce, txNonce)

	// And new transaction should cause a higher local nonce to be returned
	testTx := types.NewTransaction(txNonce, testAddr, big.NewInt(1), 21000, big.NewInt(30), []byte{})
	r.NoError(backend.SendTransaction(context.Background(), testTx))
	postTxNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	r.NoError(err)
	r.Equal(backend.localNonce, postTxNonce)
	r.Greater(postTxNonce, txNonce)
}

func TestDriftingLocal(t *testing.T) {
	r := require.New(t)

	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))

	// Given that the local nonce is higher
	apiNonce := uint64(100)
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: uint64(apiNonce + maxNonceDrift)}
	mockBackend.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(apiNonce, nil).Times(1)

	// When the nonce is requested from the backend
	txNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	// Then it should reset the local nonce and return the server nonce with no errors
	r.NoError(err)
	r.Equal(apiNonce, backend.localNonce)
	r.Equal(apiNonce, txNonce)
}

func TestReplacementErr(t *testing.T) {
	r := require.New(t)

	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))

	// Given that the local nonce is higher
	apiNonce := uint64(100)
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: apiNonce + 1}
	mockBackend.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(apiNonce, nil).Times(1)

	// When the nonce is requested from the backend
	txNonce, err := backend.PendingNonceAt(context.Background(), testAddr)
	// Then it should return the local nonce
	r.NoError(err)
	r.Equal(backend.localNonce, txNonce)
	r.Greater(txNonce, apiNonce)

	// And new transaction should cause a nonce reset when replacement tx error is returned
	testTx := types.NewTransaction(txNonce, testAddr, big.NewInt(1), 21000, big.NewInt(30), []byte{})
	mockBackend.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(errors.New(errStrReplacementTx))
	r.Error(backend.SendTransaction(context.Background(), testTx))
	r.Equal(apiNonce, backend.localNonce)
	r.Greater(txNonce, backend.localNonce)
}

func TestSuggestGasPrice_Default(t *testing.T) {
	r := require.New(t)

	suggested := big.NewInt(100)

	// Given no max price
	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 1}
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
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 1, maxPrice: maxPrice}
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
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 1, maxPrice: maxPrice}
	mockBackend.EXPECT().SuggestGasPrice(gomock.Any()).Return(suggested, nil).Times(1)

	// When the SuggestedGasPrice is called
	res, err := backend.SuggestGasPrice(context.Background())

	// Then it should default to suggested + 10%
	r.NoError(err)
	r.Equal(int64(110), res.Int64())
}

func TestSuggestGasPrice_CacheHit(t *testing.T) {
	r := require.New(t)

	maxPrice := big.NewInt(150)
	cached := big.NewInt(110)

	// Given no max price
	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 1, maxPrice: maxPrice, gasPrice: cached, gasPriceUpdated: time.Now()}

	// When the SuggestedGasPrice is called
	res, err := backend.SuggestGasPrice(context.Background())

	// Then it should default to suggested + 10%
	r.NoError(err)
	r.Equal(int64(110), res.Int64())
}

func TestSuggestGasPrice_CacheMiss(t *testing.T) {
	r := require.New(t)

	cached := big.NewInt(110)
	suggested := big.NewInt(200)

	// Given no max price
	mockBackend := mock_ethereum.NewMockContractBackend(gomock.NewController(t))
	backend := &contractBackend{ContractBackend: mockBackend, localNonce: 1, gasPrice: cached, gasPriceUpdated: time.Now().Add(-5 * time.Minute)}
	mockBackend.EXPECT().SuggestGasPrice(gomock.Any()).Return(suggested, nil).Times(1)

	// When the SuggestedGasPrice is called
	res, err := backend.SuggestGasPrice(context.Background())

	// Then it should default to suggested + 10%
	r.NoError(err)
	r.Equal(int64(220), res.Int64())
}
