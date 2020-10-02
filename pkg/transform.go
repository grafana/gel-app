package main

import (
	"encoding/json"

	"github.com/grafana/gel-app/pkg/gelpoc"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransformData takes Queries which are either GEL nodes (a.k.a expressions/transforms)
// or are datasource requests. The transform.GrafanaAPIHandler allows callbacks
// to grafana to fulfill datasource requests.
func (gp *GELPlugin) TransformData(ctx context.Context, req *backend.QueryDataRequest, callBack backend.TransformDataCallBackHandler) (*backend.QueryDataResponse, error) {
	svc := gelpoc.Service{
		CallBack: callBack,
	}
	// Build the pipeline from the request, checking for ordering issues (e.g. loops)
	// and parsing graph nodes from the queries.
	pipeline, err := svc.BuildPipeline(req.Queries)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Execute the pipeline
	responses, err := svc.ExecutePipeline(ctx, pipeline)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	// Get which queries have the Hide property so they those queries' results
	// can be excluded from the response.
	hidden, err := hiddenRefIDs(req.Queries)
	if err != nil {
		return nil, status.Error((codes.Internal), err.Error())
	}

	if len(hidden) != 0 {
		filteredRes := backend.NewQueryDataResponse()
		for refID, res := range responses.Responses {
			if _, ok := hidden[refID]; !ok {
				filteredRes.Responses[refID] = res
			}
		}
		responses = filteredRes
	}

	return responses, nil

}

func hiddenRefIDs(queries []backend.DataQuery) (map[string]struct{}, error) {
	hidden := make(map[string]struct{})

	for _, query := range queries {
		hide := struct {
			Hide bool `json:"hide"`
		}{}

		if err := json.Unmarshal(query.JSON, &hide); err != nil {
			return nil, err
		}

		if hide.Hide {
			hidden[query.RefID] = struct{}{}
		}
	}
	return hidden, nil
}
