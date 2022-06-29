package e2e_test

import (
	"github.com/forta-network/forta-node/cmd"
	"github.com/forta-network/forta-node/tests/e2e/ethaccounts"
)

// TestRegister tests what happens when registering with or without registration.
func (s *Suite) TestRegister() {
	s.forta("", "run")
	s.fortaProcess.Wait()
	s.True(s.fortaProcess.HasOutput(cmd.ErrCannotRunScanner.Error()))
	s.T().Log("as expected: could not run scan node without registration")

	s.T().Log("trying to run with --no-check")
	s.forta("", "run", "--no-check")
	s.expectUpIn(largeTimeout, runnerSupervisedContainers...)
	s.T().Log("--no-check works")
	s.stopForta()

	s.forta("", "register", "--owner-address", ethaccounts.ScannerOwnerAddress.Hex())
	s.fortaProcess.Wait()
	s.fortaProcess.HasOutput("polygonscan")

	// should work without pre-registration (false) now
	s.startForta(false)
	s.expectIn(smallTimeout, func() bool {
		return s.fortaProcess.HasOutput("container started")
	})
	s.stopForta()
}
