package e2e_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/forta-network/forta-core-go/clients/webhook/client/models"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-node/tests/e2e/agents/combinerbot/combinerbotalertid"
	"github.com/forta-network/forta-node/tests/e2e/agents/txdetectoragent/testbotalertid"
)

const localModeConfig = `chainId: 137

registry:
  checkIntervalSeconds: 1
  jsonRpc:
    url: http://localhost:8545
  containerRegistry: localhost:1970

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
const localModeAlertConfig = `chainId: 137

registry:
  checkIntervalSeconds: 1
  jsonRpc:
    url: http://localhost:8545
  containerRegistry: localhost:1970


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
    - forta-e2e-alert-test-agent
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

combiner:
  alertApiUrl: https://api-dev.forta.network/graphql

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
	logFileName := logFileName()
	logFilePath := path.Join(localModeDir, "logs", logFileName)
	_ = os.RemoveAll(logFilePath)
	s.runLocalMode(webhookURL, logFileName, func() ([]byte, bool) {
		b, err := ioutil.ReadFile(logFilePath)
		b = []byte(strings.TrimSpace(string(b)))
		return b, err == nil && len(b) > 0
	})
}

func (s *Suite) TestLocalModeAlertHandlingWithWebhookLogger() {
	webhookURL := "" // should cause the logger to be used
	logFileName := logFileName()
	logFilePath := path.Join(localModeDir, "logs", logFileName)
	_ = os.RemoveAll(logFilePath)
	s.runLocalModeAlertHandler(
		webhookURL, logFileName, func() ([]byte, bool) {
			b, err := ioutil.ReadFile(logFilePath)
			b = []byte(strings.TrimSpace(string(b)))
			return b, err == nil && len(b) > 0
		},
	)
}

func logFileName() string {
	return fmt.Sprintf("test-log-file-%d", time.Now().UnixNano())
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
	_ = os.RemoveAll(configFilePath)
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
	s.r.Len(webhookAlerts.Alerts, 2)
	s.r.NotEmpty(webhookAlerts.Metrics)

	var (
		exploiterAlert *models.Alert
		tokenAlert     *models.Alert
	)

	for _, alert := range webhookAlerts.Alerts {
		if alert.AlertID == testbotalertid.ExploiterAlertId {
			exploiterAlert = alert
		}
		if alert.AlertID == testbotalertid.TokenAlertId {
			tokenAlert = alert
		}
	}
	s.r.NotNil(exploiterAlert)
	s.r.NotNil(tokenAlert)

	_, err = security.VerifyScannerJWT(s.getTokenFromAlert(tokenAlert))
	s.r.NoError(err)

	s.T().Log(string(b))
}

func (s *Suite) runLocalModeAlertHandler(webhookURL, logFileName string, readAlertsFunc func() ([]byte, bool)) {
	// change the config accordingly so we scan the block that includes the tx
	configFilePath := path.Join(localModeDir, "config.yml")
	_ = os.RemoveAll(configFilePath)
	s.r.NoError(
		ioutil.WriteFile(
			configFilePath,
			[]byte(fmt.Sprintf(
				localModeAlertConfig, webhookURL, logFileName, 0, 0,
			)),
			0777,
		),
	)

	// make sure that non-registered local nodes also can start
	s.forta(localModeDir, "run")
	defer s.stopForta()
	// the bot in local mode list should run
	s.expectUpIn(smallTimeout, "forta-agent")

	var b []byte
	s.expectIn(
		smallTimeout, func() (ok bool) {
			b, ok = readAlertsFunc()
			return
		},
	)
	var webhookAlerts models.AlertBatch
	s.r.NoError(json.Unmarshal(b, &webhookAlerts))
	s.r.GreaterOrEqual(len(webhookAlerts.Alerts), 1)
	s.r.NotEmpty(webhookAlerts.Metrics)

	var (
		combinationAlert *models.Alert
	)

	for _, alert := range webhookAlerts.Alerts {
		if alert.AlertID == combinerbotalertid.CombinationAlertID {
			combinationAlert = alert
		}
	}

	s.r.NotNil(combinationAlert)

	s.T().Log(string(b))
}

func (s *Suite) getTokenFromAlert(tokenAlert *models.Alert) string {
	tokenAlertMeta, ok := tokenAlert.Metadata.(map[string]interface{})
	s.r.True(ok)
	tokenV := tokenAlertMeta["token"]
	s.r.NotNil(tokenV)
	token, ok := tokenV.(string)
	s.r.True(ok)
	return token
}
