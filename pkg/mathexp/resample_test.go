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
			name:     "resample series",
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
			errIs:    assert.NoError,
			seriesIs: assert.Equal,
			series: makeSeries("", nil, tp{
				unixTimePointer(0, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(5, 0), float64Pointer(1),
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
