package gelpoc

import (
	"context"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/gel-app/pkg/mathexp"
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

	// Service currently doesn't care about order of dataframes in the return.
	trans := cmp.Transformer("Sort", func(in []*dataframe.Frame) []*dataframe.Frame {
		out := append([]*dataframe.Frame(nil), in...) // Copy input to avoid mutating it
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].RefID > out[j].RefID
		})
		return out
	})
	if diff := cmp.Diff(expect, res, trans); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}
}

type mockGrafanaAPI struct {
	QueryDatasourceFn func() ([]datasource.DatasourceQueryResult, error)
}

func newMockGrafanaAPI(df ...*dataframe.Frame) *mockGrafanaAPI {
	return &mockGrafanaAPI{
		QueryDatasourceFn: func() (res []datasource.DatasourceQueryResult, err error) {
			series := make([]mathexp.Series, 0, len(df))
			for _, frame := range df {
				s, err := mathexp.SeriesFromFrame(frame)
				if err != nil {
					return res, err
				}
				series = append(series, s)
			}

			frames := make([]*dataframe.Frame, len(series))
			for idx, s := range series {
				frames[idx] = s.AsDataFrame()
			}
			return []datasource.DatasourceQueryResult{
				datasource.DatasourceQueryResult{
					DataFrames: frames,
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
