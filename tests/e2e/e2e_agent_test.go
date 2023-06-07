package e2e_test

func (s *Suite) TestLinkUnlink() {
	s.startForta(true)
	s.expectIn(smallTimeout, func() bool {
		return s.fortaProcess.HasOutput("container started")
	})

	tx, err := s.mockRegistryContract.LinkTestAgent(s.deployer)
	s.r.NoError(err)
	s.ensureTx("link agent", tx)

	s.expectUpIn(largeTimeout, agentContainerID)

	s.expectIn(
		smallTimeout, func() (ok bool) {
			b := s.alertServer.GetLogs()
			return len(b) > 0
		},
	)

	tx, err = s.mockRegistryContract.UnlinkTestAgent(s.deployer)
	s.r.NoError(err)
	s.ensureTx("unlink agent", tx)

	s.stopForta()
}
