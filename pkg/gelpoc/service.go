package gelpoc

import (
	"context"

	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/hashicorp/go-hclog"
)

var pluginLogger hclog.Logger

// SetPluginLogger sets the package level logger for the plugin
func SetPluginLogger(l hclog.Logger) {
	pluginLogger = l
}

// Service is service representation for GEL.
type Service struct {
	CallBack backend.TransformCallBackHandler
}

// BuildPipeline builds a pipeline from a request.
func (s *Service) BuildPipeline(queries []backend.DataQuery) (DataPipeline, error) {
	return buildPipeline(queries, s.CallBack)
}

// ExecutePipeline executes a GEL data pipeline and returns all the results
// as a slice of *dataframe.Frame.
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
