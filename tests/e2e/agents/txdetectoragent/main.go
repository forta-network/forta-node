package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
	jwt_provider "github.com/forta-network/forta-node/services/jwt-provider"
	"github.com/forta-network/forta-node/tests/e2e/agents/txdetectoragent/testbotalertid"
	"github.com/forta-network/forta-node/tests/e2e/ethaccounts"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.AgentGrpcPort))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	ethClient, err := ethclient.Dial(
		fmt.Sprintf("http://%s:%s", os.Getenv(config.EnvJsonRpcHost), os.Getenv(config.EnvJsonRpcPort)),
	)
	if err != nil {
		panic(err)
	}
	protocol.RegisterAgentServer(server, &agentServer{
		ethClient: ethClient,
	})

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			log.Println("new log", time.Now().UnixNano())
		}
	}()

	log.Println("Starting agent server...")
	log.Println(server.Serve(lis))
}

type agentServer struct {
	ethClient *ethclient.Client
	protocol.UnimplementedAgentServer
}

func (as *agentServer) Initialize(context.Context, *protocol.InitializeRequest) (*protocol.InitializeResponse, error) {
	return &protocol.InitializeResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func (as *agentServer) EvaluateTx(ctx context.Context, txRequest *protocol.EvaluateTxRequest) (*protocol.EvaluateTxResponse, error) {
	response := &protocol.EvaluateTxResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}

	token, err := fetchJWTToken()
	if err != nil {
		logrus.WithError(err).Warn("can't fetch token")
		return nil, err
	}

	// detect exploiter address transactions
	if strings.EqualFold(txRequest.Event.Transaction.From, ethaccounts.ExploiterAddress.Hex()) {
		balance, err := as.ethClient.BalanceAt(context.Background(), ethaccounts.ExploiterAddress, nil)
		if err != nil {
			return &protocol.EvaluateTxResponse{
				Status: protocol.ResponseStatus_ERROR,
				Errors: []*protocol.Error{
					{
						Message: fmt.Sprintf("failed to get balance: %v", err),
					},
				},
			}, nil
		}
		response.Findings = []*protocol.Finding{
			{
				Protocol:    "testchain",
				Severity:    protocol.Finding_CRITICAL,
				AlertId:     testbotalertid.ExploiterAlertId,
				Name:        "Exploiter Transaction Detected",
				Description: txRequest.Event.Receipt.TransactionHash,
				Metadata: map[string]string{
					"exploiter": ethaccounts.ExploiterAddress.Hex(),
					"balance":   balance.String(),
				},
			},
			{
				Protocol:    "testchain",
				Severity:    protocol.Finding_INFO,
				AlertId:     testbotalertid.TokenAlertId,
				Name:        "Scanner Token Retrieved",
				Description: txRequest.Event.Receipt.TransactionHash,
				Metadata: map[string]string{
					"token": token,
				},
			},
		}
	}
	return response, nil
}

func (as *agentServer) EvaluateBlock(context.Context, *protocol.EvaluateBlockRequest) (*protocol.EvaluateBlockResponse, error) {
	return &protocol.EvaluateBlockResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func fetchJWTToken() (string, error) {
	jwtProviderAddr := fmt.Sprintf(
		"%s:%s", os.Getenv(config.EnvJWTProviderHost), os.Getenv(config.EnvJWTProviderPort),
	)

	payload, err := json.Marshal(
		jwt_provider.CreateJWTMessage{
			Claims: map[string]interface{}{
				"data":    "123",
				"payload": 123,
			},
		},
	)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/create", jwtProviderAddr), "application/json", bytes.NewReader(payload),
	)
	if err != nil {
		logrus.WithError(err).Warn("can not fetch jwt token")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		reason, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("can't fetch jwt, status code: %d, reason: %s ", resp.StatusCode, string(reason))
	}

	var s jwt_provider.CreateJWTResponse
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", err
	}

	return s.Token, nil
}
