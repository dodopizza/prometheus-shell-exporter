package main

import (
	"github.com/rs/zerolog/log"

	"net/http"
)

func main() {
	exp := NewExporter()

	if err := exp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Stack().
			Err(err).
			Msg("Server stopped unexpectedly")
	}
}
