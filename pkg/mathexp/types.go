package mathexp

import (
	"fmt"
	"sort"
	"time"

	"regexp"
	"strconv"
	"strings"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/mathexp/parse"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"gonum.org/v1/gonum/stat"
)

// Results is a container for Value interfaces.
type Results struct {
	Values Values
}

// Values is a slice of Value interfaces
type Values []Value

// Value is the interface that holds different types such as a Scalar, Series, or Number.
// all Value implementations should be a *data.Frame
type Value interface {
	Type() parse.ReturnType
	Value() interface{}
	GetLabels() data.Labels
	SetLabels(data.Labels)
	GetName() string
	AsDataFrame() *data.Frame
}

// Scalar is the type that holds a single number constant.
// Before returning from an expression it will be wrapped in a
// data frame.
type Scalar struct{ *data.Frame }

// Type returns the Value type and allows it to fulfill the Value interface.
func (s Scalar) Type() parse.ReturnType { return parse.TypeScalar }

// Value returns the actual value allows it to fulfill the Value interface.
func (s Scalar) Value() interface{} { return s }

// AsDataFrame returns the underlying *data.Frame.
func (s Scalar) AsDataFrame() *data.Frame { return s.Frame }

// NewScalar creates a Scalar holding value f.
func NewScalar(f *float64) Scalar {
	return Scalar{
		&data.Frame{
			Fields: data.Fields{
				&data.Field{
					Name:   "Scalar",
					Type:   data.TypeNumber,
					Vector: &data.Float64Vector{f},
				},
			},
		},
	}
}

// NewScalarResults creates a Results holding a single Scalar
func NewScalarResults(f *float64) Results {
	return Results{
		Values: []Value{NewScalar(f)},
	}
}

// GetFloat64Value retrieves the single scalar value from the dataframe
func (s Scalar) GetFloat64Value() *float64 {
	return s.Fields[0].Vector.GetValue(0).(*float64)
}

// Number hold a labelled single number values.
type Number struct{ *data.Frame }

// Type returns the Value type and allows it to fulfill the Value interface.
func (n Number) Type() parse.ReturnType { return parse.TypeNumberSet }

// Value returns the actual value allows it to fulfill the Value interface.
func (n Number) Value() interface{} { return &n }

// AsDataFrame returns the underlying *data.Frame.
func (n Number) AsDataFrame() *data.Frame { return n.Frame }

// NewFields initalizes the Number's fields.
func (n Number) NewFields(metricName string) {
	n.Fields = data.Fields(newNumberFields(metricName))
}

// SetValue sets the value of the Number to float64 pointer f
func (n Number) SetValue(f *float64) {
	n.Fields[0].Vector.SetValue(0, f)
}

// GetFloat64Value retrieves the single scalar value from the dataframe
func (n Number) GetFloat64Value() *float64 {
	return n.Fields[0].Vector.GetValue(0).(*float64)
}

// NewNumber returns a dataframe that holds a float64Vector
func NewNumber(name string, labels data.Labels) Number {
	return Number{&data.Frame{
		Fields: newNumberFields(name),
		Labels: labels,
	}}
}

// newNumberFields creates fields for the Number type.
func newNumberFields(metricName string) data.Fields {
	v := make(data.Float64Vector, 1)
	return data.Fields{
		&data.Field{
			Name:   metricName,
			Type:   data.TypeNumber,
			Vector: &v,
		},
	}
}

// Series has *time.Time and *float64 fields.
type Series struct{ *data.Frame }

// Type returns the Value type and allows it to fulfill the Value interface.
func (s Series) Type() parse.ReturnType { return parse.TypeSeriesSet }

// Value returns the actual value allows it to fulfill the Value interface.
func (s Series) Value() interface{} { return &s }

// AsDataFrame returns the underlying *data.Frame.
func (s Series) AsDataFrame() *data.Frame { return s.Frame }

// GetPoint returns the time and value at the specified index.
func (s Series) GetPoint(pointIdx int) (*time.Time, *float64) {
	return s.GetTime(pointIdx), s.GetValue(pointIdx)
}

// SetPoint sets the time and value on the corresponding vectors at the specified index.
func (s Series) SetPoint(pointIdx int, t *time.Time, f *float64) {
	s.Fields[0].Vector.SetValue(pointIdx, t) // We switch from tsdb's package value,time to time,value
	s.Fields[1].Vector.SetValue(pointIdx, f)
}

// Len returns the length of the series.
func (s Series) Len() int {
	return s.Fields[0].Vector.Len()
}

// GetTime returns the time at the specified index.
func (s Series) GetTime(pointIdx int) *time.Time {
	return s.Fields[0].Vector.GetValue(pointIdx).(*time.Time)
}

// GetValue returns the float value at the specified index.
func (s Series) GetValue(pointIdx int) *float64 {
	return s.Fields[1].Vector.GetValue(pointIdx).(*float64)
}

