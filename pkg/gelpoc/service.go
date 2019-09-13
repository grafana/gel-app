package gelpoc

import (
	"context"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

// // GelAppReq is current request structure from the frontend for GEL.
// type GelAppReq struct {
// 	Options ReqOptions `json:"options"`
// }

// // ReqOptions is the payload of a GelAppReq.
// type ReqOptions struct {
// 	RequestID   string `json:"requestId"`
// 	Timezone    string `json:"timezone"`
// 	PanelID     int    `json:"panelId"`
// 	DashboardID int    `json:"dashboardId"`
// 	Range       struct {
// 		From time.Time `json:"from"`
// 		To   time.Time `json:"to"`
// 		Raw  struct {
// 			From string `json:"from"`
// 			To   string `json:"to"`
// 		} `json:"raw"`
// 	} `json:"range"`
// 	Interval      string `json:"interval"`
// 	IntervalMs    int    `json:"intervalMs"`
// 	MaxDataPoints int    `json:"maxDataPoints"`
// 	ScopedVars    struct {
// 		Interval struct {
// 			Text  string `json:"text"`
// 			Value string `json:"value"`
// 		} `json:"__interval"`
// 		IntervalMs struct {
// 			Text  string `json:"text"`
// 			Value int    `json:"value"`
// 		} `json:"__interval_ms"`
// 	} `json:"scopedVars"`
// 	RangeRaw struct {
// 		From string `json:"from"`
// 		To   string `json:"to"`
// 	} `json:"rangeRaw"`
// 	Targets []*simplejson.Json `json:"targets"`
// }

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
	//DatasourceCache datasources.CacheService `inject:""`
	// log             log.Logger
}

// // Init initializes the GelService.
// func (s *GelService) Init() error {
// 	s.log = log.New("gel")

// 	s.RouterRegister.Group("/api/gel", func(r routing.RouteRegister) {
// 		r.Post("/expr", binding.Bind(GelAppReq{}), api.Wrap(s.Pipeline))
// 	})

// 	return nil
// }

// Pipeline builds and executes a GEL data pipeline and returns all the results
// as a slice of *data.Frame.
func (s *Service) Pipeline(ctx context.Context, req GelAppReq) ([]*data.Frame, error) {
	timeRange := &datasource.TimeRange{
		FromRaw: req.DataSourceReq.TimeRange.GetFromRaw(),
		ToRaw:   req.DataSourceReq.TimeRange.GetToRaw(),
	}

	pipeline, err := buildPipeline(
		//req.Options.Targets,
		req.DataSourceReq.Queries,
		timeRange,
		//tsdb.NewTimeRange(req.Options.Range.Raw.From, req.Options.Range.Raw.To),
		s.DatasourceAPI,
	)
	if err != nil {
		return nil, err
	}

	vars, err := pipeline.Execute(ctx)
	if err != nil {
		return nil, err
	}

	//frames := extractDataFrames(vars, req.HiddenTargets())
	frames := extractDataFrames(vars, make(map[string]struct{}))

	return frames, nil
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
