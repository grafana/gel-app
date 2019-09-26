package mathexp

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/mathexp/parse"
	"github.com/stretchr/testify/assert"
)

func TestScalarExpr(t *testing.T) {
	var tests = []struct {
		name      string
		expr      string
		vars      Vars
		newErrIs  assert.ErrorAssertionFunc
		execErrIs assert.ErrorAssertionFunc
		resultIs  assert.ComparisonAssertionFunc
		Results   Results
	}{
		{
			name:      "a scalar",
			expr:      "1",
			vars:      Vars{},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			Results:   Results{[]Value{NewScalar(float64Pointer(1.0))}},
		},
		{
			name:      "unary: scalar",
			expr:      "! 1.2",
			vars:      Vars{},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			Results:   Results{[]Value{NewScalar(float64Pointer(0.0))}},
		},
		{
			name:      "binary: scalar Op scalar",
			expr:      "1 + 1",
			vars:      Vars{},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			Results:   Results{[]Value{NewScalar(float64Pointer(2.0))}},
		},
		{
			name:      "binary: scalar Op scalar - divide by zero",
			expr:      "1 / 0",
			vars:      Vars{},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			Results:   Results{[]Value{NewScalar(float64Pointer(math.Inf(1)))}},
		},
		{
			name:      "binary: scalar Op number",
			expr:      "1 + $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("temp", nil, float64Pointer(2.0))}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			Results:   Results{[]Value{makeNumber("", nil, float64Pointer(3.0))}},
		},
		{
			name:      "binary: number Op Scalar",
			expr:      "$A - 3",
			vars:      Vars{"A": Results{[]Value{makeNumber("temp", nil, float64Pointer(2.0))}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			Results:   Results{[]Value{makeNumber("", nil, float64Pointer(-1))}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.expr)
			tt.newErrIs(t, err)
			if e != nil {
				res, err := e.Execute(tt.vars)
				tt.execErrIs(t, err)
				tt.resultIs(t, tt.Results, res)
			}
		})
	}
}

func TestNumberExpr(t *testing.T) {
	var tests = []struct {
		name      string
		expr      string
		vars      Vars
		newErrIs  assert.ErrorAssertionFunc
		execErrIs assert.ErrorAssertionFunc
		resultIs  assert.ComparisonAssertionFunc
		results   Results
		willPanic bool
	}{
		{
			name:      "binary: number Op Scalar",
			expr:      "$A / $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("temp", nil, float64Pointer(2.0))}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results:   Results{[]Value{makeNumber("", nil, float64Pointer(1))}},
		},
		{
			name:      "unary: number",
			expr:      "- $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("temp", nil, float64Pointer(2.0))}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results:   Results{[]Value{makeNumber("", nil, float64Pointer(-2.0))}},
		},
		{
			name:      "binary: Scalar Op Number (Number will nil val) - currently Panics",
			expr:      "1 + $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("", nil, nil)}}},
			willPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBlock := func() {
				e, err := New(tt.expr)
				tt.newErrIs(t, err)
				if e != nil {
					res, err := e.Execute(tt.vars)
					tt.execErrIs(t, err)
					tt.resultIs(t, tt.results, res)
				}
			}
			if tt.willPanic {
				assert.Panics(t, testBlock)
			} else {
				assert.NotPanics(t, testBlock)
			}
		})
	}
}