// NewSeries returns a dataframe of type Series.
func NewSeries(name string, labels data.Labels, len int) Series {
	return Series{&data.Frame{
		Fields: newSeriesFields(name, len),
		Labels: labels,
	}}
}

// seriesFields are the fields expected for a Frame with a Series.
// type seriesFields data.Fields

// newSeriesFields creates fields for the Series type.
func newSeriesFields(metricName string, len int) data.Fields {
	fields := data.Fields{
		&data.Field{
			Name:   "Time",
			Type:   data.TypeTime,
			Vector: &data.TimeVector{},
		},
		&data.Field{
			Name:   metricName,
			Type:   data.TypeNumber,
			Vector: &data.Float64Vector{},
		},
	}
	fields[0].Vector.Make(len)
	fields[1].Vector.Make(len)
	return fields
}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

// Resample turns the Series into a Number based on the given reduction function
func (s Series) Resample(rule string, tr *datasource.TimeRange) (Series, error) {
	AliasToDuration := map[string]time.Duration{
		"D":   86400 * time.Second,
		"W":   604800 * time.Second,
		"MS":  2629800 * time.Second,
		"Y":   31557600 * time.Second,
		"H":   time.Hour,
		"T":   time.Minute,
		"min": time.Minute,
		"S":   time.Second,
		"L":   time.Millisecond,
		"ms":  time.Millisecond,
		"U":   time.Microsecond,
		"us":  time.Microsecond,
		"N":   time.Nanosecond,
	}
	aliases := make([]string, 0)
	for k := range AliasToDuration {
		aliases = append(aliases, k)
	}
	// Use anything other regular expressions?
	expr := strings.Join(aliases, "|")
	re := regexp.MustCompile(fmt.Sprintf(`^(\d*)(%v)$`, expr))
	match := re.FindStringSubmatch(rule)

	if len(match) == 0 {
		// What should I return instead of s?
		return s, fmt.Errorf("resample rule %v not implemented", rule)
	}
	var multiplier int64
	if match[1] != "" {
		valueInt64, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			// Different message for ErrSyntax and ErrRange
			return s, fmt.Errorf("string %v cannot be converted to integer", match[1])
		}
		multiplier = valueInt64
	} else {
		multiplier = 1
	}
	interval := time.Duration(multiplier) * AliasToDuration[match[2]]

	from, err := msToTime(tr.FromRaw)
	if err != nil {
		return s, fmt.Errorf(`failed to parse "from" field "%v": %v`, tr.FromRaw, err)
	}
	to, err := msToTime(tr.ToRaw)
	if err != nil {
		return s, fmt.Errorf(`failed to parse "to" field "%v": %v`, tr.FromRaw, err)
	}

	newSeriesLength := int(float64(to.Sub(from).Nanoseconds()) / float64(interval.Nanoseconds()))
	if newSeriesLength <= 0 {
		return s, fmt.Errorf("The series cannot be sampled further; the time range is shorter than the interval")
	}
	resampled := NewSeries(s.Name, s.Labels, newSeriesLength+1)
	bookmark := 0
	var lastSeen *float64 = nil
	idx := 0
	t := from
	for !t.After(to) && idx <= newSeriesLength {
		values := make([]float64, 0)
		sIdx := bookmark
		for {
			if sIdx == s.Len() {
				break
			}
			st, v := s.GetPoint(sIdx)
			if st.After(t) {
				break
			}
			bookmark++
			sIdx++
			lastSeen = v
			values = append(values, *v)
		}
		var value *float64 = nil
		if len(values) == 0 { // upsampling
			if lastSeen != nil { // only bfill for now
				value = lastSeen
			}
		} else { // downsampling
			tmp := stat.Mean(values, nil) // only mean for now
			value = &tmp
		}
		tv := t // his is required otherwise all points keep the latest timestamp; anything better?
		resampled.SetPoint(idx, &tv, value)
		t = t.Add(interval)
		idx++
	}
	return resampled, nil
}

// Reduce turns the Series into a Number based on the given reduction function
func (s Series) Reduce(rFunc string) (Number, error) {
	// TODO Labels....
	number := NewNumber(fmt.Sprintf("%v_%v", rFunc, s.Name), nil)
	var f *float64
	fVec := s.Fields[1].Vector.(*data.Float64Vector)
	switch rFunc {
	case "sum":
		f = fVec.Sum()
	case "mean":
		f = fVec.Avg()
	case "min":
		f = fVec.Min()
	case "max":
		f = fVec.Max()
	case "count":
		f = fVec.Count()
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
		s := NewSeries(series.Name, data.Labels(series.Tags), len(series.Points))
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

// FromTSDB converts Grafana's existing tsdb.TimeSeriesSlice type to Series stored in a data.FrameCollection
// func FromTSDB(seriesCollection tsdb.TimeSeriesSlice) Results {
// 	results := Results{[]Value{}}
// 	results.Values = make([]Value, len(seriesCollection))
// 	for seriesIdx, series := range seriesCollection {
// 		s := NewSeries(series.Name, data.Labels(series.Tags), len(series.Points))
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
