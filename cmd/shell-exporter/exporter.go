package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Exporter struct {
	*http.Server
}

func NewExporter() *Exporter {

	scripts, err := WalkMatch("/workspaces/prometheus-shell-exporter/metrics_examples", "*.json")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if len(scripts) <= 0 {
		log.Fatal().Msg("No scripts to serve")
	}

	mux := &http.ServeMux{}

	collector := NewCollector(scripts)
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
