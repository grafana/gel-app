package main

import (
	"github.com/grafana/gel-app/pkg/gelpoc"
	"github.com/grafana/grafana-plugin-model/go/datasource"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

// GELPlugin stores reference to plugin and logger
type GELPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

// Query Primary method called by grafana-server
func (gp *GELPlugin) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest, api datasource.GrafanaAPI) (*datasource.DatasourceResponse, error) {

	gService := gelpoc.Service{
		DatasourceAPI: api,
	}

	// Build Pipeline from Request

	gReq := gelpoc.GelAppReq{
		DataSourceReq: tsdbReq,
	}

	frames, err := gService.Pipeline(ctx, gReq)

	//gp.logger.Debug("resp", spew.Sdump(frames))

	if err != nil {
		gp.logger.Error("Failed to call api.QueryDatasource", "err", err)
	}

	//pbFrames := make([]*datasource.Frames, len(frames))

	pbFrames := &datasource.Frames{
		Frames: make([]*datasource.Frame, len(frames)),
	}

	for i, frame := range frames {
		pbFrames.Frames[i], err = frame.ToPBFrame()
		if err != nil {
			return nil, err
		}
	}

	resp := &datasource.DatasourceResponse{
		Results: []*datasource.QueryResult{
			&datasource.QueryResult{
				Frames: pbFrames,
			},
		},
	}

	//plugin.logger.Debug("Query", "datasource", tsdbReq.Datasource.Name, "TimeRange", tsdbReq.TimeRange)

	return resp, nil
}
