package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

var pluginLogger = hclog.New(&hclog.LoggerOptions{
	Name:  "gel-app",
	Level: hclog.LevelFromString("DEBUG"),
})

// GELPlugin stores reference to plugin and logger
type GELPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

func main() {
	err := backend.Serve(backend.ServeOpts{
		TransformDataHandler: &GELPlugin{
			logger: pluginLogger,
		},
	})
	if err != nil {
		pluginLogger.Error(err.Error())
	}
}
