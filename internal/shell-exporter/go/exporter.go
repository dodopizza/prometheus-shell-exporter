package shellexporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	scriptExecutor "github.com/dodopizza/prometheus-shell-exporter/pkg/script-executor/go"
)

func ServeMetrics(scriptsDir string, metricsHTTPEndpoint string, port int) (err error) {
	expHandler, err := NewExporterHandler(scriptsDir)
	if err != nil {
		return
	}

	mux := &http.ServeMux{}
	mux.Handle(metricsHTTPEndpoint, expHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return
	}

	return
}

func NewExporterHandler(scriptsDir string) (handler http.Handler, err error) {
	scripts, err := getMetricsScripts(scriptsDir)
	if err != nil {
		return
	}

	collector := newCollector(
		scripts,
		getDataFromShellExecutor,
	)

	registry := prometheus.NewRegistry()

	registry.MustRegister(collector)

	handler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return
}

func getDataFromShellExecutor(script string) (metricsData []shellMetric, err error) {
	exec := scriptExecutor.NewScriptExecutor(scriptExecutor.ShellTypeAutodetect)
	stdOut, _, err := exec.Execute(script)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(strings.NewReader(stdOut))
	err = decoder.Decode(&metricsData)

	return
}
