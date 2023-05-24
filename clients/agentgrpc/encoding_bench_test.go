package agentgrpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/config"
	"google.golang.org/grpc"
)

var (
	benchBlockMsg = &protocol.EvaluateBlockRequest{}
	benchTxMsg    = &protocol.EvaluateTxRequest{}
)

func init() {
	b, err := ioutil.ReadFile("./testdata/bench_block.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &benchBlockMsg.Event); err != nil {
		panic(err)
	}
	b, err = ioutil.ReadFile("./testdata/bench_tx.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &benchTxMsg.Event); err != nil {
		panic(err)
	}
}

const benchAgentReqCount = 25

func getBenchClient() agentgrpc.Client {
	agentClient := agentgrpc.NewClient()
	for {
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", config.AgentGrpcPort), grpc.WithInsecure())
		if err == nil {
			agentClient.WithConn(conn)
			var success bool
			_, err1 := agentClient.EvaluateBlock(context.Background(), benchBlockMsg)
			_, err2 := agentClient.EvaluateTx(context.Background(), benchTxMsg)
			success = (err1 == nil) && (err2 == nil)
			if success {
				break
			}
		}
		time.Sleep(time.Second * 2)
		log.Println("retrying to connect to grpc server")
	}

	return agentClient
}

func BenchmarkEvaluateBlock(b *testing.B) {
	agentClient := getBenchClient()
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchAgentReqCount; j++ {
			out, err := agentClient.EvaluateBlock(context.Background(), benchBlockMsg)
			if err != nil {
				panic(err)
			}
			_ = out
		}
	}
}

func BenchmarkEvaluateBlockWithPreparedMessage(b *testing.B) {
	agentClient := getBenchClient()
	for i := 0; i < b.N; i++ {
		preparedMsg, err := agentgrpc.EncodeMessage(benchBlockMsg)
		if err != nil {
			panic(err)
		}
		for j := 0; j < benchAgentReqCount; j++ {
			var resp protocol.EvaluateBlockResponse
			err := agentClient.Invoke(context.Background(), agentgrpc.MethodEvaluateBlock, preparedMsg, &resp)
			if err != nil {
				panic(err)
			}
		}
	}
}

func BenchmarkEvaluateTx(b *testing.B) {
	agentClient := getBenchClient()
	for i := 0; i < b.N; i++ {
		for j := 0; j < benchAgentReqCount; j++ {
			out, err := agentClient.EvaluateTx(context.Background(), benchTxMsg)
			if err != nil {
				panic(err)
			}
			_ = out
		}
	}
}

func BenchmarkEvaluateTxWithPreparedMessage(b *testing.B) {
	agentClient := getBenchClient()
	for i := 0; i < b.N; i++ {
		preparedMsg, err := agentgrpc.EncodeMessage(benchTxMsg)
		if err != nil {
			panic(err)
		}
		for j := 0; j < benchAgentReqCount; j++ {
			var resp protocol.EvaluateTxResponse
			err := agentClient.Invoke(context.Background(), agentgrpc.MethodEvaluateTx, preparedMsg, &resp)
			if err != nil {
				panic(err)
			}
		}
	}
}
