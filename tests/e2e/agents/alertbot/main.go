package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/tests/e2e/agents/alertbot/alerttestbotalertid"
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
		server, &agentServer{
		},
	)

	log.Println("Starting agent server...")
	log.Println(server.Serve(lis))
}

type agentServer struct {
	protocol.UnimplementedAgentServer
}

var (
	// alertSubscriptions subscribes to police bot alerts
	alertSubscriptions = []string{"0xe66d22cdcfe0b7e03cbd01e554727fa760aa4170e3d565b7c5a2547f587225ad"}
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

func (as *agentServer) EvaluateAlert(ctx context.Context, alertRequest *protocol.EvaluateAlertRequest) (*protocol.EvaluateAlertResponse, error) {
	response := &protocol.EvaluateAlertResponse{Status: protocol.ResponseStatus_SUCCESS}

	response.Findings = append(
		response.Findings, &protocol.Finding{
			Protocol:    "1",
			Severity:    protocol.Finding_CRITICAL,
			Metadata:    nil,
			Type:        protocol.Finding_INFORMATION,
			AlertId:     alerttestbotalertid.TraceSupportAlertId,
			Name:        "Check Trace Support",
			Description: alertRequest.Event.Alert.Metadata["containerTraceSupported"],
			EverestId:   "",
			Private:     false,
			Addresses:   nil,
			Indicators:  nil,
		},
	)

	logrus.WithField("alert", "trace check").Warn(response.Findings)

	return response, nil
}
