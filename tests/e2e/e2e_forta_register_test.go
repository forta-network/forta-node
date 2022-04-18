package e2e_test

import (
	"github.com/forta-network/forta-node/cmd"
	"github.com/forta-network/forta-node/tests/e2e/ethaccounts"
)

// TestRegister_NoRegisterRun tests the cases when the node is ran without registering
// and the check is ignored with the --no-check flag.
func (s *Suite) TestRegister_NoRegisterRun() {
	s.forta("run")
	s.fortaProcess.Wait()
	s.True(s.fortaProcess.HasOutput(cmd.ErrCannotRunScanner.Error()))
	s.T().Log("as expected: could not scanner without registration")

	s.T().Log("trying to run with --no-check")
	s.forta("run", "--no-check")
	defer s.stopForta()
	s.expectUpIn(largeTimeout, runnerSupervisedContainers...)
	s.T().Log("--no-check works")
}

// TestRegister_RegisterRun tests a run after normal registering.
func (s *Suite) TestRegister_RegisterRun() {
	s.forta("register", "--owner-address", ethaccounts.ScannerOwnerAddress.Hex())
	s.fortaProcess.Wait()
	s.fortaProcess.HasOutput("polygonscan")

	// should work without pre-registration (false) now
	s.startForta(false)
	s.expectIn(smallTimeout, func() bool {
		return s.fortaProcess.HasOutput("container started")
	})
	s.stopForta()
}
