package data

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/csv"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	st "github.com/golang/protobuf/ptypes/struct"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	"github.com/mattetti/filebuffer"
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

func (f *Frame) ToArrow() ([]byte, error) {
	arrowFields := make([]arrow.Field, len(f.Fields))
	for i, field := range f.Fields {
		at, err := field.Type.ArrowType()
		if err != nil {
			return nil, err
		}
		mdMap := map[string]string{
			"name": field.Name,
			"type": field.Type.String(),
		}
		arrowFields[i] = arrow.Field{Name: field.Name, Type: at, Metadata: arrow.MetadataFrom(mdMap), Nullable: true}
	}
	tableMDMap := map[string]string{
		"name":  f.Name,
		"refId": f.RefID,
	}
	if f.GetLabels() != nil {
		tableMDMap["labels"] = f.Labels.String()
	}

	md := arrow.MetadataFrom(tableMDMap)
	schema := arrow.NewSchema(arrowFields, &md)
	pool := memory.NewGoAllocator()
	columns := make([]array.Column, len(f.Fields))
	for fieldIdx, field := range f.Fields {
		switch field.Type {
		case TypeNumber:
			builder := array.NewFloat64Builder(pool)
			defer builder.Release()
			for _, v := range *field.Vector.(*Float64Vector) {
				if v == nil {
					builder.AppendNull()
					continue
				}
				builder.Append(*v)
			}
			chunked := array.NewChunked(arrowFields[fieldIdx].Type, []array.Interface{builder.NewArray()})

			columns[fieldIdx] = *array.NewColumn(arrowFields[fieldIdx], chunked)
			builder.Release()
		case TypeTime:
			builder := array.NewStringBuilder(pool)
			defer builder.Release()
			for _, v := range *field.Vector.(*TimeVector) {
				if v == nil {
					builder.AppendNull()
					continue
				}
				builder.Append(fmt.Sprintf("%v", v.UnixNano()))
			}
			chunked := array.NewChunked(arrowFields[fieldIdx].Type, []array.Interface{builder.NewArray()})

			columns[fieldIdx] = *array.NewColumn(arrowFields[fieldIdx], chunked)
			builder.Release()
		default:
			return nil, fmt.Errorf("unsupported field type %s for arrow converstion", field.Type)
		}
	}
	table := array.NewTable(schema, columns, -1)
	tableReader := array.NewTableReader(table, -1)
	defer tableReader.Release()

	b := bytes.Buffer{}
	writer := ipc.NewWriter(&b, ipc.WithSchema(tableReader.Schema()))
	defer writer.Close()

	outFile, err := os.OpenFile("/home/kbrandt/tmp/arrowstuff", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	CSVOutFile, err := os.OpenFile("/home/kbrandt/tmp/csvarrowstuff", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	fWriter, err := ipc.NewFileWriter(outFile, ipc.WithSchema(tableReader.Schema()))
	if err != nil {
		return nil, err
	}

	fb := filebuffer.New(nil)

	fakeFWriter, err := ipc.NewFileWriter(fb, ipc.WithSchema(tableReader.Schema()))
	if err != nil {
		return nil, err
	}
	defer fakeFWriter.Close()

	csvWriter := csv.NewWriter(CSVOutFile, schema, csv.WithHeader())

	for tableReader.Next() {
		rec := tableReader.Record()
		err := writer.Write(rec)
		if err != nil {
			return nil, err
		}
		err = fWriter.Write(rec)
		if err != nil {
			return nil, err
		}
		err = fakeFWriter.Write(rec)
		if err != nil {
			return nil, err
		}
		err = csvWriter.Write(rec)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	err = fWriter.Close()
	if err != nil {
		return nil, err
	}

	err = fakeFWriter.Close()
	if err != nil {
		return nil, err
	}

	err = CSVOutFile.Close()
	if err != nil {
		return nil, err
	}
	return fb.Buff.Bytes(), nil
}

func (f FieldType) ArrowType() (dt arrow.DataType, err error) {
	switch f {
	case TypeString:
		dt = &arrow.StringType{}
	case TypeNumber:
		dt = &arrow.Float64Type{}
	case TypeTime:
		dt = &arrow.StringType{} // we encode time as strings for now
	default:
		return dt, fmt.Errorf("unsupported type %s for arrow conversion", f)
	}
	return dt, err
}

// +func (d *DataFrame) ToArrow() *array.TableReader {
// 	+	arrowFields := make([]arrow.Field, len(d.Schema))
// 	+	for i, cs := range d.Schema {
// 	+		arrowFields[i] = arrow.Field{Name: cs.GetName(), Type: cs.ArrowType()}
// 	+	}
// 	+	schema := arrow.NewSchema(arrowFields, nil)
// 	+
// 	+	pool := memory.NewGoAllocator()
// 	+
// 	+	rb := array.NewRecordBuilder(pool, schema)
// 	+	defer rb.Release()
// 	+
// 	+	records := make([]array.Record, len(d.Records))
// 	+	for rowIdx, row := range d.Records {
// 	+		for fieldIdx, field := range row {
// 	+			switch arrowFields[fieldIdx].Type.(type) {
// 	+			case *arrow.StringType:
// 	+				rb.Field(fieldIdx).(*array.StringBuilder).Append(*(field.(*string)))
// 	+				//rb.Field(fieldIdx).(*array.StringBuilder).AppendValues([]string{*(field.(*string))}, []bool{})
// 	+			case *arrow.Float64Type:
// 	+				rb.Field(fieldIdx).(*array.Float64Builder).Append(*(field.(*float64)))
// 	+				//rb.Field(fieldIdx).(*array.Float64Builder).AppendValues([]float64{*(field.(*float64))}, []bool{})
// 	+			default:
// 	+				fmt.Println("unmatched")
// 	+			}
// 	+		}
// 	+		rec := rb.NewRecord()
// 	+		defer rec.Release()
// 	+		records[rowIdx] = rec
// 	+	}
// 	+	table := array.NewTableFromRecords(schema, records)
// 	+	defer table.Release()
// 	+	tableReader := array.NewTableReader(table, 3)
// 	+	//tableReader.Retain()
// 	+
// 	+	return tableReader
// 	+
// 	+}
