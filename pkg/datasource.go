package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/gel-app/pkg/data"
	"github.com/grafana/gel-app/pkg/gelpoc"
	"github.com/grafana/grafana-plugin-model/go/datasource"
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

// Query Primary method called by grafana-server
func (gp *GELPlugin) Query(ctx context.Context, tsdbReq *datasource.DatasourceRequest, api datasource.GrafanaAPI) (*datasource.DatasourceResponse, error) {
	gService := gelpoc.Service{
		DatasourceAPI: api,
	}

	gReq := gelpoc.GelAppReq{
		DataSourceReq: tsdbReq,
	}

	// Build Pipeline from Request
	pipeline, err := gService.BuildPipeline(gReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Execute Pipeline
	frames, err := gService.ExecutePipeline(ctx, pipeline)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	hidden, err := gReq.HiddenRefIDs()
	if err != nil {
		return nil, status.Error((codes.Internal), err.Error())
	}

	if len(hidden) != 0 {
		filteredFrames := make([]*data.Frame, 0, len(frames)-len(hidden))
		for _, frame := range frames {
			if _, ok := hidden[frame.RefID]; !ok {
				filteredFrames = append(filteredFrames, frame)
			}
		}
		frames = filteredFrames
	}

	// Convert the []*data.Frames to protobuf Frames appropriate for the plugin-model
	pbFrames := &datasource.Frames{
		Frames: make([]*datasource.Frame, len(frames)),
	}
	for i, frame := range frames {
		pbFrames.Frames[i], err = frame.ToPBFrame()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	byteFrames := make([][]byte, len(frames))

	for i, frame := range frames {
		//gp.logger.Debug("frame", spew.Sdump(frame))
		b, err := frame.ToArrow()
		if err != nil {
			return nil, err
		}
		gp.logger.Debug("ref", frame.RefID, "len", fmt.Sprintf("%v", len(b)))
		byteFrames[i] = b
	}

	jBFrames, err := json.Marshal(byteFrames) // json.Marshal b64 encodes []byte
	if err != nil {
		return nil, err
	}

	gp.logger.Debug("json", string(jBFrames))

	return &datasource.DatasourceResponse{
		Results: []*datasource.QueryResult{
			&datasource.QueryResult{
				//Frames:   pbFrames,
				MetaJson: string(jBFrames),
			},
		},
	}, nil
}
