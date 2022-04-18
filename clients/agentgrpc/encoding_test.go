package agentgrpc_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/config"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var txMsg = &protocol.EvaluateTxRequest{
	RequestId: "123",
	Event: &protocol.TransactionEvent{
		Type: protocol.TransactionEvent_BLOCK,
		Transaction: &protocol.TransactionEvent_EthTransaction{
			Hash: "0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50",
		},
	},
}

type agentServer struct {
	r                *require.Assertions
	doneCh           chan struct{}
	disableAssertion bool
	protocol.UnimplementedAgentServer
}

func (as *agentServer) Initialize(context.Context, *protocol.InitializeRequest) (*protocol.InitializeResponse, error) {
	return &protocol.InitializeResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func (as *agentServer) EvaluateTx(ctx context.Context, txRequest *protocol.EvaluateTxRequest) (*protocol.EvaluateTxResponse, error) {
	if !as.disableAssertion {
		as.r.Equal(txMsg.RequestId, txRequest.RequestId)
		as.r.Equal(txMsg.Event.Transaction.Hash, txRequest.Event.Transaction.Hash)
		close(as.doneCh)
	}
	return &protocol.EvaluateTxResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func (as *agentServer) EvaluateBlock(context.Context, *protocol.EvaluateBlockRequest) (*protocol.EvaluateBlockResponse, error) {
	return &protocol.EvaluateBlockResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func TestEncodeMessage(t *testing.T) {
	r := require.New(t)

	preparedMsg, err := agentgrpc.EncodeMessage(txMsg)
	r.NoError(err)
	log.Printf("%+v", preparedMsg)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.AgentGrpcPort))
	r.NoError(err)
	defer lis.Close()

	server := grpc.NewServer()
	as := &agentServer{r: r, doneCh: make(chan struct{})}
	protocol.RegisterAgentServer(server, as)
	go server.Serve(lis)

	agentClient := agentgrpc.NewClient()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", config.AgentGrpcPort), grpc.WithInsecure())
	r.NoError(err)
	agentClient.WithConn(conn)

	var resp protocol.EvaluateTxResponse
	r.NoError(agentClient.Invoke(context.Background(), agentgrpc.MethodEvaluateTx, preparedMsg, &resp))
	<-as.doneCh
}
