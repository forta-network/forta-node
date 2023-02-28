package scanner

import (
	"testing"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestCombinerAlertAnalyzerService_createBloomFilter(t *testing.T) {
	type args struct {
		finding *protocol.Finding
	}
	tests := []struct {
		name            string
		args            args
		wantBloomFilter *protocol.BloomFilter
		wantErr         bool
	}{
		{
			name: "combination finding",
			args: args{
				finding: &protocol.Finding{Addresses: []string{"0xaaa"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				aas := &CombinerAlertAnalyzerService{}
				gotBloomFilter, err := aas.createBloomFilter(tt.args.finding)
				assert.Equal(t, tt.wantErr, err != nil)

				bf, err := utils.RecreateBloomFilter(gotBloomFilter)
				assert.NoError(t, err)

				// check for finding addresses
				for _, findingAddr := range tt.args.finding.Addresses {
					assert.True(t, bf.Test([]byte(findingAddr)), findingAddr)
				}
			},
		)
	}
}
