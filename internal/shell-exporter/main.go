package main

import (
	"log"

	exporter "github.com/dodopizza/prometheus-shell-exporter/internal/shell-exporter/go"
)

func main() {
	err := exporter.Run("examples/", 9999)
	if err != nil {
		log.Fatal(err.Error())
	}
}
