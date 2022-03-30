package e2e_test

import "time"

func (s *Suite) TestLinkUnlink() {
	s.startForta(true)
	defer s.stopForta()

	tx, err := s.dispatchContract.Link(s.admin, agentIDBigInt, scannerIDBigInt)
	s.r.NoError(err)
	s.ensureTx("Dispatch.link() agent", tx)

	s.expectUpIn(time.Minute, agentContainerID)

	tx, err = s.dispatchContract.Unlink(s.admin, agentIDBigInt, scannerIDBigInt)
	s.r.NoError(err)
	s.ensureTx("Dispatch.unlink() agent", tx)

	s.expectDownIn(time.Minute, agentContainerID)
}
