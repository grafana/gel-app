package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/grafana/grafana-plugin-sdk-go/transform"
	"github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

var pluginLogger = hclog.New(&hclog.LoggerOptions{
	Name:  "gel-app",
	Level: hclog.LevelFromString("DEBUG"),
})

func healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func registerPProfHandlers(r *http.ServeMux) {
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// GELPlugin stores reference to plugin and logger
type GELPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", healthcheckHandler)

	go func() {
		if err := http.ListenAndServe(":6060", m); err != nil {
			log.Fatal(err)
		}
	}()

	err := transform.Serve("gel-app", &GELPlugin{
		logger: pluginLogger,
	})
	if err != nil {
		pluginLogger.Error(err.Error())
	}
}
