package main

import (
	"encoding/json"

	"github.com/grafana/gel-app/pkg/gelpoc"
	"github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/grafana/grafana-plugin-sdk-go/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/transform"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GELPlugin stores reference to plugin and logger
type GELPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

func (gp *GELPlugin) Transform(ctx context.Context, tr datasource.TimeRange, ds datasource.DataSourceInfo, queries []transform.Query, api transform.GrafanaAPIHandler) ([]transform.QueryResult, error) {
	svc := gelpoc.Service{
		GrafanaAPI: api,
	}

	// Build Pipeline from Request
	pipeline, err := svc.BuildPipeline(tr, queries)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Execute Pipeline
	frames, err := svc.ExecutePipeline(ctx, pipeline)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	hidden, err := hiddenRefIDs(queries)
	if err != nil {
		return nil, status.Error((codes.Internal), err.Error())
	}

	if len(hidden) != 0 {
		filteredFrames := make([]*dataframe.Frame, 0, len(frames)-len(hidden))
		for _, frame := range frames {
			if _, ok := hidden[frame.RefID]; !ok {
				filteredFrames = append(filteredFrames, frame)
			}
		}
		frames = filteredFrames
	}

	res := []transform.QueryResult{
		{
			DataFrames: frames,
		},
	}

	return res, nil
}

func hiddenRefIDs(queries []transform.Query) (map[string]struct{}, error) {
	hidden := make(map[string]struct{})

	for _, query := range queries {
		hide := struct {
			Hide bool `json:"hide"`
		}{}

		if err := json.Unmarshal(query.ModelJSON, &hide); err != nil {
			return nil, err
		}

		if hide.Hide {
			hidden[query.RefID] = struct{}{}
		}
	}
	return hidden, nil
}
