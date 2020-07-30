package mathexp

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
)

func TestSeriesReduce(t *testing.T) {
	var tests = []struct {
		name        string
		red         string
		vars        Vars
		varToReduce string
		errIs       assert.ErrorAssertionFunc
		resultsIs   assert.ComparisonAssertionFunc
		results     Results
	}{
		{
			name:        "sum series",
			red:         "sum",
			varToReduce: "A",
			vars:        aSeriesNullableTime,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("sum_", nil, float64Pointer(3)),
				},
			},
		},
		{
			name:        "mean series",
			red:         "mean",
			varToReduce: "A",
			vars:        aSeriesNullableTime,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("mean_", nil, float64Pointer(1.5)),
				},
			},
		},
		{
			name:        "mean series with labels",
			red:         "mean",
			varToReduce: "A",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeriesNullableTime("temp", data.Labels{"host": "a"}, nullTimeTP{
							unixTimePointer(5, 0), float64Pointer(2),
						}, nullTimeTP{
							unixTimePointer(10, 0), float64Pointer(1),
						}),
					},
				},
			},
			errIs:     assert.NoError,
			resultsIs: assert.Equal,
			results: Results{
				[]Value{
					makeNumber("mean_", data.Labels{"host": "a"}, float64Pointer(1.5)),
				},
			},
		},
	}
	//Vars{

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := Results{}
			seriesSet := tt.vars[tt.varToReduce]
			for _, series := range seriesSet.Values {
				ns, err := series.Value().(*Series).Reduce(tt.red)
				tt.errIs(t, err)
				if err != nil {
					t.Fail()
				}
				results.Values = append(results.Values, ns)
			}
			tt.resultsIs(t, tt.results, results)
		})
	}
}
