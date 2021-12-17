package utils_test

import (
	"testing"

	"github.com/forta-protocol/forta-node/utils"
	"github.com/stretchr/testify/require"
)

const testDiscoHost = "disco.forta.network"

func TestValidateDiscoImageRef(t *testing.T) {
	testCases := []struct {
		name        string
		ref         string
		fixedRef    string
		expectedErr error
	}{
		{
			name:        "valid default",
			ref:         "bafybeidslpugzaxfpvbhw3mknsohhdljgpqsimom6re7pbqwvtyzqtyi5m@sha256:6910b7c4806203b40bd8cd6e0d5b184280051b872517d8a37b4849a62ff9a014",
			fixedRef:    "disco.forta.network/bafybeidslpugzaxfpvbhw3mknsohhdljgpqsimom6re7pbqwvtyzqtyi5m",
			expectedErr: nil,
		},
		{
			name:        "long default",
			ref:         "disco.forta.network/bafybeidslpugzaxfpvbhw3mknsohhdljgpqsimom6re7pbqwvtyzqtyi5m@sha256:6910b7c4806203b40bd8cd6e0d5b184280051b872517d8a37b4849a62ff9a014",
			fixedRef:    "disco.forta.network/bafybeidslpugzaxfpvbhw3mknsohhdljgpqsimom6re7pbqwvtyzqtyi5m",
			expectedErr: nil,
		},
		{
			name:        "valid short not accepted",
			ref:         "bafybeidslpugzaxfpvbhw3mknsohhdljgpqsimom6re7pbqwvtyzqtyi5m",
			expectedErr: utils.ErrDiscoRefInvalid,
		},
		{
			name:        "bad cid",
			ref:         "bafybeidslpugzaxfpvbhw3mknsohhdljgpqsimom6re7pbqwvtyzqtyi5@sha256:6910b7c4806203b40bd8cd6e0d5b184280051b872517d8a37b4849a62ff9a014",
			expectedErr: utils.ErrDiscoRefNotIPFSCIDv1,
		},
		{
			name:        "totally invalid",
			ref:         "asdafdsgsdsf",
			expectedErr: utils.ErrDiscoRefInvalid,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := require.New(t)

			fixedRef, err := utils.ValidateDiscoImageRef(testDiscoHost, testCase.ref)
			r.ErrorIs(err, testCase.expectedErr)
			r.Equal(testCase.fixedRef, fixedRef)
		})
	}
}
