package mathexp

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
)

var seriesWithNil = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", nil, tp{
				time.Unix(5, 0), float64Pointer(2),
			}, tp{
				time.Unix(10, 0), nil,
			}),
		},
	},
}

var seriesEmpty = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", nil),
		},
	},
}

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
			name:        "sum series with a nil value",
			red:         "sum",
			varToReduce: "A",
			vars:        seriesWithNil,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("sum_", nil, NaN),
				},
			},
		},
		{
			name:        "sum empty series",
			red:         "sum",
			varToReduce: "A",
			vars:        seriesEmpty,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("sum_", nil, float64Pointer(0)),
				},
			},
		},
		{
			name:        "mean series with a nil value",
			red:         "mean",
			varToReduce: "A",
			vars:        seriesWithNil,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("mean_", nil, NaN),
				},
			},
		},
		{
			name:        "mean empty series",
			red:         "mean",
			varToReduce: "A",
			vars:        seriesEmpty,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("mean_", nil, NaN),
				},
			},
		},
		{
			name:        "min series with a nil value",
			red:         "min",
			varToReduce: "A",
			vars:        seriesWithNil,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("min_", nil, NaN),
				},
			},
		},
		{
			name:        "min empty series",
			red:         "min",
			varToReduce: "A",
			vars:        seriesEmpty,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("min_", nil, NaN),
				},
			},
		},
		{
			name:        "max series with a nil value",
			red:         "max",
			varToReduce: "A",
			vars:        seriesWithNil,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("max_", nil, NaN),
				},
			},
		},
		{
			name:        "max empty series",
			red:         "max",
			varToReduce: "A",
			vars:        seriesEmpty,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("max_", nil, NaN),
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
			opt := cmp.Comparer(func(x, y float64) bool {
				return (math.IsNaN(x) && math.IsNaN(y)) || x == y
			})
			options := append([]cmp.Option{opt}, data.FrameTestCompareOptions()...)
			if diff := cmp.Diff(tt.results, results, options...); diff != "" {
				t.Errorf("Result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
