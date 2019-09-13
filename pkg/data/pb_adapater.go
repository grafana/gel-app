package data

import (
	"fmt"
	"time"

	st "github.com/golang/protobuf/ptypes/struct"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

// ToPBFrame Converts that dataframe to our protobuf (for grpc)
// type.
// Warning: Both these types are subject to substantial change
// as the protocol is only in a branch.
func (f *Frame) ToPBFrame() (*datasource.Frame, error) {
	dsFrame := &datasource.Frame{
		Name:   f.Name,
		Labels: f.Labels,
		RefId:  f.RefID,
		Fields: make([]*datasource.Field, len(f.Fields)),
	}
	for i, field := range f.Fields {

		values := make([]*st.Value, field.Vector.Len())

		for vIdx := 0; vIdx < field.Vector.Len(); vIdx++ {
			val := field.Vector.GetValue(vIdx)
			pbVal, err := toPBVal(val)
			if err != nil {
				return nil, err
			}
			values[vIdx] = pbVal
		}

		listValue := st.ListValue{
			Values: values,
		}
		dsFrame.Fields[i] = &datasource.Field{
			Name:   field.Name,
			Type:   datasource.Field_FieldType(field.Type),
			Values: &listValue,
		}
	}
	return dsFrame, nil
}

func toPBVal(v interface{}) (*st.Value, error) {
	switch v := v.(type) {
	case nil:
		return nil, nil
	case *float64:
		return &st.Value{
			Kind: &st.Value_NumberValue{
				NumberValue: *v,
			},
		}, nil
	case *time.Time:
		return &st.Value{
			Kind: &st.Value_StringValue{
				StringValue: fmt.Sprintf("%v", v.UnixNano()),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type %T for frame to protobuf conversion", v)

	}
}
