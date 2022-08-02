package e2e_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/forta-network/forta-core-go/clients/webhook/client/models"
)

const localModeConfig = `chainId: 137

registry:
  checkIntervalSeconds: 1
  jsonRpc:
    url: http://localhost:8545

publish:
  batch:
    intervalSeconds: 1
    metricsBucketIntervalSeconds: 1

scan:
  jsonRpc:
    url: http://localhost:8545

localMode:
  enable: true
  includeMetrics: true
  webhookUrl: %s
  logFileName: %s
  botImages:
    - forta-e2e-test-agent
  runtimeLimits:
    startBlock: %d
    stopBlock: %d
    stopTimeoutSeconds: 30

autoUpdate:
  disable: true

trace:
  enabled: false

ens:
  override: true

telemetry:
  disable: true

agentLogs:
  disable: true

log:
  level: trace
`

const localModeDir = ".forta-local"

func (s *Suite) TestLocalModeWithWebhookClient() {
	webhookURL := "http://localhost:9090/batch/webhook"
	s.runLocalMode(webhookURL, "", func() ([]byte, bool) {
		return s.alertServer.GetAlert("webhook")
	})
}

func (s *Suite) TestLocalModeWithWebhookLogger() {
	webhookURL := "" // should cause the logger to be used
	logFileName := "test-log-file"
	logFilePath := path.Join(localModeDir, "logs", logFileName)
	os.Remove(logFilePath)
	s.runLocalMode(webhookURL, logFileName, func() ([]byte, bool) {
		b, err := ioutil.ReadFile(logFilePath)
		b = []byte(strings.TrimSpace(string(b)))
		return b, err == nil && len(b) > 0
	})
}

func (s *Suite) runLocalMode(webhookURL, logFileName string, readAlertsFunc func() ([]byte, bool)) {
	// get the start block number
	startBlockNumber, err := s.ethClient.BlockNumber(s.ctx)
	s.r.NoError(err)

	// send a transaction and let a block be mined
	s.sendExploiterTx()

	// get the stop block number
	lastBlockNumber, err := s.ethClient.BlockNumber(s.ctx)
	s.r.NoError(err)
	stopBlockNumber := lastBlockNumber + 1

	// change the config accordingly so we scan the block that includes the tx
	configFilePath := path.Join(localModeDir, "config.yml")
	os.Remove(configFilePath)
	s.r.NoError(
		ioutil.WriteFile(
			configFilePath,
			[]byte(fmt.Sprintf(localModeConfig, webhookURL, logFileName, startBlockNumber, stopBlockNumber)),
			0777,
		),
	)

	// make sure that non-registered local nodes also can start
	s.forta(localModeDir, "run")
	defer s.stopForta()
	// the bot in local mode list should run
	s.expectUpIn(smallTimeout, "forta-agent")

	var b []byte
	s.expectIn(smallTimeout, func() (ok bool) {
		b, ok = readAlertsFunc()
		return
	})
	var webhookAlerts models.AlertBatch
	s.r.NoError(json.Unmarshal(b, &webhookAlerts))
	s.r.NotEmpty(webhookAlerts.Alerts)
	s.r.NotEmpty(webhookAlerts.Metrics)

	// TODO: Find token from metadata of one of the alerts
	// _, err := security.VerifyScannerJWT(request.Header.Get("Authorization"))
	// s.r.NoError(err)

	s.T().Log(string(b))
}
