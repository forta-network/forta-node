package main

import (
	"context"
	"fmt"
	"net"

	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-node/config"
	"google.golang.org/grpc"
)

type AgentServer struct {
	protocol.UnimplementedAgentServer
}

func (as *AgentServer) Initialize(context.Context, *protocol.InitializeRequest) (*protocol.InitializeResponse, error) {
	return &protocol.InitializeResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func (as *AgentServer) EvaluateTx(ctx context.Context, txRequest *protocol.EvaluateTxRequest) (*protocol.EvaluateTxResponse, error) {
	return &protocol.EvaluateTxResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func (as *AgentServer) EvaluateBlock(context.Context, *protocol.EvaluateBlockRequest) (*protocol.EvaluateBlockResponse, error) {
	return &protocol.EvaluateBlockResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.AgentGrpcPort))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	as := &AgentServer{}
	protocol.RegisterAgentServer(server, as)
	server.Serve(lis)
}
