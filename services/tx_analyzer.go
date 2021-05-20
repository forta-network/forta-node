package services

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/feeds"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/store"
)

const keyPath = "/.keys"

// TxAnalyzerService reads TX info, calls agents, and emits results
type TxAnalyzerService struct {
	queryNode protocol.QueryNodeClient
	cfg       TxAnalyzerServiceConfig
	ctx       context.Context
	agents    []AnalyzerAgent
	key       *keystore.Key
}

type AnalyzerAgent struct {
	config config.AgentConfig
	client protocol.AgentClient
}

type TxAnalyzerServiceConfig struct {
	TxChannel     <-chan *feeds.TransactionEvent
	AgentConfigs  []config.AgentConfig
	QueryNodeAddr string
}

type responseWrapper struct {
	request  *protocol.EvaluateRequest
	agent    AnalyzerAgent
	response *protocol.EvaluateResponse
}

// newAgentStream creates a agent transaction handler (sends and receives request)
func newAgentStream(ctx context.Context, agent AnalyzerAgent, input <-chan *protocol.EvaluateRequest, output chan<- *responseWrapper) func() error {
	return func() error {
		for request := range input {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			resp, err := agent.client.Evaluate(ctx, request)
			cancel()
			if err != nil {
				log.Error("error invoking agent", err)
				continue
			}
			resp.Metadata["image"] = agent.config.Image
			output <- &responseWrapper{
				agent:    agent,
				response: resp,
				request:  request,
			}

			//TODO: this print logic is just for test purposes
			m := jsonpb.Marshaler{}
			resStr, err := m.MarshalToString(resp)
			if err != nil {
				log.Error("error marshaling response", err)
				continue
			}
			log.Infof(resStr)
		}
		return nil
	}
}

func loadKey() (*keystore.Key, error) {
	f, err := os.OpenFile("/passphrase", os.O_RDONLY, 400)
	if err != nil {
		return nil, err
	}

	pw, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	passphrase := string(pw)

	files, err := ioutil.ReadDir(keyPath)
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		return nil, errors.New("there must be only one key in key directory")
	}

	keyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", keyPath, files[0].Name()))
	if err != nil {
		return nil, err
	}

	return keystore.DecryptKey(keyBytes, passphrase)
}

func (t *TxAnalyzerService) sign(alert *protocol.Alert) (*protocol.SignedAlert, error) {
	b, err := proto.Marshal(alert)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256(b)
	sig, err := crypto.Sign(hash, t.key.PrivateKey)
	if err != nil {
		return nil, err
	}
	signature := fmt.Sprintf("0x%s", hex.EncodeToString(sig))
	return &protocol.SignedAlert{
		Alert: alert,
		Signature: &protocol.Signature{
			Signature: signature,
			Algorithm: "ECDSA",
		},
	}, nil
}

func (t *TxAnalyzerService) Start() error {
	log.Infof("Starting %s", t.Name())
	grp, ctx := errgroup.WithContext(t.ctx)

	//TODO: change this protocol when we know more about query-node delivery
	// Gear 3: handles responses (sends to query node)
	output := make(chan *responseWrapper, 100)
	grp.Go(func() error {
		for resp := range output {
			ts := time.Now().UTC()

			//TODO: validate finding returned is well-formed
			for _, f := range resp.response.Findings {
				b, err := proto.Marshal(f)
				if err != nil {
					return err
				}
				alertID := base58.Encode(sha3.New256().Sum(b))
				r := resp.request.Event.Receipt
				alert := &protocol.Alert{
					Id:        alertID,
					Finding:   f,
					Timestamp: ts.Format(store.AlertTimeFormat),
					Type:      protocol.AlertType_TRANSACTION,
					Agent: &protocol.AgentInfo{
						Name:      resp.agent.config.Name,
						Image:     resp.agent.config.Image,
						ImageHash: resp.agent.config.ImageHash,
					},
					Tags: map[string]string{
						"chainId":     resp.request.Event.Network.ChainId,
						"blockHash":   r.BlockHash,
						"blockNumber": r.BlockNumber,
						"txHash":      r.TransactionHash,
					},
					Scanner: &protocol.ScannerInfo{
						Address: t.key.Address.Hex(),
					},
				}

				signedAlert, err := t.sign(alert)
				if err != nil {
					return err
				}
				_, err = t.queryNode.Notify(ctx, &protocol.NotifyRequest{
					SignedAlert: signedAlert,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	// Gear 2: set of agents pulling from their own individual channels
	var agentChannels []chan *protocol.EvaluateRequest
	for _, agt := range t.agents {
		agent := agt
		input := make(chan *protocol.EvaluateRequest, 100)
		agentChannels = append(agentChannels, input)
		grp.Go(newAgentStream(ctx, agent, input, output))
	}

	// Gear 1: loops over transactions and distributes to all agents
	grp.Go(func() error {
		defer func() {
			for _, agtCh := range agentChannels {
				close(agtCh)
			}
			close(output)
		}()

		// for each transaction
		for tx := range t.cfg.TxChannel {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// convert to message
			msg, err := tx.ToMessage()
			if err != nil {
				log.Error("error converting tx event to message", err)
				continue
			}

			// create a request
			requestId := uuid.Must(uuid.NewUUID())
			request := &protocol.EvaluateRequest{RequestId: requestId.String(), Event: msg}

			// forward to each agent channel
			for _, agtCh := range agentChannels {
				agtCh <- request
			}
		}
		return nil
	})

	return grp.Wait()
}

func (t *TxAnalyzerService) Stop() error {
	log.Infof("Stopping %s", t.Name())
	return nil
}

func (t *TxAnalyzerService) Name() string {
	return "TxAnalyzerService"
}

func NewTxAnalyzerService(ctx context.Context, cfg TxAnalyzerServiceConfig) (*TxAnalyzerService, error) {
	var agents []AnalyzerAgent
	for _, agt := range cfg.AgentConfigs {
		conn, err := grpc.Dial(fmt.Sprintf("%s:50051", agt.ContainerName()), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect to %s, %v", agt.ContainerName(), err)
		}
		client := protocol.NewAgentClient(conn)
		agents = append(agents, AnalyzerAgent{
			config: agt,
			client: client,
		})
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:8770", cfg.QueryNodeAddr), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	qn := protocol.NewQueryNodeClient(conn)

	k, err := loadKey()
	if err != nil {
		return nil, err
	}

	log.Infof("loaded address %s", k.Address.Hex())

	return &TxAnalyzerService{
		cfg:       cfg,
		ctx:       ctx,
		agents:    agents,
		queryNode: qn,
		key:       k,
	}, nil
}
