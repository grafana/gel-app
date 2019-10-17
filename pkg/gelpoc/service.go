package gelpoc

import (
	"context"

	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-sdk-go"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

// Service is service representation for GEL.
type Service struct {
	GrafanaAPI grafana.GrafanaAPI
}

// BuildPipeline builds a pipeline from a request.
func (s *Service) BuildPipeline(tr grafana.TimeRange, queries []grafana.Query) (DataPipeline, error) {
	return buildPipeline(queries, tr, s.GrafanaAPI)
}

// ExecutePipeline executes a GEL data pipeline and returns all the results
// as a slice of *dataframe.Frame. Queries that are marked has hidden should be executed
// but should not returned (TODO: currently hidden is ignored).
func (s *Service) ExecutePipeline(ctx context.Context, pipeline DataPipeline) ([]*dataframe.Frame, error) {
	vars, err := pipeline.execute(ctx)
	if err != nil {
		return nil, err
	}
	return extractDataFrames(vars), nil
}

func extractDataFrames(vars mathexp.Vars) []*dataframe.Frame {
	res := []*dataframe.Frame{}
	for refID, results := range vars {
		for _, val := range results.Values {
			df := val.AsDataFrame()
			df.RefID = refID
			res = append(res, df)
		}
	}
	return res
}
