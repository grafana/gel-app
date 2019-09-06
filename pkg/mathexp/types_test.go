package mathexp

import (
	"testing"
	"time"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/stretchr/testify/assert"
)

func unixTimePointer(sec, nsec int64) *time.Time {
	t := time.Unix(sec, nsec)
	return &t
}

func float64Pointer(f float64) *float64 {
	return &f
}

func TestFromTSDB(t *testing.T) {
	unixTimePointer := func(sec, nsec int64) *time.Time {
		t := time.Unix(sec, nsec)
		return &t
	}
	float64Pointer := func(f float64) *float64 {
		return &f
	}
	var tests = []struct {
		name            string
		tsdbSeriesSlice tsdb.TimeSeriesSlice
		seriesSetIs     assert.ComparisonAssertionFunc
		Results         Results
	}{
		{
			name: "it work maybe?",
			tsdbSeriesSlice: tsdb.TimeSeriesSlice{
				&tsdb.TimeSeries{
					Name: "temp",
					Points: []tsdb.TimePoint{
						{
							null.NewFloat(2.0, true),      // Value
							null.NewFloat(5.0*1000, true), // Time
						},
						{
							null.NewFloat(2.0, true),               // Value
							null.NewFloat(1566560853132.000, true), // Time
						},
					},
				},
			},
			seriesSetIs: assert.Equal,
			Results: Results{
				Values{
					Series{
						&data.Frame{Fields: data.Fields{
							&data.Field{
								Name:   "Time",
								Type:   data.TypeTime,
								Vector: &data.TimeVector{unixTimePointer(5, 0), unixTimePointer(1566560853, 132*int64(time.Millisecond))},
							},
							&data.Field{
								Name:   "temp",
								Type:   data.TypeNumber,
								Vector: &data.Float64Vector{float64Pointer(2.0), float64Pointer(2.0)},
							},
						},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := FromTSDB(tt.tsdbSeriesSlice)
			tt.seriesSetIs(t, tt.Results, ss)
		})
	}
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
			series: makeSeries("", nil, tp{
				unixTimePointer(3, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(1, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeries("", nil, tp{
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
			series: makeSeries("", nil, tp{
				unixTimePointer(3, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(1, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeries("", nil, tp{
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
