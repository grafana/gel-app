package data

import "time"

// Vector is an Interface for the Vectors within a Field of Frame.
type Vector interface {
	// SetValue sets val at idx of the Vector. It will panic if the vector is not initialized, idx is out of range, or
	// val is not of the expected type.
	SetValue(idx int, val interface{})
	// GetValue returns the value at idx of the Vector. It will panic if the vector is not initialized or if idx is out of range.
	GetValue(idx int) interface{}
	// Make initializes the vector and makes of length len.
	Make(len int)
	Len() int
	copy() Vector
}

// Consider New function

// TimeVector is a Vector of time.Time pointers.
type TimeVector []*time.Time

// SetValue sets the value at index idx of the vector.
// It will panic if idx is out of range, the vector is not initialized,
// or t is not a *time.Time.
func (tv *TimeVector) SetValue(idx int, t interface{}) {
	(*tv)[idx] = t.(*time.Time)
}

// GetValue returns the time.Time pointer from the vector at the index idx.
// It will panic if idx is out of range or the vector is not initialized.
func (tv *TimeVector) GetValue(idx int) interface{} {
	return (*tv)[idx]
}

// Make initializes the Vector with length len.
func (tv *TimeVector) Make(len int) {
	*tv = make(TimeVector, len)
}

// Len returns the length of the vector.
// It will panic if the vector is not initialized.
func (tv *TimeVector) Len() int {
	return len(*tv)
}

// Copy returns a copy of the TimeVector that can be mutated without changing the original.
func (tv *TimeVector) Copy() *TimeVector {
	newTV := make(TimeVector, tv.Len())
	for i, v := range *tv {
		newTV[i] = v
	}
	return &newTV
}

func (tv *TimeVector) copy() Vector {
	return tv.Copy()
}

// Float64Vector is a Vector of float64 pointers.
type Float64Vector []*float64

// SetValue sets the value at index idx of the vector.
// It will panic if the idx is out of range, the vector is not initialized,
// or t is not a *float64.
func (fv *Float64Vector) SetValue(idx int, t interface{}) {
	(*fv)[idx] = t.(*float64)
}

// GetValue returns the float64 pointer from the vector at the index idx.
// It will panic if the idx is out of range or if the vector is not initialized.
func (fv *Float64Vector) GetValue(idx int) interface{} {
	return (*fv)[idx]
}

// Make initializes the Vector with length len.
func (fv *Float64Vector) Make(len int) {
	*fv = make(Float64Vector, len)
}

// Len returns the length of the vector.
// It will panic if the vector is not initialized.
func (fv *Float64Vector) Len() int {
	return len(*fv)
}

// Count returns the numbers of points of the vector.
// It will panic if the vector is not initialized.
func (fv *Float64Vector) Count() *float64 {
	f := float64(len(*fv))
	return &f
}

// Copy returns a copy of the Float64Vector that can be mutated without changing the original.
func (fv *Float64Vector) Copy() *Float64Vector {
	newFV := make(Float64Vector, fv.Len())
	for i, v := range *fv {
		newFV[i] = v
	}
	return &newFV
}

func (fv *Float64Vector) copy() Vector {
	return fv.Copy()
}

// Sum returns the sum of all floats in the Vector
// as a float64 pointer. null values will currently trigger a panic.
func (fv *Float64Vector) Sum() *float64 {
	var f float64
	for i := 0; i < fv.Len(); i++ {
		v := fv.GetValue(i).(*float64)
		f += *v
	}
	return &f
}

// Avg returns the arithmetic mean of all floats in the Vector
// a pointer. null values will currently trigger a panic.
func (fv *Float64Vector) Avg() *float64 {
	f := *fv.Sum() / float64(fv.Len())
	return &f
}

// Min returns the min of all floats in the Vector
// a pointer. null values will currently trigger a panic.
func (fv *Float64Vector) Min() *float64 {
	var f float64
	for i := 0; i < fv.Len(); i++ {
		v := *(fv.GetValue(i).(*float64))
		if i == 0 || v < f {
			f = v
		}
	}
	return &f
}

// Max returns the max of all floats in the Vector
// a pointer. null values will currently trigger a panic.
func (fv *Float64Vector) Max() *float64 {
	var f float64
	for i := 0; i < fv.Len(); i++ {
		v := *(fv.GetValue(i).(*float64))
		if i == 0 || v > f {
			f = v
		}
	}
	return &f
}

func ToFloat64Vector(vals []float64) *Float64Vector {
	fVec := Float64Vector{}
	(&fVec).Make(len(vals))
	for i, v := range vals {
		tmp := v
		(&fVec).SetValue(i, &tmp)
	}
	return &fVec
}
