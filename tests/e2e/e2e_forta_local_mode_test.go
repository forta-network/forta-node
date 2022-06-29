package e2e_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
  webhookUrl: http://localhost:9090/batch/webhook
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

func (s *Suite) TestLocalMode() {
	const localModeDir = ".forta-local"

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
			configFilePath, []byte(fmt.Sprintf(localModeConfig, startBlockNumber, stopBlockNumber)), 0777,
		),
	)

	// make sure that non-registered local nodes also can start
	s.forta(localModeDir, "run")
	defer s.stopForta()
	// the bot in local mode list should run
	s.expectUpIn(smallTimeout, "forta-agent")

	// an alert should be detected and sent to the webhook url
	var b []byte
	s.expectIn(smallTimeout, func() (ok bool) {
		// try to receive the alert
		b, ok = s.alertServer.GetAlert("webhook")
		return ok
	})
	var webhookAlerts models.AlertBatch
	s.r.NoError(json.Unmarshal(b, &webhookAlerts))
	s.r.NotEmpty(webhookAlerts.Alerts)
	s.r.NotEmpty(webhookAlerts.Metrics)
	s.T().Log(string(b))
}
