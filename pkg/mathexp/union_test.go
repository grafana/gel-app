package mathexp

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/stretchr/testify/assert"
)

func Test_union(t *testing.T) {
	var tests = []struct {
		name      string
		aResults  Results
		bResults  Results
		unionsAre assert.ComparisonAssertionFunc
		unions    []*Union
	}{
		{
			name: "equal tags single union",
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: dataframe.Labels{"id": "1"},
					A:      makeSeries("a", dataframe.Labels{"id": "1"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1"}),
				},
			},
		},
		{
			name: "equal tags keys with no matching values will result in no unions",
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "2"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions:    []*Union{},
		},
		{
			name:      "empty results will result in no unions",
			aResults:  Results{},
			bResults:  Results{},
			unionsAre: assert.EqualValues,
			unions:    []*Union{},
		},
		{
			name: "incompatible tags of different length with will result in no unions",
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"ID": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions:    []*Union{},
		},
		{
			name: "A is subset of B results in single union with Labels of B",
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"}, // Union gets the labels that is not the subset
					A:      makeSeries("a", dataframe.Labels{"id": "1"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
			},
		},
		{
			name: "B is subset of A results in single union with Labels of A",
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"}, // Union gets the labels that is not the subset
					A:      makeSeries("a", dataframe.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1"}),
				},
			},
		},
		{
			name: "single valued A is subset of many valued B, results in many union with Labels of B",
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
					makeSeries("b", dataframe.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("a", dataframe.Labels{"id": "1"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
				{
					Labels: dataframe.Labels{"id": "1", "fish": "red snapper"},
					A:      makeSeries("a", dataframe.Labels{"id": "1"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
		},
		{
			name: "A with different tags keys lengths to B makes 3 unions (with two unions have matching tags)",
			// Is this the behavior we want? A result within the results will no longer
			// be uniquely identifiable.
			aResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1"}),
					makeSeries("aa", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
					makeSeries("bb", dataframe.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("a", dataframe.Labels{"id": "1"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
				{
					Labels: dataframe.Labels{"id": "1", "fish": "red snapper"},
					A:      makeSeries("a", dataframe.Labels{"id": "1"}),
					B:      makeSeries("bb", dataframe.Labels{"id": "1", "fish": "red snapper"}),
				},
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("aa", dataframe.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
			},
		},
		{
			name: "B with different tags keys lengths to A makes 3 unions (with two unions have matching tags)",
			// Is this the behavior we want? A result within the results will no longer
			// be uniquely identifiable.
			aResults: Results{
				Values: Values{
					makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
					makeSeries("bb", dataframe.Labels{"id": "1", "fish": "red snapper"}),
				},
			},
			bResults: Results{
				Values: Values{
					makeSeries("a", dataframe.Labels{"id": "1"}),
					makeSeries("aa", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
			},
			unionsAre: assert.EqualValues,
			unions: []*Union{
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("a", dataframe.Labels{"id": "1"}),
				},
				{
					Labels: dataframe.Labels{"id": "1", "fish": "herring"},
					A:      makeSeries("b", dataframe.Labels{"id": "1", "fish": "herring"}),
					B:      makeSeries("aa", dataframe.Labels{"id": "1", "fish": "herring"}),
				},
				{
					Labels: dataframe.Labels{"id": "1", "fish": "red snapper"},
					A:      makeSeries("bb", dataframe.Labels{"id": "1", "fish": "red snapper"}),
					B:      makeSeries("a", dataframe.Labels{"id": "1"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unions := union(tt.aResults, tt.bResults)
			tt.unionsAre(t, tt.unions, unions)
		})
	}
}
