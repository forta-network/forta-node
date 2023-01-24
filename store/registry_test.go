package store

import (
	"testing"
)

func Test_calculateShardID(t *testing.T) {
	type args struct {
		target uint
		shards uint
		idx    uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		// [0,0,0,1,1,1]
		{
			name: "50% sharded",
			args: args{target: 3, shards: 2, idx: 3},
			want: 1,
		},
		// [0,0,1,1,2,2]
		{
			name: "66% sharded",
			args: args{target: 2, shards: 3, idx: 3},
			want: 1,
		},
		// [0,1,2,3,4,5]
		{
			name: "max redundancy",
			args: args{target: 1, shards: 6, idx: 5},
			want: 5,
		},
		// 	[0,0,0,0,0,0]
		{
			name: "no redundancy",
			args: args{target: 6, shards: 1, idx: 4},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := calculateShardID(tt.args.target, tt.args.idx); got != tt.want {
					t.Errorf("calculateShardID() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
