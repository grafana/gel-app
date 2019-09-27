package mathexp

import (
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/stretchr/testify/assert"
)

func TestResampleSeries(t *testing.T) {
	var tests = []struct {
		name             string
		interval         string
		timeRange        *datasource.TimeRange
		seriesToResample Series
		errIs            assert.ErrorAssertionFunc
		seriesIs         assert.ComparisonAssertionFunc
		series           Series
	}{
		{
			name:     "resample series: time range shorter than the rule interval",
			interval: "5S",
			timeRange: &datasource.TimeRange{
				FromRaw: fmt.Sprintf("%v", time.Unix(0, 0).Unix()+1e3),
				ToRaw:   fmt.Sprintf("%v", time.Unix(11, 0).Unix()+1e3),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}),
			errIs: assert.Error,
		},
		{
			name:     "resample series: invalid time range",
			interval: "5S",
			timeRange: &datasource.TimeRange{
				FromRaw: fmt.Sprintf("%v", time.Unix(11, 0).Unix()+1e3),
				ToRaw:   fmt.Sprintf("%v", time.Unix(0, 0).Unix()+1e3),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}),
			errIs: assert.Error,
		},
		{
			name:     "resample series: downsampling (mean aggregation)",
			interval: "5S",
			timeRange: &datasource.TimeRange{
				FromRaw: fmt.Sprintf("%v", time.Unix(0, 0).Unix()*1e3),
				ToRaw:   fmt.Sprintf("%v", time.Unix(16, 0).Unix()*1e3),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(4, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(9, 0), float64Pointer(2),
			}),
			errIs:    assert.NoError,
			seriesIs: assert.Equal,
			series: makeSeries("", nil, tp{
				unixTimePointer(0, 0), nil,
			}, tp{
				unixTimePointer(5, 0), float64Pointer(2.5),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(1.5),
			}, tp{
				unixTimePointer(15, 0), float64Pointer(2),
			}),
		},
		{
			name:     "resample series: upsampling (bfill)",
			interval: "2S",
			timeRange: &datasource.TimeRange{
				FromRaw: fmt.Sprintf("%v", time.Unix(0, 0).Unix()*1e3),
				ToRaw:   fmt.Sprintf("%v", time.Unix(11, 0).Unix()*1e3),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}),
			errIs:    assert.NoError,
			seriesIs: assert.Equal,
			series: makeSeries("", nil, tp{
				unixTimePointer(0, 0), nil,
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(4, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(6, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(8, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(1),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			series, err := tt.seriesToResample.Resample(tt.interval, tt.timeRange)
			tt.errIs(t, err)
			if err == nil {
				tt.seriesIs(t, tt.series, series)
			}
		})
	}
}
