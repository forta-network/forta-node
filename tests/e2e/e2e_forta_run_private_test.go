package e2e_test

import (
	"encoding/json"
	"time"

	"github.com/forta-network/forta-core-go/clients/webhook/client/models"
)

func (s *Suite) TestPrivateMode() {
	const localModeDir = ".forta-local"

	// make sure that non-registered local nodes also can start
	s.forta(localModeDir, "run")
	defer s.stopForta()
	// the bot in local mode list should run
	s.expectUpIn(smallTimeout, "forta-agent")

	// an alert should be detected and sent to the webhook url
	var b []byte
	s.expectIn(smallTimeout, func() (ok bool) {
		// trigger an alert
		s.sendExploiterTx()
		time.Sleep(time.Second * 2)
		// try to receive the alert
		b, ok = s.alertServer.GetAlert("webhook")
		return ok
	})
	var webhookAlerts models.AlertList
	s.r.NoError(json.Unmarshal(b, &webhookAlerts))
	s.T().Log(string(b))
}
