package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	plugin "github.com/hashicorp/go-plugin"
)

// GELPlugin stores reference to plugin and logger
type GELPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	logger log.Logger
}

func main() {
	err := backend.Serve(backend.ServeOpts{
		TransformDataHandler: &GELPlugin{
			logger: backend.Logger,
		},
	})
	if err != nil {
		backend.Logger.Error(err.Error())
	}
}
