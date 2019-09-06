package data

import (
	"bytes"
	"sort"
	"strings"
)

// TODO: think about what should be pointers in terms of allocations

// Frames is a collection of Frame
type Frames []*Frame

// Frame is a columnar container for data.
type Frame struct {
	Name   string `json:"name"`
	Fields Fields `json:"fields"`
	Labels Labels `json:"labels"`
	RefID  string `json:"refId"`
	// GrafanaType
}

// GetLabels returns the Labels for the Frame.
func (f *Frame) GetLabels() Labels {
	return f.Labels
}

// GetName returns the name for the Frame.
func (f *Frame) GetName() string {
	return f.Name
}

// SetLabels returns the Labels for the Frame.
func (f *Frame) SetLabels(l Labels) {
	f.Labels = l
}

// Copy returns a new copy of the Frame that can be mutated without
// changing the original.
func (f *Frame) Copy() *Frame {
	newFrame := &Frame{
		Name:  f.Name,
		RefID: f.RefID,
	}
	if f.Fields != nil {
		newFrame.Fields = f.Fields.Copy()
	}
	if f.Labels != nil {
		newFrame.Labels = f.Labels.Copy()
	}
	return newFrame
}

// Fields is a collection of Field.
type Fields []*Field

// Copy returns a new copy of the Fields that can be mutated without
// changing the original.
func (f *Fields) Copy() Fields {
	newFields := make(Fields, len(*f))
	for i, field := range *f {
		newFields[i] = field.Copy()
	}
	return newFields
}

// Field holds information and Type and its values.
type Field struct {
	Name   string    `json:"name"`
	Type   FieldType `json:"type"`
	Vector Vector    `json:"values"` // KMB thinks this should be vectors
}

// Copy returns a new copy of the Field that can be mutated without
// changing the original.
func (f *Field) Copy() *Field {
	newField := Field{
		Name: f.Name,
		Type: f.Type,
	}
	if f.Vector != nil {
		newField.Vector = f.Vector.copy()
	}
	return &newField
}

// FieldType is the type of the data for a Fields Vector
type FieldType int

const (
	// TypeOther is for a FieldType of unknown or other type.
	TypeOther FieldType = iota
	// TypeNumber is a FieldType for number types, and is currently float64 pointer values.
	TypeNumber
	// TypeString is a FieldType for string values.
	TypeString
	// TypeBoolean is a FieldType for string values.
	TypeBoolean
	// TypeTime is a FieldType for time value types, and is currently represented as time.Time pointer values.
	TypeTime
)

func (f FieldType) String() string {
	switch f {
	case TypeOther:
		return "other"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeBoolean:
		return "boolean"
	case TypeTime:
		return "time"
	default:
		return "unknown"
	}
}

// MarshalJSON makes the FieldType a string when represented as JSON.
func (f FieldType) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(f.String())
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// Labels are an identifier attached to a Frame.
type Labels map[string]string

func (l Labels) String() string {
	// Better structure, should be sorted, copy prom probably
	keys := make([]string, len(l))
	i := 0
	for k := range l {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	sb := strings.Builder{}
	i = 0
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(l[k])
		if i != len(keys)-1 {
			sb.WriteString(", ")
		}
		i++
	}
	return sb.String()
}

// Equal returns true if the argument has the same k=v pairs as the receiver.
func (l Labels) Equal(arg Labels) bool {
	if len(l) != len(arg) {
		return false
	}
	for k, v := range l {
		if argVal, ok := arg[k]; !ok || argVal != v {
			return false
		}
	}
	return true
}

// Subset returns true if all k=v pairs of the argument are in the receiver.
func (l Labels) Subset(arg Labels) bool {
	if len(arg) > len(l) {
		return false
	}
	for k, v := range arg {
		if argVal, ok := l[k]; !ok || argVal != v {
			return false
		}
	}
	return true
}

// Copy returns a copy of the Labels
func (l Labels) Copy() Labels {
	newLabels := make(Labels, len(l))
	for k, v := range l {
		newLabels[k] = v
	}
	return newLabels
}
