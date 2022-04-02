package agentgrpc_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-node/clients/agentgrpc"
	"github.com/forta-protocol/forta-node/config"
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
	r      *require.Assertions
	doneCh chan struct{}
	protocol.UnimplementedAgentServer
}

func (as *agentServer) Initialize(context.Context, *protocol.InitializeRequest) (*protocol.InitializeResponse, error) {
	return &protocol.InitializeResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func (as *agentServer) EvaluateTx(ctx context.Context, txRequest *protocol.EvaluateTxRequest) (*protocol.EvaluateTxResponse, error) {
	as.r.Equal(txMsg.RequestId, txRequest.RequestId)
	as.r.Equal(txMsg.Event.Transaction.Hash, txRequest.Event.Transaction.Hash)
	close(as.doneCh)
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

var benchBlockMsg = &protocol.EvaluateBlockRequest{
	RequestId: "123",
	Event: &protocol.BlockEvent{
		BlockNumber: "0x123",
		Block:       &protocol.BlockEvent_EthBlock{},
	},
}

var preparedBlockMsg *grpc.PreparedMsg

func init() {
	for i := 0; i < 10000; i++ {
		benchBlockMsg.Event.Block.Uncles = append(benchBlockMsg.Event.Block.Uncles, "0xf779d223db50593a463e4c73cdfc2b5aa4d780cfceb39288662027d8df061ab4")
	}
	preparedMsg, err := agentgrpc.EncodeMessage(benchBlockMsg)
	if err != nil {
		panic(err)
	}
	preparedBlockMsg = preparedMsg
}

const benchAgentReqCount = 25

func BenchmarkEvaluateBlock(b *testing.B) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.AgentGrpcPort))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	as := &agentServer{}
	protocol.RegisterAgentServer(server, as)
	go server.Serve(lis)

	time.Sleep(time.Second * 10)

	agentClient := agentgrpc.NewClient()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", config.AgentGrpcPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	agentClient.WithConn(conn)

	ch := make(chan struct{}, benchAgentReqCount)
	for i := 0; i < benchAgentReqCount; i++ {
		go func() {
			for range ch {
				out, err := agentClient.EvaluateBlock(context.Background(), benchBlockMsg)
				if err != nil {
					panic(err)
				}
				_ = out
			}
		}()
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < benchAgentReqCount; j++ {
			ch <- struct{}{}
		}
	}
	close(ch)
}

func BenchmarkEvaluateBlockWithPreparedMessage(b *testing.B) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.AgentGrpcPort))
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	as := &agentServer{}
	protocol.RegisterAgentServer(server, as)
	go server.Serve(lis)

	time.Sleep(time.Second * 10)

	agentClient := agentgrpc.NewClient()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", config.AgentGrpcPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	agentClient.WithConn(conn)

	ch := make(chan struct{}, benchAgentReqCount)
	for i := 0; i < benchAgentReqCount; i++ {
		go func() {
			for range ch {
				var resp protocol.EvaluateBlockResponse
				err := agentClient.Invoke(context.Background(), agentgrpc.MethodEvaluateBlock, preparedBlockMsg, &resp)
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < benchAgentReqCount; j++ {
			ch <- struct{}{}
		}
	}
	close(ch)
}
