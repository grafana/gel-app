package mathexp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			series: makeSeriesNullableTime("", nil, nullTimeTP{
				unixTimePointer(3, 0), float64Pointer(3),
			}, nullTimeTP{
				unixTimePointer(1, 0), float64Pointer(1),
			}, nullTimeTP{
				unixTimePointer(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeriesNullableTime("", nil, nullTimeTP{
				unixTimePointer(1, 0), float64Pointer(1),
			}, nullTimeTP{
				unixTimePointer(2, 0), float64Pointer(2),
			}, nullTimeTP{
				unixTimePointer(3, 0), float64Pointer(3),
			}),
		},
		{
			name:       "unordered series should sort by time descending",
			descending: true,
			series: makeSeriesNullableTime("", nil, nullTimeTP{
				unixTimePointer(3, 0), float64Pointer(3),
			}, nullTimeTP{
				unixTimePointer(1, 0), float64Pointer(1),
			}, nullTimeTP{
				unixTimePointer(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeriesNullableTime("", nil, nullTimeTP{
				unixTimePointer(3, 0), float64Pointer(3),
			}, nullTimeTP{
				unixTimePointer(2, 0), float64Pointer(2),
			}, nullTimeTP{
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
