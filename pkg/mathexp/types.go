package mathexp

import (
	"github.com/grafana/gel-app/pkg/mathexp/parse"
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

func (s Scalar) GetLabels() dataframe.Labels { return nil }

func (s Scalar) SetLabels(ls dataframe.Labels) { return }

func (s Scalar) GetName() string { return s.Frame.Name }

// AsDataFrame returns the underlying *dataframe.Frame.
func (s Scalar) AsDataFrame() *dataframe.Frame { return s.Frame }

// NewScalar creates a Scalar holding value f.
func NewScalar(f *float64) Scalar {
	frame := dataframe.New("",
		dataframe.NewField("Scalar", nil, []*float64{f}),
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

func (n Number) GetLabels() dataframe.Labels { return n.Frame.Fields[0].Labels }

func (n Number) SetLabels(ls dataframe.Labels) { n.Frame.Fields[0].Labels = ls }

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
		dataframe.New("",
			dataframe.NewField(name, labels, make([]*float64, 1)),
		),
	}
}
