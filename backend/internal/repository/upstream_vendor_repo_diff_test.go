package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffAccountGroupIDs(t *testing.T) {
	tests := []struct {
		name     string
		existing []int64
		desired  []int64
		wantAdd  []int64
		wantDel  []int64
	}{
		{
			name:     "keeps unchanged rows untouched",
			existing: []int64{2, 3},
			desired:  []int64{2, 3},
			wantAdd:  nil,
			wantDel:  nil,
		},
		{
			name:     "adds and removes without touching the overlap",
			existing: []int64{1, 2, 3},
			desired:  []int64{2, 3, 4},
			wantAdd:  []int64{4},
			wantDel:  []int64{1},
		},
		{
			name:     "clears all groups",
			existing: []int64{1, 2},
			desired:  nil,
			wantAdd:  nil,
			wantDel:  []int64{1, 2},
		},
		{
			name:     "attaches to an account with no groups",
			existing: nil,
			desired:  []int64{5, 6},
			wantAdd:  []int64{5, 6},
			wantDel:  nil,
		},
		{
			name:     "ignores duplicates and non-positive ids",
			existing: []int64{1},
			desired:  []int64{0, -1, 1, 7, 7},
			wantAdd:  []int64{7},
			wantDel:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			add, remove := diffAccountGroupIDs(tt.existing, tt.desired)
			require.Equal(t, tt.wantAdd, add)
			require.Equal(t, tt.wantDel, remove)
		})
	}
}
