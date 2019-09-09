package gelpoc

import (
	"time"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/mathexp"
	"github.com/grafana/grafana/pkg/components/simplejson"
)

// GelAppReq is current request structure from the frontend for GEL.
type GelAppReq struct {
	Options ReqOptions `json:"options"`
}

// ReqOptions is the payload of a GelAppReq.
type ReqOptions struct {
	RequestID   string `json:"requestId"`
	Timezone    string `json:"timezone"`
	PanelID     int    `json:"panelId"`
	DashboardID int    `json:"dashboardId"`
	Range       struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
		Raw  struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"raw"`
	} `json:"range"`
	Interval      string `json:"interval"`
	IntervalMs    int    `json:"intervalMs"`
	MaxDataPoints int    `json:"maxDataPoints"`
	ScopedVars    struct {
		Interval struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		} `json:"__interval"`
		IntervalMs struct {
			Text  string `json:"text"`
			Value int    `json:"value"`
		} `json:"__interval_ms"`
	} `json:"scopedVars"`
	RangeRaw struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"rangeRaw"`
	Targets []*simplejson.Json `json:"targets"`
}

// HiddenTargets returns map if refId strings for targets
// that have hide: true in the request.
func (gr *GelAppReq) HiddenTargets() map[string]struct{} {
	hidden := make(map[string]struct{})
	for _, target := range gr.Options.Targets {
		refID := target.Get("refId").MustString()
		if target.Get("hide").MustBool() {
			hidden[refID] = struct{}{}
		}
	}
	return hidden
}

// Service is service representation for GEL.
type Service struct {
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
// func (s *Service) Pipeline(c *models.ReqContext, req GelAppReq) api.Response {
// 	pipeline, err := buildPipeline(
// 		req.Options.Targets,
// 		tsdb.NewTimeRange(req.Options.Range.Raw.From, req.Options.Range.Raw.To),
// 		s.DatasourceCache,
// 	)
// 	if err != nil {
// 		return api.Error(500, "error building pipeline", err)
// 	}

// 	vars, err := pipeline.Execute(c)
// 	if err != nil {
// 		return api.Error(500, "failed to execute pipeline", err)
// 	}

// 	frames := extractDataFrames(vars, req.HiddenTargets())

// 	return api.JSON(200, frames)
// }

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
