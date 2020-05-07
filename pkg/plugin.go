package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	plugin "github.com/hashicorp/go-plugin"
)

// GELPlugin stores reference to plugin
type GELPlugin struct {
	plugin.NetRPCUnsupportedPlugin
}

func main() {
	err := backend.Serve(backend.ServeOpts{
		TransformDataHandler: &GELPlugin{},
	})
	if err != nil {
		backend.Logger.Error(err.Error())
	}
}
