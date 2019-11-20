package mathexp

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
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
		results   Results
	}{
		{
			name:      "unary series",
			expr:      "! ! $A",
			vars:      aSeries,
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
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
			results: Results{
				[]Value{
					makeSeries("id=1", dataframe.Labels{"id": "1"}, tp{
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
			results: Results{
				[]Value{
					makeSeries("id=1", dataframe.Labels{"id": "1"}, tp{
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
			results: Results{
				[]Value{
					makeSeries("sensor=a, turbine=1", dataframe.Labels{"sensor": "a", "turbine": "1"}, tp{
						unixTimePointer(5, 0), float64Pointer(6 * .5),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(8 * .2),
					}),
					makeSeries("sensor=b, turbine=1", dataframe.Labels{"sensor": "b", "turbine": "1"}, tp{
						unixTimePointer(5, 0), float64Pointer(10 * .5),
					}, tp{
						unixTimePointer(10, 0), float64Pointer(16 * .2),
					}),
				},
			},
		},
		// Length of resulting series is A when A + B. However, only points where the time matches
		// for A and B are added to the result
		{
			name: "series Op series with sparse time join",
			expr: "$A + $B",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("temp", dataframe.Labels{}, tp{
							unixTimePointer(5, 0), float64Pointer(1),
						}, tp{
							unixTimePointer(10, 0), float64Pointer(2),
						}),
					},
				},
				"B": Results{
					[]Value{
						makeSeries("efficiency", dataframe.Labels{}, tp{
							unixTimePointer(5, 0), float64Pointer(3),
						}, tp{
							unixTimePointer(9, 0), float64Pointer(4),
						}),
					},
				},
			},

			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{ // Not sure about preserving names...
						unixTimePointer(5, 0), float64Pointer(4),
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
				if diff := cmp.Diff(tt.results, res, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("Result mismatch (-want +got):\n%s", diff)
				}
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
			makeNumber("volt", dataframe.Labels{"id": "1"}, float64Pointer(7)),
		},
	},
}

var twoSeriesSets = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", dataframe.Labels{"sensor": "a", "turbine": "1"}, tp{
				unixTimePointer(5, 0), float64Pointer(6),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(8),
			}),
			makeSeries("temp", dataframe.Labels{"sensor": "b", "turbine": "1"}, tp{
				unixTimePointer(5, 0), float64Pointer(10),
			}, tp{
				unixTimePointer(10, 0), float64Pointer(16),
			}),
		},
	},
	"B": Results{
		[]Value{
			makeSeries("efficiency", dataframe.Labels{"turbine": "1"}, tp{
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

func makeSeries(name string, labels dataframe.Labels, points ...tp) Series {
	newSeries := NewSeries(name, labels, true, len(points))
	for idx, p := range points {
		newSeries.SetPoint(idx, p.t, p.f)
	}
	return newSeries
}

func makeNumber(name string, labels dataframe.Labels, f *float64) Number {
	newNumber := NewNumber(name, labels)
	newNumber.SetValue(f)
	return newNumber
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

func TestNullValues(t *testing.T) {
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
			name:      "scalar: unary ! null(): is null",
			expr:      "! null()",
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   NewScalarResults(nil),
		},
		{
			name:      "scalar: binary null() + null(): is null",
			expr:      "null() + null()",
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   NewScalarResults(nil),
		},
		{
			name:      "scalar: binary 1 + null(): is null",
			expr:      "1 + null()",
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results:   NewScalarResults(nil),
		},
		{
			name: "series: unary with a null value in it has a null value in result",
			expr: "- $A",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("", nil, tp{
							unixTimePointer(5, 0), float64Pointer(1),
						}, tp{
							unixTimePointer(10, 0), nil,
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
						unixTimePointer(10, 0), nil,
					}),
				},
			},
		},
		{
			name: "series: binary with a null value in it has a null value in result",
			expr: "$A - $A",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("", nil, tp{
							unixTimePointer(5, 0), float64Pointer(1),
						}, tp{
							unixTimePointer(10, 0), nil,
						}),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), float64Pointer(0),
					}, tp{
						unixTimePointer(10, 0), nil,
					}),
				},
			},
		},
		{
			name: "series and scalar: binary with a null value in it has a nil value in result",
			expr: "$A - 1",
			vars: Vars{
				"A": Results{
					[]Value{
						makeSeries("", nil, tp{
							unixTimePointer(5, 0), float64Pointer(1),
						}, tp{
							unixTimePointer(10, 0), nil,
						}),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), float64Pointer(0),
					}, tp{
						unixTimePointer(10, 0), nil,
					}),
				},
			},
		},
		{
			name: "number: unary ! null number: is null",
			expr: "! $A",
			vars: Vars{
				"A": Results{
					[]Value{
						makeNumber("", nil, nil),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeNumber("", nil, nil),
				},
			},
		},
		{
			name: "number: binary null number and null number: is null",
			expr: "$A + $A",
			vars: Vars{
				"A": Results{
					[]Value{
						makeNumber("", nil, nil),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeNumber("", nil, nil),
				},
			},
		},
		{
			name: "number: binary non-null number and null number: is null",
			expr: "$A * $B",
			vars: Vars{
				"A": Results{
					[]Value{
						makeNumber("", nil, nil),
					},
				},
				"B": Results{
					[]Value{
						makeNumber("", nil, float64Pointer(1)),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeNumber("", nil, nil),
				},
			},
		},
		{
			name: "number and series: binary non-null number and series with a null: is null",
			expr: "$A * $B",
			vars: Vars{
				"A": Results{
					[]Value{
						makeNumber("", nil, float64Pointer(1)),
					},
				},
				"B": Results{
					[]Value{
						makeSeries("", nil, tp{
							unixTimePointer(5, 0), float64Pointer(1),
						}, tp{
							unixTimePointer(10, 0), nil,
						}),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), float64Pointer(1),
					}, tp{
						unixTimePointer(10, 0), nil,
					}),
				},
			},
		},
		{
			name: "number and series: binary null number and series with non-null and null: is null and null",
			expr: "$A * $B",
			vars: Vars{
				"A": Results{
					[]Value{
						makeNumber("", nil, nil),
					},
				},
				"B": Results{
					[]Value{
						makeSeries("", nil, tp{
							unixTimePointer(5, 0), float64Pointer(1),
						}, tp{
							unixTimePointer(10, 0), nil,
						}),
					},
				},
			},
			newErrIs:  assert.NoError,
			execErrIs: assert.NoError,
			results: Results{
				[]Value{
					makeSeries("", nil, tp{
						unixTimePointer(5, 0), nil,
					}, tp{
						unixTimePointer(10, 0), nil,
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
					if diff := cmp.Diff(tt.results, res, opt); diff != "" {
						t.Errorf("Result mismatch (-want +got):\n%s", diff)
					}
				}
			}
			if tt.willPanic {
				assert.Panics(t, testBlock)
			} else {
				testBlock()
			}
		})
	}
}
