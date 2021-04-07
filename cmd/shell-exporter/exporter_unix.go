// +build !windows

package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	exporter "github.com/dodopizza/prometheus-shell-exporter/internal/shell-exporter/go"
	"github.com/rs/zerolog/log"
)

func processAppArguments() (scriptsDir string, port int) {
	var showHelpInfo bool
	var showAppVersion bool

	flag.StringVar(&scriptsDir, "f", "../../metrics", "scripts dir")
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

func run() {

	scriptsDir, port := processAppArguments()

	err := exporter.Run(scriptsDir, port)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
