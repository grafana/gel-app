package mathexp

import (
	"fmt"
	"sort"
	"time"

	"github.com/grafana/gel-app/pkg/mathexp/parse"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

// Results is a container for Value interfaces.
type Results struct {
	Values Values
}

// Values is a slice of Value interfaces
type Values []Value

// Value is the interface that holds different types such as a Scalar, Series, or Number.
// all Value implementations should be a *dataframe.Frame
type Value interface {
	Type() parse.ReturnType
	Value() interface{}
	GetLabels() dataframe.Labels
	SetLabels(dataframe.Labels)
	GetName() string
	AsDataFrame() *dataframe.Frame
}

// Scalar is the type that holds a single number constant.
// Before returning from an expression it will be wrapped in a
// data frame.
type Scalar struct{ Frame *dataframe.Frame }

// Type returns the Value type and allows it to fulfill the Value interface.
func (s Scalar) Type() parse.ReturnType { return parse.TypeScalar }

// Value returns the actual value allows it to fulfill the Value interface.
func (s Scalar) Value() interface{} { return s }

func (s Scalar) GetLabels() dataframe.Labels { return s.Frame.Labels }

func (s Scalar) SetLabels(ls dataframe.Labels) { s.Frame.Labels = ls }

func (s Scalar) GetName() string { return s.Frame.Name }

// AsDataFrame returns the underlying *dataframe.Frame.
func (s Scalar) AsDataFrame() *dataframe.Frame { return s.Frame }

// NewScalar creates a Scalar holding value f.
func NewScalar(f *float64) Scalar {
	frame := dataframe.New("", nil,
		dataframe.NewField("Scalar", dataframe.FieldTypeNumber, []*float64{f}),
	)
	return Scalar{frame}
}

// NewScalarResults creates a Results holding a single Scalar
func NewScalarResults(f *float64) Results {
	return Results{
		Values: []Value{NewScalar(f)},
	}
}

// GetFloat64Value retrieves the single scalar value from the dataframe
func (s Scalar) GetFloat64Value() *float64 {
	return s.Frame.Fields[0].Vector.At(0).(*float64)
}

// Number hold a labelled single number values.
type Number struct{ Frame *dataframe.Frame }

// Type returns the Value type and allows it to fulfill the Value interface.
func (n Number) Type() parse.ReturnType { return parse.TypeNumberSet }

// Value returns the actual value allows it to fulfill the Value interface.
func (n Number) Value() interface{} { return &n }

func (n Number) GetLabels() dataframe.Labels { return n.Frame.Labels }

func (n Number) SetLabels(ls dataframe.Labels) { n.Frame.Labels = ls }

func (n Number) GetName() string { return n.Frame.Name }

// AsDataFrame returns the underlying *dataframe.Frame.
func (n Number) AsDataFrame() *dataframe.Frame { return n.Frame }

// SetValue sets the value of the Number to float64 pointer f
func (n Number) SetValue(f *float64) {
	n.Frame.Fields[0].Vector.Set(0, f)
}

// GetFloat64Value retrieves the single scalar value from the dataframe
func (n Number) GetFloat64Value() *float64 {
	return n.Frame.Fields[0].Vector.At(0).(*float64)
}

// NewNumber returns a dataframe that holds a float64Vector
func NewNumber(name string, labels dataframe.Labels) Number {
	return Number{
		dataframe.New("", labels,
			dataframe.NewField(name, dataframe.FieldTypeNumber, make([]*float64, 1)),
		),
	}
}

// Series has *time.Time and *float64 fields.
type Series struct{ Frame *dataframe.Frame }

// Type returns the Value type and allows it to fulfill the Value interface.
func (s Series) Type() parse.ReturnType { return parse.TypeSeriesSet }

// Value returns the actual value allows it to fulfill the Value interface.
func (s Series) Value() interface{} { return &s }

func (s Series) GetLabels() dataframe.Labels { return s.Frame.Labels }

func (s Series) SetLabels(ls dataframe.Labels) { s.Frame.Labels = ls }

func (s Series) GetName() string { return s.Frame.Name }

// AsDataFrame returns the underlying *dataframe.Frame.
func (s Series) AsDataFrame() *dataframe.Frame { return s.Frame }

// GetPoint returns the time and value at the specified index.
func (s Series) GetPoint(pointIdx int) (*time.Time, *float64) {
	return s.GetTime(pointIdx), s.GetValue(pointIdx)
}

// SetPoint sets the time and value on the corresponding vectors at the specified index.
func (s Series) SetPoint(pointIdx int, t *time.Time, f *float64) {
	s.Frame.Fields[0].Vector.Set(pointIdx, t) // We switch from tsdb's package value,time to time,value
	s.Frame.Fields[1].Vector.Set(pointIdx, f)
}

// Len returns the length of the series.
func (s Series) Len() int {
	return s.Frame.Fields[0].Vector.Len()
}

// GetTime returns the time at the specified index.
func (s Series) GetTime(pointIdx int) *time.Time {
	return s.Frame.Fields[0].Vector.At(pointIdx).(*time.Time)
}

// GetValue returns the float value at the specified index.
func (s Series) GetValue(pointIdx int) *float64 {
	return s.Frame.Fields[1].Vector.At(pointIdx).(*float64)
}

// NewSeries returns a dataframe of type Series.
func NewSeries(name string, labels dataframe.Labels, size int) Series {
	return Series{
		dataframe.New("", labels,
			dataframe.NewField("Time", dataframe.FieldTypeTime, make([]*time.Time, size)),
			dataframe.NewField(name, dataframe.FieldTypeNumber, make([]*float64, size)),
		),
	}
}

func Sum(v dataframe.Vector) *float64 {
	var sum float64
	for i := 0; i < v.Len(); i++ {
		if f, ok := v.At(i).(*float64); ok {
			sum += *f
		}
	}
	return &sum
}

func Avg(v dataframe.Vector) *float64 {
	sum := Sum(v)
	f := *sum / float64(v.Len())
	return &f
}

func Min(fv dataframe.Vector) *float64 {
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

func Max(fv dataframe.Vector) *float64 {
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

func Count(fv dataframe.Vector) *float64 {
	f := float64(fv.Len())
	return &f
}

// Reduce turns the Series into a Number based on the given reduction function
func (s Series) Reduce(rFunc string) (Number, error) {
	// TODO Labels....
	number := NewNumber(fmt.Sprintf("%v_%v", rFunc, s.GetName()), nil)
	var f *float64
	fVec := s.Frame.Fields[1].Vector
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

// FromGRPC converts time series only (at the moment) from a
// GRPC TimeSeries type to a Series Type
func FromGRPC(seriesCollection []*datasource.TimeSeries) Results {
	results := Results{[]Value{}}
	results.Values = make([]Value, len(seriesCollection))
	for seriesIdx, series := range seriesCollection {
		s := NewSeries(series.Name, dataframe.Labels(series.Tags), len(series.Points))
		for pointIdx, point := range series.Points {
			t, f := convertDSTimePoint(point)
			s.SetPoint(pointIdx, t, f)
		}
		results.Values[seriesIdx] = s
	}
	return results
}

func convertDSTimePoint(point *datasource.Point) (t *time.Time, f *float64) {
	tI := int64(point.Timestamp)
	uT := time.Unix(tI/int64(1e+3), (tI%int64(1e+3))*int64(1e+6)) // time.Time from millisecond unix ts
	t = &uT
	f = &point.Value
	return t, f
}

// SortByTime sorts the series by the time from oldest to newest.
// If desc is true, it will sort from newest to oldest.
// If any time values are nil, it will panic.
func (s Series) SortByTime(desc bool) {
	if desc {
		sort.Sort(sort.Reverse(SortSeriesByTime(s)))
		return
	}
	sort.Sort(SortSeriesByTime(s))
}

// SortSeriesByTime allows a Series to be sorted by time
// the sort interface will panic if any timestamps are null
type SortSeriesByTime Series

func (ss SortSeriesByTime) Len() int { return Series(ss).Len() }

func (ss SortSeriesByTime) Swap(i, j int) {
	iTimeVal, iFVal := Series(ss).GetPoint(i)
	jTimeVal, jFVal := Series(ss).GetPoint(j)
	Series(ss).SetPoint(j, iTimeVal, iFVal)
	Series(ss).SetPoint(i, jTimeVal, jFVal)
}

func (ss SortSeriesByTime) Less(i, j int) bool {
	iTimeVal := Series(ss).GetTime(i)
	jTimeVal := Series(ss).GetTime(j)
	return iTimeVal.Before(*jTimeVal)
}

// Kept for reference of difference Series types, commented out to avoid import

// FromTSDB converts Grafana's existing tsdb.TimeSeriesSlice type to Series stored in a dataframe.FrameCollection
// func FromTSDB(seriesCollection tsdb.TimeSeriesSlice) Results {
// 	results := Results{[]Value{}}
// 	results.Values = make([]Value, len(seriesCollection))
// 	for seriesIdx, series := range seriesCollection {
// 		s := NewSeries(series.Name, dataframe.Labels(series.Tags), len(series.Points))
// 		for pointIdx, point := range series.Points {
// 			t, f := convertTSDBTimePoint(point)
// 			s.SetPoint(pointIdx, t, f)
// 		}
// 		results.Values[seriesIdx] = s
// 	}
// 	return results
// }

// convertTSDBTimePoint coverts a tsdb.TimePoint into two values appropriate
// for Series values.
// func convertTSDBTimePoint(point tsdb.TimePoint) (t *time.Time, f *float64) {
// 	timeIdx, valueIdx := 1, 0
// 	if point[timeIdx].Valid { // Assuming valid is null?
// 		tI := int64(point[timeIdx].Float64)
// 		uT := time.Unix(tI/int64(1e+3), (tI%int64(1e+3))*int64(1e+6)) // time.Time from millisecond unix ts
// 		t = &uT
// 	}
// 	if point[valueIdx].Valid {
// 		f = &point[valueIdx].Float64
// 	}
// 	return
// }
