package mathexp

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
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

func TestSeriesFromFrame(t *testing.T) {
	var tests = []struct {
		name   string
		frame  *dataframe.Frame
		errIs  assert.ErrorAssertionFunc
		Is     assert.ComparisonAssertionFunc
		Series Series
	}{
		{
			name: "[]time, []float frame should convert",
			frame: &dataframe.Frame{
				Name:   "test",
				Labels: nil,
				Fields: []*dataframe.Field{
					dataframe.NewField("time", []time.Time{}),
					dataframe.NewField("value", []float64{}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &dataframe.Frame{
					Name:   "test",
					Labels: nil,
					Fields: []*dataframe.Field{
						dataframe.NewField("time", []time.Time{}),
						dataframe.NewField("value", []float64{}),
					},
				},
				TimeIdx:        0,
				TimeIsNullable: false,
				ValueIdx:       1,
				ValueIsNullabe: false,
			},
		},
		{
			name: "[]*float, []*time frame should convert",
			frame: &dataframe.Frame{
				Name:   "test",
				Labels: nil,
				Fields: []*dataframe.Field{
					dataframe.NewField("value", []*float64{float64Pointer(5)}),
					dataframe.NewField("time", []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &dataframe.Frame{
					Name:   "test",
					Labels: nil,
					Fields: []*dataframe.Field{
						dataframe.NewField("value", []*float64{float64Pointer(5)}),
						dataframe.NewField("time", []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:        1,
				TimeIsNullable: true,
				ValueIdx:       0,
				ValueIsNullabe: true,
			},
		},
		{
			name: "[]*time, []*time frame should error",
			frame: &dataframe.Frame{
				Name:   "test",
				Labels: nil,
				Fields: []*dataframe.Field{
					dataframe.NewField("time", []*time.Time{}),
					dataframe.NewField("time", []*time.Time{}),
				},
			},
			errIs: assert.Error,
		},
		{
			name: "[]*float64, []float64 frame should error",
			frame: &dataframe.Frame{
				Name:   "test",
				Labels: nil,
				Fields: []*dataframe.Field{
					dataframe.NewField("value", []*float64{}),
					dataframe.NewField("value", []*float64{}),
				},
			},
			errIs: assert.Error,
		},
		{
			name: "[]*float64 frame should error",
			frame: &dataframe.Frame{
				Name:   "test",
				Labels: nil,
				Fields: []*dataframe.Field{
					dataframe.NewField("value", []*float64{}),
				},
			},
			errIs: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := SeriesFromFrame(tt.frame)
			tt.errIs(t, err)
			if err == nil {
				tt.Is(t, s, tt.Series)
			}
		})
	}
}
