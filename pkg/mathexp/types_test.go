package mathexp

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

func TestSeriesSort(t *testing.T) {
	var tests = []struct {
		name           string
		descending     bool
		series         Series
		sortedSeriesIs assert.ComparisonAssertionFunc
		sortedSeries   Series
		panics         assert.PanicTestFunc
	}{
		{
			name:       "unordered series should sort by time ascending",
			descending: false,
			series: makeSeries("", nil, true, tp{
				unixTimePointer(3, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(1, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeries("", nil, true, tp{
				unixTimePointer(1, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(3, 0), float64Pointer(3),
			}),
		},
		{
			name:       "unordered series should sort by time descending",
			descending: true,
			series: makeSeries("", nil, true, tp{
				unixTimePointer(3, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(1, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeries("", nil, true, tp{
				unixTimePointer(3, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(1, 0), float64Pointer(1),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.series.SortByTime(tt.descending)
			tt.sortedSeriesIs(t, tt.series, tt.sortedSeries)
		})
	}
}
