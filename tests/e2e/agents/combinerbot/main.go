package main

import (
	"context"
	"fmt"
	"log"
	"net"

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

	protocol.RegisterAgentServer(
		server, &agentServer{},
	)

	log.Println("Starting agent server...")
	log.Println(server.Serve(lis))
}

type agentServer struct {
	protocol.UnimplementedAgentServer
}

var (
	subscription = "0x5e13c2f3a97c292695b598090056ba5d52f9dcc7790bcdaa8b6cd87c1a1ebc0f"
	// alertSubscriptions subscribes to police bot alerts
	alertSubscriptions = []*protocol.CombinerBotSubscription{
		{
			BotId: subscription,
		},
	}
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
func (as *agentServer) EvaluateAlert(ctx context.Context, request *protocol.EvaluateAlertRequest) (*protocol.EvaluateAlertResponse, error) {
	response := &protocol.EvaluateAlertResponse{Status: protocol.ResponseStatus_SUCCESS}

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
			RelatedAlerts: []string{subscription},
		},
	)

	logrus.WithField("alert", "combiner alert").Warn(response.Findings)

	return response, nil
}