func TestSeriesExpr(t *testing.T) {
	var tests = []struct {
		name      string
		expr      string
		vars      Vars
		newErrIs  assert.ErrorAssertionFunc
		execErrIs assert.ErrorAssertionFunc
		resultIs  assert.ComparisonAssertionFunc
		results   Results
	}{
		{
			name:      "unary series",
			expr:      "! ! $A",
			vars:      aSeries,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{ // Not sure about preservering names...
						unixTimePointer(5, 0), float64Pointer(1),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(1),
					}),
				},
			},
		},
		{
			name:      "binary scalar Op series",
			expr:      "98 + $A",
			vars:      aSeries,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{ // Not sure about preservering names...
						unixTimePointer(5, 0), float64Pointer(100),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(99),
					}),
				},
			},
		},
		{
			name:      "binary series Op scalar",
			expr:      "$A + 98",
			vars:      aSeries,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{ // Not sure about preservering names...
						unixTimePointer(5, 0), float64Pointer(100),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(99),
					}),
				},
			},
		},
		{
			name:      "series Op series",
			expr:      "$A + $A",
			vars:      aSeries,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{ // Not sure about preservering names...
						unixTimePointer(5, 0), float64Pointer(4),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(2),
					}),
				},
			},
		},
		{
			name:      "series Op number",
			expr:      "$A + $B",
			vars:      aSeriesbNumber,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("id=1", data.Labels{"id": "1"}, tp{
						unixTimePointer(5, 0), float64Pointer(9),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(8),
					}),
				},
			},
		},
		{
			name:      "number Op series",
			expr:      "$B + $A",
			vars:      aSeriesbNumber,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("id=1", data.Labels{"id": "1"}, tp{
						unixTimePointer(5, 0), float64Pointer(9),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(8),
					}),
				},
			},
		},
		{
			name:      "series Op series with label union",
			expr:      "$A * $B",
			vars:      twoSeriesSets,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("sensor=a, turbine=1", data.Labels{"sensor": "a", "turbine": "1"}, tp{
						unixTimePointer(5, 0), float64Pointer(6 * .5),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(8 * .2),
					}),
					makeSeries("sensor=b, turbine=1", data.Labels{"sensor": "b", "turbine": "1"}, tp{
						unixTimePointer(5, 0), float64Pointer(10 * .5),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(16 * .2),
					}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.expr)
			tt.newErrIs(t, err)
			if e != nil {
				res, err := e.Execute(tt.vars)
				tt.execErrIs(t, err)
				tt.resultIs(t, tt.results, res)
			}
		})
	}
}

var aSeries = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", nil, tp{
				unixTimePointer(5, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(1),
			}),
		},
	},
}

var aSeriesbNumber = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", nil, tp{
				unixTimePointer(5, 0), float64Pointer(2),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(1),
			}),
		},
	},
	"B": Results{
		[]Value{
			makeNumber("volt", data.Labels{"id": "1"}, float64Pointer(7)),
		},
	},
}

var twoSeriesSets = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", data.Labels{"sensor": "a", "turbine": "1"}, tp{
				unixTimePointer(5, 0), float64Pointer(6),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(8),
			}),
			makeSeries("temp", data.Labels{"sensor": "b", "turbine": "1"}, tp{
				unixTimePointer(5, 0), float64Pointer(10),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(16),
			}),
		},
	},
	"B": Results{
		[]Value{
			makeSeries("efficiency", data.Labels{"turbine": "1"}, tp{
				unixTimePointer(5, 0), float64Pointer(.5),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(.2),
			}),
		},
	},
}

type tp struct {
	t *time.Time
	f *float64
}

func makeSeries(name string, labels data.Labels, points ...tp) Series {
	newSeries := NewSeries(name, labels, len(points))
	for idx, p := range points {
		newSeries.SetPoint(idx, p.t, p.f)
	}
	return newSeries
}

func makeNumber(name string, labels data.Labels, f *float64) Number {
	newNumber := NewNumber(name, labels)
	newNumber.SetValue(f)
	return newNumber
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
			vars:        aSeries,
			errIs:       assert.NoError,
			resultsIs:   assert.Equal,
			results: Results{
				[]Value{
					makeNumber("sum_", nil, float64Pointer(3)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := Results{}
			seriesSet := tt.vars[tt.varToReduce]
			for _, series := range seriesSet.Values {
				if series.Type() == parse.TypeSeriesSet {
					ns, err := series.Value().(*Series).Reduce(tt.red)
					tt.errIs(t, err)
					if err != nil {
						t.Fail()
					}
					results.Values = append(results.Values, ns)
				}
			}
			tt.resultsIs(t, tt.results, results)
		})
	}
}

