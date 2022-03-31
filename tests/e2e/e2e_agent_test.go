package e2e_test

import "time"

func (s *Suite) TestLinkUnlink() {
	fortaMain := s.startForta(true)
	s.r.NoError(fortaMain.ErrorAfter(time.Second * 5))
	defer s.stopForta()

	tx, err := s.dispatchContract.Link(s.admin, agentIDBigInt, scannerIDBigInt)
	s.r.NoError(err)
	s.ensureTx("Dispatch.link() agent", tx)

	s.expectUpIn(largeTimeout, agentContainerID)
	s.r.NoError(fortaMain.ErrorAfter(time.Second * 5))

	tx, err = s.dispatchContract.Unlink(s.admin, agentIDBigInt, scannerIDBigInt)
	s.r.NoError(err)
	s.ensureTx("Dispatch.unlink() agent", tx)

	s.expectDownIn(largeTimeout, agentContainerID)
	s.r.NoError(fortaMain.ErrorAfter(time.Second * 5))
}
