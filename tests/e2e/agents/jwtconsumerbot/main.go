package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/config"
	jwt_provider "github.com/forta-network/forta-node/services/jwt-provider"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	externalAPIAddr = "localhost:7474"
	AlertId         = "EXTERNAL_REQUEST"
)

var (
	jwtProviderAddr = fmt.Sprintf(
		"%s:%s", os.Getenv(config.EnvJWTProviderHost), os.Getenv(config.EnvJWTProviderPort),
	)

	log = logrus.WithFields(
		logrus.Fields{
			"container": "jwt-provider",
		},
	)
)

func serverWithJWTAuth() {
	logrus.Infof("starting mock api on %s", externalAPIAddr)
	err := http.ListenAndServe(externalAPIAddr, http.HandlerFunc(jwtAuthHandler))
	if err != nil {
		panic(err)
	}
}

func jwtAuthHandler(writer http.ResponseWriter, request *http.Request) {
	log := log.WithField("source", "external-api")
	token, err := security.VerifyScannerJWT(request.Header.Get("Authorization"))
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Token.Claims.(jwt.MapClaims)
	if !ok {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	botID := claims["bot-id"]

	if botID != os.Getenv(config.EnvFortaBotID) {
		log.Infof("%d \t %s->%s", http.StatusUnauthorized, request.RemoteAddr, botID)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Infof("%d \t %s->%s", http.StatusOK, request.RemoteAddr, botID)
	writer.WriteHeader(http.StatusOK)
}

func main() {
	go serverWithJWTAuth()

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
	protocol.RegisterAgentServer(
		server, &agentServer{
			ethClient: ethClient,
		},
	)

	log.Println("Starting agent server...")
	log.Println(server.Serve(lis))
}

type agentServer struct {
	ethClient *ethclient.Client
	protocol.UnimplementedAgentServer
}

func (as *agentServer) Initialize(context.Context, *protocol.InitializeRequest) (
	*protocol.InitializeResponse, error,
) {
	return &protocol.InitializeResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}, nil
}

func fetchJWTToken() (string, error) {
	payload, err := json.Marshal(
		jwt_provider.CreateJWTMessage{
			Claims: map[string]interface{}{
				"data":    "123",
				"payload": 123,
			},
		},
	)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/create", jwtProviderAddr), "application/json", bytes.NewReader(payload),
	)
	if err != nil {
		logrus.WithError(err).Warn("can not fetch jwt token")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		reason, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("can't fetch jwt, status code: %d, reason: %s ", resp.StatusCode, string(reason))
	}

	var s jwt_provider.CreateJWTResponse
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", err
	}

	return s.Token, nil
}

func (as *agentServer) EvaluateTx(
	_ context.Context, _ *protocol.EvaluateTxRequest,
) (*protocol.EvaluateTxResponse, error) {
	response := &protocol.EvaluateTxResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}

	return response, nil
}

func sendRequestWithJWT(token string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", externalAPIAddr), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to use jwt token")
	}

	return resp.Body, err
}

func (as *agentServer) EvaluateBlock(
	_ context.Context, blockRequest *protocol.EvaluateBlockRequest,
) (*protocol.EvaluateBlockResponse, error) {
	response := &protocol.EvaluateBlockResponse{
		Status: protocol.ResponseStatus_SUCCESS,
	}

	token, err := fetchJWTToken()
	if err != nil {
		logrus.WithError(err).Warn("can't fetch token")
		return nil, err
	}

	_, err = sendRequestWithJWT(token)
	if err != nil {
		return nil, err
	}

	log.Warn("server accepted jwt token")
	response.Findings = append(
		response.Findings, &protocol.Finding{
			Protocol:    "testchain",
			Severity:    protocol.Finding_INFO,
			Metadata:    nil,
			AlertId:     AlertId,
			Name:        "Use JWT token",
			Description: blockRequest.RequestId,
		},
	)
	return response, nil
}
