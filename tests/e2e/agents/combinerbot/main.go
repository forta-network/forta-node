package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/forta-network/forta-core-go/clients/graphql"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/bothttp"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/tests/e2e/agents/combinerbot/combinerbotalertid"
	"github.com/gorilla/mux"
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

	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/health", HandleHealthCheck)

		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", config.DefaultBotHealthCheckPort), r)
		if err != nil {
			panic(fmt.Errorf("error while listening health check %w", err))
		}
	}()

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

func (as *agentServer) HealthCheck(context.Context, *protocol.HealthCheckRequest) (*protocol.HealthCheckResponse, error) {
	logrus.Infof("requesting to subscribe bot alerts: %s", alertSubscriptions)
	return &protocol.HealthCheckResponse{
		Status: protocol.HealthCheckResponse_SUCCESS,
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
	logrus.WithField("source", "handle alert feed").Infof("incoming alert %s for bot %s", request.Event.Alert.Hash, subscribedBot)
	response := &protocol.EvaluateAlertResponse{Status: protocol.ResponseStatus_SUCCESS}

	alerts, err := queryPublicAPI(ctx, subscribedBot)
	if err != nil {
		logrus.WithError(err).Warn("failed to fetch latest alerts")
		return &protocol.EvaluateAlertResponse{Status: protocol.ResponseStatus_ERROR}, err
	}

	logrus.WithField("source", "public api proxy").Infof("succesfully fetched %d alerts for bot %s", len(alerts), subscribedBot)

	response.Findings = append(
		response.Findings, &protocol.Finding{
			Protocol:      "1",
			Severity:      protocol.Finding_CRITICAL,
			Metadata:      nil,
			Type:          protocol.Finding_INFORMATION,
			AlertId:       combinerbotalertid.CombinationAlertID,
			Name:          "Combination Alert",
			Description:   request.Event.Alert.Hash,
			Private:       false,
			Addresses:     nil,
			Indicators:    nil,
			RelatedAlerts: []string{subscribedBot},
		},
	)

	logrus.WithField("alert", "combiner alert").Warn(response.Findings)

	return response, nil
}

func queryPublicAPI(ctx context.Context, bot string) ([]*protocol.AlertEvent, error) {
	publicAPIAddr := fmt.Sprintf(
		"http://%s:%s/graphql", os.Getenv(config.EnvPublicAPIProxyHost), os.Getenv(config.EnvPublicAPIProxyPort),
	)
	graphqlClient := graphql.NewClient(publicAPIAddr)

	return graphqlClient.GetAlertsBatch(ctx, []*graphql.AlertsInput{{Bots: []string{bot}}}, nil)
}

func HandleHealthCheck(rw http.ResponseWriter, r *http.Request) {
	healthResponse := bothttp.HealthResponse{
		Metrics: []bothttp.Metrics{
			{
				ChainID: 1,
				DataPoints: map[string][]float64{
					domain.MetricBlockDrop: {1, 2, 3},
				},
			},
			{
				ChainID: 2,
				DataPoints: map[string][]float64{
					domain.MetricBlockDrop: {2},
				},
			},
		},
	}

	err := json.NewEncoder(rw).Encode(healthResponse)
	if err != nil {
		logrus.WithError(err).Warn("can't encode health response")
	}
}
