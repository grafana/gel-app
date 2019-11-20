package mathexp

import (
	"fmt"
	"sort"
	"time"

	"github.com/grafana/gel-app/pkg/mathexp/parse"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

// Series has time.Time and ...? *float64 fields.
type Series struct {
	Frame          *dataframe.Frame
	TimeIsNullable bool
	// TODO:
	// - Multiple Value Fields
	// - Field Index of time column doesn't matter
	// - Validate Method?
	// - Value can be different number types
}

// SeriesFromFrame validates that the dataframe can be considered a Series type
// and populate meta information on Series about the frame.
func SeriesFromFrame(frame *dataframe.Frame) (s Series, err error) {
	if len(frame.Fields) != 2 {
		return s, fmt.Errorf("frame must have two fields to be a series, has %v", len(frame.Fields))
	}

	switch frame.Fields[0].Vector.PrimitiveType() {
	case dataframe.VectorPTypeTime:
		s.TimeIsNullable = false
	case dataframe.VectorPTypeNullableTime:
		s.TimeIsNullable = true
	default:
		return s, fmt.Errorf("first column of dataframe must be []time.Time or []*time.Time")
	}
	s.Frame = frame
	return
}

// NewSeries returns a dataframe of type Series.
func NewSeries(name string, labels dataframe.Labels, nullableTime bool, size int) Series {
	var timeColumn *dataframe.Field
	if nullableTime {
		timeColumn = dataframe.NewField("Time", make([]*time.Time, size))
	} else {
		timeColumn = dataframe.NewField("Time", make([]time.Time, size))
	}
	return Series{
		Frame: dataframe.New("", labels,
			timeColumn,
			dataframe.NewField(name, make([]*float64, size)),
		),
		TimeIsNullable: nullableTime,
	}
}

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
func (s Series) SetPoint(pointIdx int, t *time.Time, f *float64) (err error) {
	if s.TimeIsNullable {
		s.Frame.Fields[0].Vector.Set(pointIdx, t)
	} else {
		if t == nil {
			return fmt.Errorf("can not set null time value on non-nullable time field for series name %v", s.Frame.Name)
		}
		s.Frame.Fields[0].Vector.Set(pointIdx, *t)
	}
	s.Frame.Fields[1].Vector.Set(pointIdx, f)
	return
}

// AppendPoint appends a point (time/value).
func (s Series) AppendPoint(pointIdx int, t *time.Time, f *float64) (err error) {
	if s.TimeIsNullable {
		s.Frame.Fields[0].Vector.Append(t)
	} else {
		if t == nil {
			return fmt.Errorf("can not append null time value on non-nullable time field for series name %v", s.Frame.Name)
		}
		s.Frame.Fields[0].Vector.Append(*t)

	}
	s.Frame.Fields[1].Vector.Append(f)
	return
}

// Len returns the length of the series.
func (s Series) Len() int {
	return s.Frame.Fields[0].Vector.Len()
}

// GetTime returns the time at the specified index.
func (s Series) GetTime(pointIdx int) *time.Time {
	if s.TimeIsNullable {
		return s.Frame.Fields[0].Vector.At(pointIdx).(*time.Time)
	}
	t := s.Frame.Fields[0].Vector.At(pointIdx).(time.Time)
	return &t
}

// GetValue returns the float value at the specified index.
func (s Series) GetValue(pointIdx int) *float64 {
	return s.Frame.Fields[1].Vector.At(pointIdx).(*float64)
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