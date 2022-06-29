package e2e_test

// func (s *Suite) TestLinkUnlink() {
// 	s.startForta(true)
// 	s.expectIn(smallTimeout, func() bool {
// 		return s.fortaProcess.HasOutput("container started")
// 	})

// 	tx, err := s.dispatchContract.Link(s.admin, agentIDBigInt, scannerIDBigInt)
// 	s.r.NoError(err)
// 	s.ensureTx("Dispatch.link() agent", tx)

// 	s.expectUpIn(largeTimeout, agentContainerID)

// 	tx, err = s.dispatchContract.Unlink(s.admin, agentIDBigInt, scannerIDBigInt)
// 	s.r.NoError(err)
// 	s.ensureTx("Dispatch.unlink() agent", tx)

// 	s.stopForta()
// }
