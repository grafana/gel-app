package mathexp

import (
	"testing"

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
			vars:        aSeries,
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
			vars:        aSeries,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("mean_", nil, float64Pointer(1.5)),
				},
			},
		},
	}
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
