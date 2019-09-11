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

	// Process App Request once we can get it from the Plugin-API
	needGelAppReqPlz := gelpoc.GelAppReq{}
	_ = needGelAppReqPlz

	gService := gelpoc.Service{}
	_ = gService

	// Build Pipeline from Request

	// Execute Pipeline
	//	Executing the pipeline will require the bi-directional calls
	_, err := api.QueryDatasource(ctx, &datasource.QueryDatasourceRequest{
		DatasourceId: 1,
		Queries:      tsdbReq.Queries,
		TimeRange:    tsdbReq.TimeRange,
	})

	if err != nil {
		gp.logger.Error("Failed to call api.QueryDatasource", "err", err)
	}

	//plugin.logger.Debug("Query", "datasource", tsdbReq.Datasource.Name, "TimeRange", tsdbReq.TimeRange)

	return nil, nil
}
