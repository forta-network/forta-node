package query

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"time"

	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"
	"OpenZeppelin/fortify-node/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	defaultInterval   = time.Second * 15
	defaultBatchLimit = 500
)

// AlertListener allows retrieval of alerts from the database
type AlertListener struct {
	protocol.UnimplementedQueryNodeServer
	ctx       context.Context
	store     store.AlertStore
	cfg       AlertListenerConfig
	contract  AlertsContract
	ipfs      IPFS
	ethClient EthClient

	skipEmpty     bool
	batchInterval time.Duration
	batchLimit    int
	latestBlock   uint64
	latestChainID uint64
	notifCh       chan *protocol.NotifyRequest
}

// EthClient interacts with the Ethereum API.
type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

// AlertsContract stores alerts.
type AlertsContract interface {
	AddAlertBatch(_chainId *big.Int, _blockStart *big.Int, _blockEnd *big.Int, _alertCount *big.Int, _maxSeverity *big.Int, _ref string) (*types.Transaction, error)
}

// IPFS interacts with an IPFS node/gateway.
type IPFS interface {
	Add(r io.Reader, options ...ipfsapi.AddOpts) (string, error)
}

type AlertListenerConfig struct {
	Port            int
	PrivateKey      *ecdsa.PrivateKey
	PublisherConfig config.PublisherConfig
}

func (al *AlertListener) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	al.notifCh <- req
	return &protocol.NotifyResponse{}, nil
}

func (al *AlertListener) publishAlerts() {
	ticker := time.NewTicker(al.batchInterval)
	var err error
	for {
		if err != nil {
			log.Errorf("failed to publish alert batch: %v", err)
			// Sleep
			ticker.Reset(al.batchInterval)
			<-ticker.C
		}

		batch := al.getLatestBatch()

		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(batch); err != nil {
			err = fmt.Errorf("failed to encode the signed alert: %v", err)
			continue
		}
		cid, err := al.ipfs.Add(&buf, ipfsapi.Pin(true))
		if err != nil {
			err = fmt.Errorf("failed to store alert data to ipfs: %v", err)
			continue
		}

		tx, err := al.contract.AddAlertBatch(
			big.NewInt(0).SetUint64(batch.ChainID),
			big.NewInt(0).SetUint64(batch.BlockStart),
			big.NewInt(0).SetUint64(batch.BlockEnd),
			big.NewInt(int64(batch.AlertCount)),
			big.NewInt(0).SetUint64(batch.MaxSeverity),
			cid,
		)
		if err != nil {
			err = fmt.Errorf("failed to send the alert tx: %v", err)
			continue
		}

		for _, alert := range batch.Alerts {
			alert.TxHash = tx.Hash().Hex()
			if err = al.store.AddAlert(alert); err != nil {
				continue
			}
		}

		<-ticker.C
	}
}

// AlertBatch contains a batch of alerts along with some extra data.
type AlertBatch struct {
	ChainID     uint64                  `json:"chainId"`
	BlockStart  uint64                  `json:"blockStart"`
	BlockEnd    uint64                  `json:"blockEnd"`
	AlertCount  int                     `json:"alertCount"`
	MaxSeverity uint64                  `json:"maxSeverity"`
	Alerts      []*protocol.SignedAlert `json:"alerts"`
}

func (al *AlertListener) getLatestBatch() (batch AlertBatch) {
	var done bool
	for i := 0; i < defaultBatchLimit; i++ {
		select {
		case notif := <-al.notifCh:
			alert := notif.SignedAlert
			log.Infof("alert: %s", alert.Alert.Id)

			// TODO: Separate batches by chain ID later?
			chainID, err := hexutil.DecodeUint64(alert.ChainId)
			if err != nil {
				log.Errorf("failed to parse alert chain id: %v", err)
				continue
			}
			batch.ChainID = chainID

			alertBlockNum, err := hexutil.DecodeUint64(alert.BlockNumber)
			if err != nil {
				log.Errorf("failed to parse alert block number: %v", err)
				continue
			}
			if batch.BlockStart == 0 || (batch.BlockStart > 0 && alertBlockNum < batch.BlockStart) {
				batch.BlockStart = alertBlockNum
			}
			if batch.BlockEnd == 0 || (batch.BlockEnd > 0 && alertBlockNum > batch.BlockEnd) {
				batch.BlockEnd = alertBlockNum
			}

			if uint64(alert.Alert.Finding.Severity) > batch.MaxSeverity {
				batch.MaxSeverity = uint64(alert.Alert.Finding.Severity)
			}

			batch.Alerts = append(batch.Alerts, alert)

		default:
			done = true // If we don't receive anymore notifs
		}
		if done {
			break
		}
	}

	batch.AlertCount = len(batch.Alerts)

	// We use single chain ID for now.
	if batch.ChainID > 0 {
		al.latestChainID = batch.ChainID
	} else {
		batch.ChainID = al.latestChainID
	}

	if batch.BlockStart == 0 {
		latestBlock, err := al.ethClient.BlockNumber(al.ctx)
		if err != nil {
			log.Errorf("failed to get the latest block for batch: %v", err)
			return
		}
		batch.BlockStart = al.latestBlock
		batch.BlockEnd = latestBlock
		al.latestBlock = latestBlock
	}

	return
}

func (al *AlertListener) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:8770")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterQueryNodeServer(grpcServer, al)

	go al.publishAlerts()

	return grpcServer.Serve(lis)
}

func (al *AlertListener) Stop() error {
	log.Infof("Stopping %s", al.Name())
	return nil
}

func (al *AlertListener) Name() string {
	return "AlertListener"
}

func NewAlertListener(ctx context.Context, store store.AlertStore, cfg AlertListenerConfig) (*AlertListener, error) {
	rpcClient, err := rpc.Dial(cfg.PublisherConfig.Ethereum.JsonRpcUrl)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	txOpts := bind.NewKeyedTransactor(cfg.PrivateKey)
	contract, err := contracts.NewAlertsTransactor(common.HexToAddress(cfg.PublisherConfig.ContractAddress), ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the alerts contract: %v", err)
	}
	ats := &contracts.AlertsTransactorSession{
		Contract:     contract,
		TransactOpts: *txOpts,
	}

	var ipfsClient IPFS
	if len(cfg.PublisherConfig.IPFS.Username) > 0 && len(cfg.PublisherConfig.IPFS.Password) > 0 {
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.Ethereum.JsonRpcUrl, &http.Client{
			Transport: utils.NewBasicAuthTransport(cfg.PublisherConfig.IPFS.Username, cfg.PublisherConfig.IPFS.Password),
		})
	} else {
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.Ethereum.JsonRpcUrl, http.DefaultClient)
	}

	latestBlock, err := ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest block: %v", err)
	}

	batchInterval := defaultInterval
	if cfg.PublisherConfig.Batch.IntervalSeconds != nil {
		batchInterval = (time.Duration)(*cfg.PublisherConfig.Batch.IntervalSeconds) * time.Second
	}

	batchLimit := defaultBatchLimit
	if cfg.PublisherConfig.Batch.MaxAlerts != nil {
		batchLimit = *cfg.PublisherConfig.Batch.MaxAlerts
	}

	return &AlertListener{
		ctx:       ctx,
		store:     store,
		cfg:       cfg,
		contract:  ats,
		ipfs:      ipfsClient,
		ethClient: ethClient,

		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		batchInterval: batchInterval,
		batchLimit:    batchLimit,
		latestBlock:   latestBlock,
		notifCh:       make(chan *protocol.NotifyRequest, batchLimit),
	}, nil
}
