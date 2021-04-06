package main

import (
	"log"
	"os"
)

func main() {

	scripts, err := WalkMatch("/workspaces/prometheus-shell-exporter/examples", "*.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	pe := NewPromExporter()

	for _, script := range scripts {
		p := PromMetrics{}
		p.ReadFromFile(script)
		pe.NewGaugeVecFromPromMetrics(sanitizePromLabelName(GetFileName(script)), p)
	}

	err = pe.Serve()
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(0)
}
