package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func unixTimePointer(sec, nsec int64) *time.Time {
	t := time.Unix(sec, nsec)
	return &t
}

func float64Pointer(f float64) *float64 {
	return &f
}

func TestTimeVector_Copy(t *testing.T) {
	tests := []struct {
		name     string
		tv       *TimeVector
		tvCopyIs assert.ComparisonAssertionFunc
		tvCopy   *TimeVector
	}{
		{
			name:     "copy should not be mutated",
			tv:       &TimeVector{unixTimePointer(0, 0), unixTimePointer(1, 0)},
			tvCopyIs: assert.Equal,
			tvCopy:   &TimeVector{unixTimePointer(0, 0), unixTimePointer(1, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.tv.Copy()
			tt.tv.SetValue(0, unixTimePointer(100, 100))
			tt.tvCopyIs(t, tt.tvCopy, v)
		})
	}
}
