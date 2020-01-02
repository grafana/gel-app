package gelpoc

import (
	"context"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {

	dsDF := dataframe.New("test",
		dataframe.NewField("time", nil, []*time.Time{utp(1)}),
		dataframe.NewField("value", nil, []*float64{fp(2)}))

	m := newMockTransformCallBack(dsDF)

	s := Service{m}

	queries := []backend.DataQuery{
		backend.DataQuery{
			RefID: "A",
			JSON:  json.RawMessage(`{ "datasource": "test", "datasourceId": 3, "orgId": 1, "intervalMs": 1000, "maxDataPoints": 1000 }`),
		},
		backend.DataQuery{
			RefID: "B",
			JSON:  json.RawMessage(`{ "datasource": "__expr__", "datasourceId": -100, "type": "math", "expression": "$A * 2" }`),
		},
	}

	pl, err := s.BuildPipeline(queries)
	require.NoError(t, err)

	res, err := s.ExecutePipeline(context.Background(), pl)
	require.NoError(t, err)

	bDF := dataframe.New("",
		dataframe.NewField("Time", nil, []*time.Time{utp(1)}),
		dataframe.NewField("", nil, []*float64{fp(4)}))
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

type mockTransformCallBack struct {
	DataQueryFn func() (*backend.DataQueryResponse, error)
}

func newMockTransformCallBack(df ...*dataframe.Frame) *mockTransformCallBack {
	return &mockTransformCallBack{
		DataQueryFn: func() (res *backend.DataQueryResponse, err error) {
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
			return &backend.DataQueryResponse{
				Frames: frames,
			}, nil

		},
	}
}

func (m *mockTransformCallBack) DataQuery(ctx context.Context, pc backend.PluginConfig, headers map[string]string, queries []backend.DataQuery) (*backend.DataQueryResponse, error) {
	return m.DataQueryFn()
}

func utp(sec int64) *time.Time {
	t := time.Unix(sec, 0)
	return &t
}

func fp(f float64) *float64 {
	return &f
}
