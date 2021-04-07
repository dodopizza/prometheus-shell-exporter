package main

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
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

	pe.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "HelloKitty",
			ConstLabels: map[string]string{},
		},
		func() float64 {
			return 1
		},
	)

	err = pe.Serve()
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(0)
}
