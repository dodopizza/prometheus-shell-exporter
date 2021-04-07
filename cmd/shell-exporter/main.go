package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"

	"net/http"
)

func processAppArguments() (scriptsDir string, port int) {
	var needHelp bool
	flag.StringVar(&scriptsDir, "f", "../metrics", "scripts dir")
	flag.IntVar(&port, "port", 9360, "exporter port")
	flag.BoolVar(&needHelp, "help", false, "help info")
	flag.Parse()

	if needHelp {
		fmt.Printf("Usage of %s:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	return
}

func main() {

	scriptsDir, port := processAppArguments()

	exp := NewExporter(scriptsDir, port)

	if err := exp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Stack().
			Err(err).
			Msg("Server stopped unexpectedly")
	}
}
