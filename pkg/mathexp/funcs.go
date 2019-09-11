package mathexp

import (
	"math"

	"github.com/grafana/gel-app/pkg/mathexp/parse"
)

var builtins = map[string]parse.Func{
	"abs": {
		Args:          []parse.ReturnType{parse.TypeVariantSet},
		VariantReturn: true,
		F:             Abs,
	},
	"log": {
		Args:          []parse.ReturnType{parse.TypeVariantSet},
		VariantReturn: true,
		F:             Log,
	},
}

// Abs returns the absolute value for each result in NumberSet, SeriesSet, or Scalar
func Abs(e *State, varSet Results) Results {
	newRes := Results{}
	for _, res := range varSet.Values {
		newVal := perFloat(res, math.Abs)
		newRes.Values = append(newRes.Values, newVal)
	}
	return newRes
}

// Log returns the natural logarithm value for each result in NumberSet, SeriesSet, or Scalar
func Log(e *State, varSet Results) Results {
	newRes := Results{}
	for _, res := range varSet.Values {
		newVal := perFloat(res, math.Log)
		newRes.Values = append(newRes.Values, newVal)
	}
	return newRes
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
		newSeries := NewSeries(resSeries.Name, resSeries.Labels, resSeries.Len())
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