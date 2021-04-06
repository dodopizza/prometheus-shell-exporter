package main

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	p := PromMetrics{
		metrics: []PromMetric{
			{
				value: 0,
				labels: map[string]string{
					"hello": "hi",
				},
			},
		},
	}
	// p.ReadFromFile("/workspaces/powershell_exporter/examples/pse_tcp_connection_metrics.example.json")
	p.SaveConfig("/tmp/123.json")

	pe := NewPromExporter()

	pg := pe.NewGaugeVec(prometheus.GaugeOpts{
		Name: "Test",
		Help: "TestHelp",
	},
		[]string{
			"one",
			"two",
		})
	// pg.Set(123)

	pg.With(prometheus.Labels{"one": "3", "two": "4"}).Set(444)
	// pg.WithLabelValues("hello", "kitty").Set(123)

	pe.Serve()

	os.Exit(0)
}
