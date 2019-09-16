package gelpoc

import (
	"context"
	"encoding/json"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

type GelAppReq struct {
	DataSourceReq *datasource.DatasourceRequest
}

// HiddenRefIDs returns map if refId strings for targets
// that have hide: true in the request.
func (gr *GelAppReq) HiddenRefIDs() (map[string]struct{}, error) {
	hidden := make(map[string]struct{})
	for _, query := range gr.DataSourceReq.GetQueries() {
		refID := query.GetRefId()
		hide := struct {
			Hide bool `json:"hide"`
		}{
			false,
		}
		err := json.Unmarshal([]byte(query.ModelJson), &hide)
		if err != nil {
			return nil, err
		}
		if hide.Hide {
			hidden[refID] = struct{}{}
		}
	}
	return hidden, nil
}

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
	frames := extractDataFrames(vars)

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

func extractDataFrames(vars mathexp.Vars) []*data.Frame {
	res := []*data.Frame{}
	for refID, results := range vars {
		for _, val := range results.Values {
			df := val.AsDataFrame()
			df.RefID = refID
			res = append(res, df)
		}
	}
	return res
}
