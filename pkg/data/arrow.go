package data

import (
	"fmt"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/mattetti/filebuffer"
)

// ToArrow converts the Frame to an arrow table and returns
// a byte representation of that table.
func (f *Frame) ToArrow() ([]byte, error) {
	// Create arrow schema with metadata
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

	// Build the arrow columns
	pool := memory.NewGoAllocator()
	columns := make([]array.Column, len(f.Fields))
	for fieldIdx, field := range f.Fields {
		// build each column depending on the type
		switch field.Type {
		case TypeNumber:
			builder := array.NewFloat64Builder(pool)
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
			chunked.Release()
		case TypeTime:
			builder := array.NewTimestampBuilder(pool, &arrow.TimestampType{
				Unit: arrow.Nanosecond,
			})
			for _, v := range *field.Vector.(*TimeVector) {
				if v == nil {
					builder.AppendNull()
					continue
				}
				builder.Append(arrow.Timestamp(v.UnixNano()))
			}
			chunked := array.NewChunked(arrowFields[fieldIdx].Type, []array.Interface{builder.NewArray()})

			columns[fieldIdx] = *array.NewColumn(arrowFields[fieldIdx], chunked)
			builder.Release()
			chunked.Release()
		default:
			return nil, fmt.Errorf("unsupported field type %s for arrow converstion", field.Type)
		}
	}

	// Create a table from the schema and columns
	table := array.NewTable(schema, columns, -1)
	defer table.Release()
	tableReader := array.NewTableReader(table, -1)
	defer tableReader.Release()

	// arrow tables with the go API are written to files, so we create a fake
	// file buffer that the FileWriter can write to.
	// In the future and with streaming, I think will likely be using the arrow
	// message type some how.
	fb := filebuffer.New(nil)

	fakeFWriter, err := ipc.NewFileWriter(fb, ipc.WithSchema(tableReader.Schema()))
	if err != nil {
		return nil, err
	}
	defer fakeFWriter.Close()

	for tableReader.Next() {
		rec := tableReader.Record()
		err = fakeFWriter.Write(rec)
		rec.Release()
		if err != nil {
			return nil, err
		}

	}

	err = fakeFWriter.Close()
	if err != nil {
		return nil, err
	}

	return fb.Buff.Bytes(), nil
}

// ArrowType returns the correspeonding primitive type to the
// fields Vector's primitives.
func (f FieldType) ArrowType() (dt arrow.DataType, err error) {
	switch f {
	case TypeString:
		dt = &arrow.StringType{}
	case TypeNumber:
		dt = &arrow.Float64Type{}
	case TypeTime:
		dt = &arrow.TimestampType{}
	default:
		return dt, fmt.Errorf("unsupported type %s for arrow conversion", f)
	}
	return dt, err
}
