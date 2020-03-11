package mathexp

import (
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func Sum(v *data.Field) *float64 {
	var sum float64
	for i := 0; i < v.Len(); i++ {
		if f, ok := v.At(i).(*float64); ok {
			sum += *f
		}
	}
	return &sum
}

func Avg(v *data.Field) *float64 {
	sum := Sum(v)
	f := *sum / float64(v.Len())
	return &f
}

func Min(fv *data.Field) *float64 {
	var f float64
	for i := 0; i < fv.Len(); i++ {
		if v, ok := fv.At(i).(*float64); ok {
			if i == 0 || *v < f {
				f = *v
			}
		}
	}
	return &f
}

func Max(fv *data.Field) *float64 {
	var f float64
	for i := 0; i < fv.Len(); i++ {
		if v, ok := fv.At(i).(*float64); ok {
			if i == 0 || *v > f {
				f = *v
			}
		}
	}
	return &f
}

func Count(fv *data.Field) *float64 {
	f := float64(fv.Len())
	return &f
}

// Reduce turns the Series into a Number based on the given reduction function
func (s Series) Reduce(rFunc string) (Number, error) {
	// TODO Labels....
	number := NewNumber(fmt.Sprintf("%v_%v", rFunc, s.GetName()), nil)
	var f *float64
	fVec := s.Frame.Fields[1]
	switch rFunc {
	case "sum":
		f = Sum(fVec)
	case "mean":
		f = Avg(fVec)
	case "min":
		f = Min(fVec)
	case "max":
		f = Max(fVec)
	case "count":
		f = Count(fVec)
	default:
		return number, fmt.Errorf("reduction %v not implemented", rFunc)
	}
	number.SetValue(f)
	return number, nil
}
