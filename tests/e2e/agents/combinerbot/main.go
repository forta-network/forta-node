package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/tests/e2e/agents/combinerbot/combinerbotalertid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.AgentGrpcPort))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()

	t := time.NewTicker(time.Minute)
	protocol.RegisterAgentServer(
		server, &agentServer{
			ticker: t,
		},
	)

	log.Println("Starting agent server...")
	log.Println(server.Serve(lis))
}

type agentServer struct {
	protocol.UnimplementedAgentServer
	ticker *time.Ticker
}

var (
	// alertSubscriptions subscribes to police bot alerts
	alertSubscriptions = []string{"0x5e13c2f3a97c292695b598090056ba5d52f9dcc7790bcdaa8b6cd87c1a1ebc0f"}
)

func (as *agentServer) Initialize(context.Context, *protocol.InitializeRequest) (*protocol.InitializeResponse, error) {
	logrus.Infof("requesting to subscribe bot alerts: %s", alertSubscriptions)
	return &protocol.InitializeResponse{
		Status: protocol.ResponseStatus_SUCCESS,
		AlertConfig: &protocol.AlertConfig{
			Subscriptions: alertSubscriptions,
		},
	}, nil
}

func (as *agentServer) EvaluateTx(ctx context.Context, txRequest *protocol.EvaluateTxRequest) (*protocol.EvaluateTxResponse, error) {
	response := &protocol.EvaluateTxResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}

	return response, nil
}
func (as *agentServer) EvaluateBlock(ctx context.Context, txRequest *protocol.EvaluateBlockRequest) (*protocol.EvaluateBlockResponse, error) {
	response := &protocol.EvaluateBlockResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}

	return response, nil
}
func (as *agentServer) EvaluateCombination(ctx context.Context, request *protocol.EvaluateCombinationRequest) (*protocol.EvaluateCombinationResponse, error) {
	response := &protocol.EvaluateCombinationResponse{Status: protocol.ResponseStatus_SUCCESS}

	select {
	case <-as.ticker.C:
		response.Findings = append(
			response.Findings, &protocol.Finding{
				Protocol:      "1",
				Severity:      protocol.Finding_CRITICAL,
				Metadata:      nil,
				Type:          protocol.Finding_INFORMATION,
				AlertId:       combinerbotalertid.CombinationAlertID,
				Name:          "Combination Alert",
				Description:   request.Event.Alert.Hash,
				EverestId:     "",
				Private:       false,
				Addresses:     nil,
				Indicators:    nil,
				RelatedAlerts: alertSubscriptions,
			},
		)
	default:

	}

	logrus.WithField("alert", "combiner alert").Warn(response.Findings)

	return response, nil
}
