package main

import (
	"log"
	"os"
)

func main() {

	scripts, err := WalkMatch("/workspaces/prometheus-shell-exporter/metrics_examples", "*.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	pe := NewPromExporter()
	pm := []*PromMetrics{}

	for _, script := range scripts {
		p := PromMetrics{
			Name: sanitizePromLabelName(GetFileName(script)),
		}
		p.ReadFromFile(script)
		pe.NewGaugeVecFromPromMetrics(p)
		pm = append(pm, &p)
	}

	err = pe.Serve()
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(0)
}
