package mathexp

import (
	"math"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

var aSeriesNullableTime = Vars{
	"A": Results{
		[]Value{
			makeSeriesNullableTime("temp", nil, nullTimeTP{
				unixTimePointer(5, 0), float64Pointer(2),
			}, nullTimeTP{
				unixTimePointer(10, 0), float64Pointer(1),
			}),
		},
	},
}

var aSeries = Vars{
	"A": Results{
		[]Value{
			makeSeries("temp", nil, tp{
				time.Unix(5, 0), float64Pointer(2),
			}, tp{
				time.Unix(10, 0), float64Pointer(1),
			}),
		},
	},
}

var aSeriesbNumber = Vars{
	"A": Results{
		[]Value{
			makeSeriesNullableTime("temp", nil, nullTimeTP{
				unixTimePointer(5, 0), float64Pointer(2),
			}, nullTimeTP{
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
			makeSeriesNullableTime("temp", dataframe.Labels{"sensor": "a", "turbine": "1"}, nullTimeTP{
				unixTimePointer(5, 0), float64Pointer(6),
			}, nullTimeTP{
				unixTimePointer(10, 0), float64Pointer(8),
			}),
			makeSeriesNullableTime("temp", dataframe.Labels{"sensor": "b", "turbine": "1"}, nullTimeTP{
				unixTimePointer(5, 0), float64Pointer(10),
			}, nullTimeTP{
				unixTimePointer(10, 0), float64Pointer(16),
			}),
		},
	},
	"B": Results{
		[]Value{
			makeSeriesNullableTime("efficiency", dataframe.Labels{"turbine": "1"}, nullTimeTP{
				unixTimePointer(5, 0), float64Pointer(.5),
			}, nullTimeTP{
				unixTimePointer(10, 0), float64Pointer(.2),
			}),
		},
	},
}

// NaN is just to make the calls a little cleaner, the one
// call is not for any sort of equality side effect in tests.
// note: cmp.Equal must be used to test Equality for NaNs.
var NaN = float64Pointer(math.NaN())
