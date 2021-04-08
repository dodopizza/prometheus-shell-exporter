package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	exporter "github.com/dodopizza/prometheus-shell-exporter/internal/shell-exporter/go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var appVersion = "0.0.000" // go build -ldflags "-X main.appVersion=1.2.345"

// Process command-line flag parsing.
func processAppFlags() (scriptsDir string, port int) {
	var showHelpInfo bool
	var showAppVersion bool

	flag.StringVar(&scriptsDir, "f", "./metrics", "scripts dir")
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

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func main() {
	log.Debug().Msg("start exporter")
	scriptsDir, port := processAppFlags()

	err := exporter.ServeMetrics(scriptsDir, "/metrics", port)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
