package query

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"
	"OpenZeppelin/fortify-node/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// AlertListener allows retrieval of alerts from the database
type AlertListener struct {
	protocol.UnimplementedQueryNodeServer
	ctx      context.Context
	store    store.AlertStore
	cfg      AlertListenerConfig
	contract AlertsContract
	ipfs     IPFS

	notifCh chan *protocol.NotifyRequest
}

// AlertsContract stores alerts.
type AlertsContract interface {
	AddAlert(_poolId [32]byte, _agentId [32]byte, _alertId [32]byte, _alertRef string) (*types.Transaction, error)
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
	notif := <-al.notifCh
	var err error
	for {
		if err != nil {
			log.Errorf("failed to publish alert '%s': %v", notif.SignedAlert.Alert.Id)
			time.Sleep(time.Second * 10)
		}

		log.Infof("alert: %s", notif.SignedAlert.Alert.Id)

		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(notif.SignedAlert); err != nil {
			err = fmt.Errorf("failed to encode the signed alert: %v", err)
			continue
		}
		cid, err := al.ipfs.Add(&buf, ipfsapi.Pin(true))
		if err != nil {
			err = fmt.Errorf("failed to store alert data to ipfs: %v", err)
			continue
		}

		// TODO: We won't have pools. Update/remove pool ID argument.
		tx, err := al.contract.AddAlert(([32]byte)(common.Hash{}), ([32]byte)(common.HexToHash(notif.SignedAlert.Alert.Agent.Name)), common.HexToHash(notif.SignedAlert.Alert.Id), cid)
		if err != nil {
			err = fmt.Errorf("failed to send the alert tx: %v", err)
			continue
		}
		notif.SignedAlert.TxHash = tx.Hash().Hex()
		if err = al.store.AddAlert(notif.SignedAlert); err != nil {
			continue
		}
	}
}

func (al *AlertListener) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:8770")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterQueryNodeServer(grpcServer, al)
	return grpcServer.Serve(lis)
}

func (al *AlertListener) Stop() error {
	log.Infof("Stopping %s", al.Name())
	return nil
}

func (al *AlertListener) Name() string {
	return "AlertListener"
}

func NewAlertListener(ctx context.Context, store store.AlertStore, cfg AlertListenerConfig, msgClient clients.MessageClient) (*AlertListener, error) {
	rpcClient, err := rpc.Dial(cfg.PublisherConfig.JSONRPCURL)
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
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.JSONRPCURL, &http.Client{
			Transport: utils.NewBasicAuthTransport(cfg.PublisherConfig.IPFS.Username, cfg.PublisherConfig.IPFS.Password),
		})
	} else {
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.JSONRPCURL, http.DefaultClient)
	}

	return &AlertListener{
		ctx:      ctx,
		store:    store,
		cfg:      cfg,
		contract: ats,
		ipfs:     ipfsClient,
	}, nil
}
