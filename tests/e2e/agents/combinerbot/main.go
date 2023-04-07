package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/forta-network/forta-core-go/clients/graphql"
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
	subscribedBot = "0xbe1872858e63b6ed4ef7b84fc453970dc8d89968715797662a4f43c01d598aab"
	// alertSubscriptions subscribes to police bot alerts
	alertSubscriptions = []*protocol.CombinerBotSubscription{
		{
			BotId: subscribedBot,
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
	logrus.Infof("evaluating alert %s", request.Event.Alert.Hash)
	response := &protocol.EvaluateAlertResponse{Status: protocol.ResponseStatus_SUCCESS}

	err := queryPublicAPI(ctx)
	if err != nil {
		logrus.WithError(err).Warn("failed to fetch latest alerts")
		return &protocol.EvaluateAlertResponse{Status: protocol.ResponseStatus_ERROR}, err
	}

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
			RelatedAlerts: []string{subscribedBot},
		},
	)

	logrus.WithField("alert", "combiner alert").Warn(response.Findings)

	return response, nil
}

func queryPublicAPI(ctx context.Context) error {
	publicAPIAddr := fmt.Sprintf(
		"http://%s:%s/graphql", os.Getenv(config.EnvPublicAPIProxyHost), os.Getenv(config.EnvPublicAPIProxyPort),
	)
	graphqlClient := graphql.NewClient(publicAPIAddr)

	_, err := graphqlClient.GetAlerts(ctx, &graphql.AlertsInput{Bots: []string{subscribedBot}}, nil)

	return err
}