func Test_union(t *testing.T) {
	var tests = []struct {
		name      string
		aResults  Results
		bResults  Results
		unionsAre assert.ComparisonAssertionFunc
		unions    []*Union
	}{
		{
			name: "equal tags single union",
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: data.Labels{"id": "1"},
					A:      makeSeries("a", data.Labels{"id": "1"}),
					B:      makeSeries("b", data.Labels{"id": "1"}),
				},
			},
		},
		{
			name: "equal tags keys with no matching values will result in no unions",
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "2"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions:    []*Union{},
		},
		{
			name:      "empty results will result in no unions",
			aResults:  Results{},
			bResults:  Results{},
			unionsAre: assert.EqualValues,
			unions:    []*Union{},
		},
		{
			name: "incompatible tags of different length with will result in no unions",
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"ID": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions:    []*Union{},
		},
		{
			name: "A is subset of B results in single union with Labels of B",
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: data.Labels{"id": "1", "fish": "herring"}, // Union gets the labels that is not the subset
					A:      makeSeries("a", data.Labels{"id": "1"}),
					B:      makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
				},
			},
		},
		{
			name: "B is subset of A results in single union with Labels of A",
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1", "fish": "herring"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: data.Labels{"id": "1", "fish": "herring"}, // Union gets the labels that is not the subset
					A:      makeSeries("a", data.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("b", data.Labels{"id": "1"}),
				},
			},
		},
		{
			name: "single valued A is subset of many valued B, results in many union with Labels of B",
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
					makeSeries("b", data.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: data.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("a", data.Labels{"id": "1"}),
					B:      makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
				},
				{
					Labels: data.Labels{"id": "1", "fish": "red snapper"},
					A:      makeSeries("a", data.Labels{"id": "1"}),
					B:      makeSeries("b", data.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
		},
		{
			name: "A with different tags keys lengths to B makes 3 unions (with two unions have matching tags)",
			// Is this the behavior we want? A result within the results will no longer
			// be uniquely identifiable.
			aResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1"}),
					makeSeries("aa", data.Labels{"id": "1", "fish": "herring"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
					makeSeries("bb", data.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: data.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("a", data.Labels{"id": "1"}),
					B:      makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
				},
				{
					Labels: data.Labels{"id": "1", "fish": "red snapper"},
					A:      makeSeries("a", data.Labels{"id": "1"}),
					B:      makeSeries("bb", data.Labels{"id": "1", "fish": "red snapper"}),
				},
				{
					Labels: data.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("aa", data.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
				},
			},
		},
		{
			name: "B with different tags keys lengths to A makes 3 unions (with two unions have matching tags)",
			// Is this the behavior we want? A result within the results will no longer
			// be uniquely identifiable.
			aResults: Results{
				Values: Values{
					makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
					makeSeries("bb", data.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("a", data.Labels{"id": "1"}),
					makeSeries("aa", data.Labels{"id": "1", "fish": "herring"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: data.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("a", data.Labels{"id": "1"}),
				},
				{
					Labels: data.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("b", data.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("aa", data.Labels{"id": "1", "fish": "herring"}),
				},
				{
					Labels: data.Labels{"id": "1", "fish": "red snapper"},
					A:      makeSeries("bb", data.Labels{"id": "1", "fish": "red snapper"}),
					B:      makeSeries("a", data.Labels{"id": "1"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unions := union(tt.aResults, tt.bResults)
			tt.unionsAre(t, tt.unions, unions)
		})
	}
}

func TestFunc(t *testing.T) {
	var tests = []struct {
		name      string
		expr      string
		vars      Vars
		newErrIs  assert.ErrorAssertionFunc
		execErrIs assert.ErrorAssertionFunc
		resultIs  assert.ComparisonAssertionFunc
		results   Results
	}{
		{
			name: "abs on number",
			expr: "abs($A)",
			vars: Vars{
				"A": Results{
					[]Value{
						makeNumber("", nil, float64Pointer(-7)),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results:   Results{[]Value{makeNumber("", nil, float64Pointer(7))}},
		},
		{
			name:      "abs on scalar",
			expr:      "abs(-1)",
			vars:      Vars{},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results:   Results{[]Value{NewScalar(float64Pointer(1.0))}},
		},
		{
			name: "abs on series",
			expr: "abs($A)",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("", nil, tp{
							unixTimePointer(5, 0), float64Pointer(-2),
						}, tp{
							unixTimePointer(10, 0), float64Pointer(-1),
						}),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			resultIs:  assert.Equal,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), float64Pointer(2),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(1),
					}),
				},
			},
		},
		{
			name:     "abs on string - should error",
			expr:     `abs("hi")`,
			vars:     Vars{},
			newErrIs: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.expr)
			tt.newErrIs(t, err)
			if e != nil {
				res, err := e.Execute(tt.vars)
				tt.execErrIs(t, err)
				tt.resultIs(t, tt.results, res)
			}
		})
	}
}

// NaN is just to make the calls a little cleaner, the one
// call is not for any sort of equality side effect in tests.
// note: cmp.Equal must be used to test Equality for NaNs.
var NaN = float64Pointer(math.NaN())

func TestNaN(t *testing.T) {
	var tests = []struct {
		name      string
		expr      string
		vars      Vars
		newErrIs  assert.ErrorAssertionFunc
		execErrIs assert.ErrorAssertionFunc
		results   Results
		willPanic bool
	}{
		{
			name:      "unary !: Op Number(NaN) is NaN",
			expr:      "! $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("", nil, NaN)}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   Results{[]Value{makeNumber("", nil, NaN)}},
		},
		{
			name:      "unary -: Op Number(NaN) is NaN",
			expr:      "! $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("", nil, NaN)}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   Results{[]Value{makeNumber("", nil, NaN)}},
		},
		{
			name:      "binary: Scalar Op(Non-AND/OR) Number(NaN) is NaN",
			expr:      "1 * $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("", nil, NaN)}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   Results{[]Value{makeNumber("", nil, NaN)}},
		},
		{
			name:      "binary: Scalar Op(AND/OR) Number(NaN) is 0/1",
			expr:      "1 || $A",
			vars:      Vars{"A": Results{[]Value{makeNumber("", nil, NaN)}}},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   Results{[]Value{makeNumber("", nil, float64Pointer(1))}},
		},
		{
			name: "binary: Scalar Op(Non-AND/OR) Series(with NaN value) is NaN)",
			expr: "1 - $A",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("temp", nil, tp{
							unixTimePointer(5, 0), float64Pointer(2),
						}, tp{
							unixTimePointer(10, 0), NaN,
						}),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), float64Pointer(-1),
					}, tp{
						unixTimePointer(10, 0), NaN,
					}),
				},
			},
		},
		{
			name: "binary: Number Op(Non-AND/OR) Series(with NaN value) is Series with NaN",
			expr: "$A == $B",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("temp", nil, tp{
							unixTimePointer(5, 0), float64Pointer(2),
						}, tp{
							unixTimePointer(10, 0), NaN,
						}),
					},
				},
				"B": Results{[]Value{makeNumber("", nil, float64Pointer(0))}},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), float64Pointer(0),
					}, tp{
						unixTimePointer(10, 0), NaN,
					}),
				},
			},
		},
		{
			name: "binary: Number(NaN) Op Series(with NaN value) is Series with NaN",
			expr: "$A + $B",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("temp", nil, tp{
							unixTimePointer(5, 0), float64Pointer(2),
						}, tp{
							unixTimePointer(10, 0), NaN,
						}),
					},
				},
				"B": Results{[]Value{makeNumber("", nil, NaN)}},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), NaN,
					}, tp{
						unixTimePointer(10, 0), NaN,
					}),
				},
			},
		},
	}

	// go-cmp instead of testify assert is used to compare results here
	// because it supports an option for NaN equality.
	// https://github.com/stretchr/testify/pull/691#issuecomment-528457166
	opt := cmp.Comparer(func(x, y float64) bool {
		return (math.IsNaN(x) && math.IsNaN(y)) || x == y
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBlock := func() {
				e, err := New(tt.expr)
				tt.newErrIs(t, err)
				if e != nil {
					res, err := e.Execute(tt.vars)
					tt.execErrIs(t, err)
					if !cmp.Equal(res, tt.results, opt) {
						assert.FailNow(t, tt.name, cmp.Diff(res, tt.results, opt))
					}
				}
			}
			if tt.willPanic {
				assert.Panics(t, testBlock)
			} else {
				assert.NotPanics(t, testBlock)
			}
		})
	}
}

