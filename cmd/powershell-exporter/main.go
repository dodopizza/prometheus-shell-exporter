package main

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	p := PromMetrics{}
	p.ReadFromFile("/workspaces/powershell_exporter/examples/pse_tcp_connection_metrics.example.json")

	var metric_labels []string

	for lk, _ := range p.Metrics[0].Labels {
		metric_labels = append(metric_labels, lk)
	}

	pe := NewPromExporter()

	pg := pe.NewGaugeVec(prometheus.GaugeOpts{
		Name: "Test",
		Help: "TestHelp",
	},
		metric_labels)

	for mi, _ := range p.Metrics {
		pg.With(p.Metrics[mi].Labels).Set(float64(p.Metrics[mi].Value))
	}

	pe.Serve()

	os.Exit(0)
}
