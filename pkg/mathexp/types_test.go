package mathexp

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/stretchr/testify/assert"
)

// Common Test Constructor Utils
type nullTimeTP struct {
	t *time.Time
	f *float64
}

type tp struct {
	t time.Time
	f *float64
}

func makeSeriesNullableTime(name string, labels dataframe.Labels, points ...nullTimeTP) Series {
	newSeries := NewSeries(name, labels, true, len(points))
	for idx, p := range points {
		newSeries.SetPoint(idx, p.t, p.f)
	}
	return newSeries
}

func makeSeries(name string, labels dataframe.Labels, points ...tp) Series {
	newSeries := NewSeries(name, labels, false, len(points))
	for idx, p := range points {
		err := newSeries.SetPoint(idx, &p.t, p.f)
		if err != nil {
			panic(err)
		}
	}
	return newSeries
}

func makeNumber(name string, labels dataframe.Labels, f *float64) Number {
	newNumber := NewNumber(name, labels)
	newNumber.SetValue(f)
	return newNumber
}

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
