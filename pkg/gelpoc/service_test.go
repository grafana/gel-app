package gelpoc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/grafana/grafana-plugin-sdk-go/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/transform"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {

	dsDF := dataframe.New("test", nil,
		dataframe.NewField("time", []*time.Time{utp(1)}),
		dataframe.NewField("value", []*float64{fp(2)}))

	m := newMockGrafanaAPI(dsDF)

	s := Service{m}

	tr := datasource.TimeRange{
		From: time.Unix(0, 0),
		To:   time.Unix(1, 0),
	}
	queries := []transform.Query{
		transform.Query{
			RefID:     "A",
			ModelJSON: json.RawMessage(`{ "datasource": "test", "datasourceId": 3, "orgId": 1, "intervalMs": 1000, "maxDataPoints": 1000 }`),
		},
		transform.Query{
			RefID:     "B",
			ModelJSON: json.RawMessage(`{ "datasource": "__expr__", "datasourceId": -100, "type": "math", "expression": "$A * 2" }`),
		},
	}

	pl, err := s.BuildPipeline(tr, queries)
	require.NoError(t, err)

	res, err := s.ExecutePipeline(context.Background(), pl)
	require.NoError(t, err)

	bDF := dataframe.New("", nil,
		dataframe.NewField("Time", []*time.Time{utp(1)}),
		dataframe.NewField("", []*float64{fp(4)}))
	bDF.RefID = "B"

	expect := []*dataframe.Frame{dsDF, bDF}

	if diff := cmp.Diff(expect, res); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}

}

type mockGrafanaAPI struct {
	QueryDatasourceFn func() ([]datasource.DatasourceQueryResult, error)
}

func newMockGrafanaAPI(df ...*dataframe.Frame) *mockGrafanaAPI {
	return &mockGrafanaAPI{
		QueryDatasourceFn: func() ([]datasource.DatasourceQueryResult, error) {
			return []datasource.DatasourceQueryResult{
				datasource.DatasourceQueryResult{
					DataFrames: df,
				},
			}, nil

		},
	}
}

func (m *mockGrafanaAPI) QueryDatasource(ctx context.Context, orgID int64, datasourceID int64, tr datasource.TimeRange, queries []datasource.Query) ([]datasource.DatasourceQueryResult, error) {
	return m.QueryDatasourceFn()
}

func utp(sec int64) *time.Time {
	t := time.Unix(sec, 0)
	return &t
}

func fp(f float64) *float64 {
	return &f
}
