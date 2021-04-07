package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	scriptExecutor "github.com/dodopizza/prometheus-shell-exporter/pkg/script-executor/go"
)

type Exporter struct {
	*http.Server
}

func NewExporter(scriptsDir string) *Exporter {

	scripts, err := WalkMatch(scriptsDir, "*")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if len(scripts) <= 0 {
		log.Fatal().Msg("No scripts to serve")
	}

	mux := &http.ServeMux{}

	collector := NewCollector(
		scripts,
		getDataFromShellExecutor,
	)

	registry := prometheus.NewRegistry()

	registry.MustRegister(collector)

	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 9090),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	return &Exporter{server}
}

func getDataFromShellExecutor(script string) (metricsData []shellMetric, err error) {
	exec := scriptExecutor.NewScriptExecutor(scriptExecutor.ShellTypeAutodetect)
	stdOut, _, err := exec.Execute(script)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	decoder := json.NewDecoder(strings.NewReader(stdOut))
	err = decoder.Decode(&metricsData)

	return
}
