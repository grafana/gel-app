package mathexp

import (
	"math"

	"github.com/grafana/gel-app/pkg/mathexp/parse"
)

var builtins = map[string]parse.Func{
	"abs": {
		Args:          []parse.ReturnType{parse.TypeVariantSet},
		VariantReturn: true,
		F:             abs,
	},
	"log": {
		Args:          []parse.ReturnType{parse.TypeVariantSet},
		VariantReturn: true,
		F:             log,
	},
	"nan": {
		Return: parse.TypeScalar,
		F:      nan,
	},
	"inf": {
		Return: parse.TypeScalar,
		F:      inf,
	},
	"null": {
		Return: parse.TypeScalar,
		F:      null,
	},
}

// abs returns the absolute value for each result in NumberSet, SeriesSet, or Scalar
func abs(e *State, varSet Results) Results {
	newRes := Results{}
	for _, res := range varSet.Values {
		newVal := perFloat(res, math.Abs)
		newRes.Values = append(newRes.Values, newVal)
	}
	return newRes
}

// log returns the natural logarithm value for each result in NumberSet, SeriesSet, or Scalar
func log(e *State, varSet Results) Results {
	newRes := Results{}
	for _, res := range varSet.Values {
		newVal := perFloat(res, math.Log)
		newRes.Values = append(newRes.Values, newVal)
	}
	return newRes
}

// nan returns a scalar nan value
func nan(e *State) Results {
	aNaN := math.NaN()
	return NewScalarResults(&aNaN)
}

// inf returns a scalar positive infinity value
func inf(e *State) Results {
	aInf := math.Inf(0)
	return NewScalarResults(&aInf)
}

// null returns a null scalar value
func null(e *State) Results {
	return NewScalarResults(nil)
}

func perFloat(val Value, floatF func(x float64) float64) Value {
	var newVal Value
	switch val.Type() {
	case parse.TypeNumberSet:
		n := NewNumber(val.GetName(), val.GetLabels())
		f := val.(Number).GetFloat64Value()
		nF := math.NaN()
		if f != nil {
			nF = floatF(*f)
		}
		n.SetValue(&nF)
		newVal = n
	case parse.TypeScalar:
		f := val.(Scalar).GetFloat64Value()
		nF := math.NaN()
		if f != nil {
			nF = floatF(*f)
		}
		newVal = NewScalar(&nF)
	case parse.TypeSeriesSet:
		resSeries := val.(Series)
		newSeries := NewSeries(resSeries.GetName(), resSeries.GetLabels(), resSeries.TimeIdx, resSeries.TimeIsNullable, resSeries.ValueIdx, resSeries.Len())
		for i := 0; i < resSeries.Len(); i++ {
			t, f := resSeries.GetPoint(i)
			nF := math.NaN()
			if f != nil {
				nF = floatF(*f)
			}
			newSeries.SetPoint(i, t, &nF)
		}
		newVal = newSeries
	}
	return newVal
}
