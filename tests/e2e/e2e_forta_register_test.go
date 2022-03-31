package e2e_test

import (
	"time"

	"github.com/forta-protocol/forta-node/cmd"
	"github.com/forta-protocol/forta-node/tests/e2e/ethaccounts"
)

// TestRegister_NoRegisterRun tests the cases when the node is ran without registering
// and the check is ignored with the --no-check flag.
func (s *Suite) TestRegister_NoRegisterRun() {
	fortaMain := s.forta("run")
	s.r.ErrorIs(fortaMain.ErrorAfter(time.Second*5), cmd.ErrCannotRunScanner)
	s.T().Log("successfully got the error")

	fortaMain = s.forta("run", "--no-check")
	defer s.stopForta()
	s.T().Log("ran with --no-check")
	s.r.NoError(fortaMain.ErrorAfter(time.Second * 5))
	s.expectUpIn(time.Minute, serviceContainers...)
}

// TestRegister_RegisterRun tests a run after normal registering.
func (s *Suite) TestRegister_RegisterRun() {
	fortaMain := s.forta("register", "--owner-address", ethaccounts.ScannerOwnerAddress.Hex())
	fortaMain.Wait()
	s.r.NoError(fortaMain.ErrorAfter(time.Second * 5))

	s.startForta(false)
	s.stopForta()
}
