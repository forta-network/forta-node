package scanner

import (
	"context"
	"testing"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestTxAnalyzerService_createBloomFilter(t1 *testing.T) {
	type fields struct {
		ctx                context.Context
		cfg                TxAnalyzerServiceConfig
		lastInputActivity  health.TimeTracker
		lastOutputActivity health.TimeTracker
	}
	type args struct {
		finding *protocol.Finding
		event   *protocol.TransactionEvent
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantBloomFilter *protocol.BloomFilter
		wantErr         bool
	}{
		{
			name: "tx finding",
			args: args{
				finding: &protocol.Finding{Addresses: []string{"0xaaa"}},
				event: &protocol.TransactionEvent{
					Addresses: map[string]bool{
						"0xaaa": true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(
			tt.name, func(t1 *testing.T) {
				t := &TxAnalyzerService{}
				gotBloomFilter, err := t.createBloomFilter(tt.args.finding, tt.args.event)
				assert.Equal(t1, tt.wantErr, err != nil)

				bf, err := utils.RecreateBloomFilter(gotBloomFilter)
				assert.NoError(t1, err)

				// check for finding addresses
				for _, findingAddr := range tt.args.finding.Addresses {
					if !bf.Test([]byte(findingAddr)) {
						t1.Errorf("finding address %s does not exists in bloom filter", findingAddr)
					}
				}

				// check for tx addresses
				for txAddr := range tt.args.event.Addresses {
					if !bf.Test([]byte(txAddr)) {
						t1.Errorf("tx address %s does not exists in bloom filter", txAddr)
					}
				}
			},
		)
	}
}
