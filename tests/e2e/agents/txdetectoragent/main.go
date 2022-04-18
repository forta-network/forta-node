package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/tests/e2e/ethaccounts"
	"google.golang.org/grpc"
)

const (
	AlertId = "EXPLOITER_TRANSACTION"
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
				AlertId:     AlertId,
				Name:        "Exploiter Transaction Detected",
				Description: txRequest.Event.Receipt.TransactionHash,
				Metadata: map[string]string{
					"exploiter": ethaccounts.ExploiterAddress.Hex(),
					"balance":   balance.String(),
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
