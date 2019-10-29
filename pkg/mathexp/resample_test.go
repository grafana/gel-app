package mathexp

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/datasource"
	"github.com/stretchr/testify/assert"
)

func TestResampleSeries(t *testing.T) {
	var tests = []struct {
		name             string
		interval         string
		downsampler      string
		upsampler        string
		timeRange        datasource.TimeRange
		seriesToResample Series
		errIs            assert.ErrorAssertionFunc
		seriesIs         assert.ComparisonAssertionFunc
		series           Series
	}{
		{
			name:        "resample series: time range shorter than the rule interval",
			interval:    "5S",
			downsampler: "mean",
			upsampler:   "fillna",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(11, 0),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}),
			errIs: assert.Error,
		},
		{
			name:        "resample series: invalid time range",
			interval:    "5S",
			downsampler: "mean",
			upsampler:   "fillna",
			timeRange: datasource.TimeRange{
				From: time.Unix(11, 0),
				To:   time.Unix(0, 0),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}),
			errIs: assert.Error,
		},
		{
			name:        "resample series: downsampling (mean / pad)",
			interval:    "5S",
			downsampler: "mean",
			upsampler:   "pad",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(16, 0),
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
			name:        "resample series: downsampling (max / fillna)",
			interval:    "5S",
			downsampler: "max",
			upsampler:   "fillna",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(16, 0),
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
				unixTimePointer(5, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(15, 0), nil,
			}),
		},
		{
			name:        "resample series: downsampling (min / fillna)",
			interval:    "5S",
			downsampler: "min",
			upsampler:   "fillna",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(16, 0),
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
				unixTimePointer(5, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(15, 0), nil,
			}),
		},
		{
			name:        "resample series: downsampling (sum / fillna)",
			interval:    "5S",
			downsampler: "sum",
			upsampler:   "fillna",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(16, 0),
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
				unixTimePointer(5, 0), float64Pointer(5),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(3),
			}, tp{
				unixTimePointer(15, 0), nil,
			}),
		},
		{
			name:        "resample series: downsampling (mean / fillna)",
			interval:    "5S",
			downsampler: "mean",
			upsampler:   "fillna",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(16, 0),
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
				unixTimePointer(15, 0), nil,
			}),
		},
		{
			name:        "resample series: upsampling (mean / pad )",
			interval:    "2S",
			downsampler: "mean",
			upsampler:   "pad",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(11, 0),
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
		{
			name:        "resample series: upsampling (mean / backfilling )",
			interval:    "2S",
			downsampler: "mean",
			upsampler:   "backfilling",
			timeRange: datasource.TimeRange{
				From: time.Unix(0, 0),
				To:   time.Unix(11, 0),
			},
			seriesToResample: makeSeries("", nil, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(7, 0), float64Pointer(1),
			}),
			errIs:    assert.NoError,
			seriesIs: assert.Equal,
			series: makeSeries("", nil, tp{
				unixTimePointer(0, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(2, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(4, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(6, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(8, 0), float64Pointer(1),
			}, tp{
				unixTimePointer(10, 0), nil,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			series, err := tt.seriesToResample.Resample(tt.interval, tt.downsampler, tt.upsampler, tt.timeRange)
			tt.errIs(t, err)
			if err == nil {
				tt.seriesIs(t, tt.series, series)
			}
		})
	}
}
