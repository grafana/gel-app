package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrame_Copy(t *testing.T) {
	tests := []struct {
		name        string
		frame       *Frame
		frameCopyIs assert.ComparisonAssertionFunc
		frameCopy   *Frame
	}{
		{
			name: "copy should not be mutated",
			frame: &Frame{
				Name:   "testFrame",
				RefID:  "A",
				Labels: Labels{"foo": "bar"},
				Fields: Fields{
					&Field{
						Name:   "testField",
						Type:   TypeTime,
						Vector: &TimeVector{unixTimePointer(0, 0), unixTimePointer(1, 0)},
					},
				},
			},
			frameCopyIs: assert.Equal,
			frameCopy: &Frame{
				Name:   "testFrame",
				RefID:  "A",
				Labels: Labels{"foo": "bar"},
				Fields: Fields{
					&Field{
						Name:   "testField",
						Type:   TypeTime,
						Vector: &TimeVector{unixTimePointer(0, 0), unixTimePointer(1, 0)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.frame.Copy()
			tt.frame.Name = "somethingElseFrame"
			tt.frame.RefID = "B"
			tt.frame.Labels = Labels{"someThn": "zelse"}
			tt.frame.Fields[0].Vector.SetValue(0, unixTimePointer(1e5, 0))
			tt.frame.Fields = append(tt.frame.Fields, &Field{})
			tt.frameCopyIs(t, tt.frameCopy, v)
		})
	}
}
