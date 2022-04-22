package e2e_test

import (
	"encoding/json"
	"time"

	"github.com/forta-network/forta-core-go/clients/webhook/client/models"
	"github.com/forta-network/forta-node/cmd"
)

func (s *Suite) TestPrivateMode() {
	const privateModeDir = ".forta-private"

	// make sure that non-registered private nodes also cannot start
	s.forta(privateModeDir, "run")
	s.fortaProcess.Wait()
	s.True(s.fortaProcess.HasOutput(cmd.ErrCannotRunScanner.Error()))
	s.T().Log("as expected: could not run scan node without registration")

	s.registerNode()
	s.forta(privateModeDir, "run")
	defer s.stopForta()
	// the bot in private mode list should run
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
