package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"net/http"
)

func main() {
	var scriptsDir string
	flag.StringVar(&scriptsDir, "f", "/workspaces/prometheus-shell-exporter/metrics", "scripts dir")
	flag.Parse()

	exp := NewExporter(scriptsDir)

	if err := exp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Stack().
			Err(err).
			Msg("Server stopped unexpectedly")
	}
}
