package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBuildReleaseSummary(t *testing.T) {
	r := require.New(t)

	summary, ok := GetBuildReleaseSummary()
	r.False(ok)
	r.Nil(summary)
}

func TestGetBuildReleaseSummary_NonEmpty(t *testing.T) {
	r := require.New(t)

	CommitHash = "some hash"
	Version = "some version"
	ReleaseCid = "some release cid"
	summary, ok := GetBuildReleaseSummary()
	r.True(ok)
	r.NotNil(summary)
	r.Equal(summary.Commit, CommitHash)
	r.Equal(summary.Version, Version)
	r.Equal(summary.IPFS, ReleaseCid)
}
