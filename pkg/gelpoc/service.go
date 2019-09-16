package gelpoc

import (
	"context"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

type GelAppReq struct {
	DataSourceReq *datasource.DatasourceRequest
}

// HiddenTargets returns map if refId strings for targets
// that have hide: true in the request.
// func (gr *GelAppReq) HiddenTargets() map[string]struct{} {
// 	hidden := make(map[string]struct{})
// 	for _, target := range gr.Options.Targets {
// 		refID := target.Get("refId").MustString()
// 		if target.Get("hide").MustBool() {
// 			hidden[refID] = struct{}{}
// 		}
// 	}
// 	return hidden
// }

// Service is service representation for GEL.
type Service struct {
	DatasourceAPI datasource.GrafanaAPI
}

// BuildPipeline builds a pipeline from a request.
func (s *Service) BuildPipeline(req GelAppReq) (DataPipeline, error) {
	timeRange := &datasource.TimeRange{
		FromRaw: req.DataSourceReq.TimeRange.GetFromRaw(),
		ToRaw:   req.DataSourceReq.TimeRange.GetToRaw(),
	}
	return buildPipeline(
		req.DataSourceReq.Queries,
		timeRange,
		s.DatasourceAPI,
	)
}

// ExecutePipeline executes a GEL data pipeline and returns all the results
// as a slice of *data.Frame. Queries that are marked has hidden should be executed
// but should not returned (TODO: currently hidden is ignored).
func (s *Service) ExecutePipeline(ctx context.Context, pipeline DataPipeline) ([]*data.Frame, error) {
	vars, err := pipeline.execute(ctx)
	if err != nil {
		return nil, err
	}
	//TODO: hide targets. frames := extractDataFrames(vars, req.HiddenTargets())
	frames := extractDataFrames(vars, make(map[string]struct{}))

	return frames, nil
}

// BuildAndExecutePipeline builds and executes a GEL data pipeline and returns all the results
// as a slice of *data.Frame.
func (s *Service) BuildAndExecutePipeline(ctx context.Context, req GelAppReq) ([]*data.Frame, error) {
	pipeline, err := s.BuildPipeline(req)
	if err != nil {
		return nil, err
	}
	return s.ExecutePipeline(ctx, pipeline)
}

func extractDataFrames(vars mathexp.Vars, hidden map[string]struct{}) []*data.Frame {
	res := []*data.Frame{}
	for refID, results := range vars {
		if _, ok := hidden[refID]; ok {
			continue // do not return hidden results
		}
		for _, val := range results.Values {
			df := val.AsDataFrame()
			df.RefID = refID
			res = append(res, df)
		}
	}
	return res
}