func TestNull(t *testing.T) {
	var tests = []struct {
		name      string
		expr      string
		vars      Vars
		newErrIs  assert.ErrorAssertionFunc
		execErrIs assert.ErrorAssertionFunc
		results   Results
		willPanic bool
	}{
		{
			name:      "scalar: unary ! null(): is null()",
			expr:      "! null()",
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   NewScalarResults(nil),
		},
		{
			name:      "scalar: binary null() + null(): is nan()", // odd behavior?
			expr:      "null() + null()",
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   NewScalarResults(NaN),
		},
		{
			name:      "scalar: binary 1 + null(): is nan()", // odd behavior?
			expr:      "1 + null()",
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   NewScalarResults(NaN),
		},
	}

	// go-cmp instead of testify assert is used to compare results here
	// because it supports an option for NaN equality.
	// https://github.com/stretchr/testify/pull/691#issuecomment-528457166
	opt := cmp.Comparer(func(x, y float64) bool {
		return (math.IsNaN(x) && math.IsNaN(y)) || x == y
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBlock := func() {
				e, err := New(tt.expr)
				tt.newErrIs(t, err)
				if e != nil {
					res, err := e.Execute(tt.vars)
					tt.execErrIs(t, err)
					if !cmp.Equal(res, tt.results, opt) {
						assert.FailNow(t, tt.name, cmp.Diff(res, tt.results, opt))
					}
				}
			}
			if tt.willPanic {
				assert.Panics(t, testBlock)
			} else {
				assert.NotPanics(t, testBlock)
			}
		})
	}
}
