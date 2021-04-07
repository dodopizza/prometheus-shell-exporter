package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"

	"net/http"
)

var appVersion = "0.0.000" // go build -ldflags "-X main.appConfigVersion=1.2.345"

func processAppArguments() (scriptsDir string, port int) {
	var showHelpInfo bool
	var showAppVersion bool

	flag.StringVar(&scriptsDir, "f", "metrics", "scripts dir")
	flag.IntVar(&port, "port", 9360, "exporter port")
	flag.BoolVar(&showHelpInfo, "help", false, "help info")
	flag.BoolVar(&showAppVersion, "version", false, "app version info")
	flag.Parse()

	if showHelpInfo {
		fmt.Printf("Usage of %s:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	if showAppVersion {
		fmt.Println(appVersion)
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
