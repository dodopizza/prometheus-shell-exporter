package shell_exporter

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetric struct {
	Value  int               `json:"value"`
	Labels map[string]string `json:"labels"`
}

type PromMetrics struct {
	Metrics []PromMetric
}

func (pm *PromMetrics) ReadFromFile(fname string) (err error) {
	file, err := os.Open(fname)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&pm.Metrics)
	if err != nil {
		return
	}
	return
}

type PromExporter struct {
	registry    *prometheus.Registry
	httpHandler http.Handler
}

func NewPromExporter() *PromExporter {
	var pe = new(PromExporter)
	pe.registry = prometheus.NewRegistry()
	pe.httpHandler = promhttp.HandlerFor(pe.registry, promhttp.HandlerOpts{})
	return pe
}

func (pe *PromExporter) NewGaugeVecFromPromMetrics(metricName string, promMetrics PromMetrics) {
	var metric_labels []string

	for lk, _ := range promMetrics.Metrics[0].Labels {
		metric_labels = append(metric_labels, lk)
	}

	pg := pe.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: "TestHelp",
	},
		metric_labels)

	for mi, _ := range promMetrics.Metrics {
		pg.With(promMetrics.Metrics[mi].Labels).Set(float64(promMetrics.Metrics[mi].Value))
	}
}

func (pe *PromExporter) Serve() (err error) {
	http.Handle("/metrics", pe.httpHandler)
	err = http.ListenAndServe(":4567", nil)
	if err != nil {
		return
	}
	return nil
}

func (pe *PromExporter) NewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) (pg *prometheus.GaugeVec) {
	pg = prometheus.NewGaugeVec(opts, labelNames)
	pe.registry.MustRegister(pg)
	return
}

// func prom() {
// 	promRegistry := prometheus.NewRegistry()
// 	promHandler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})

// 	gaugeOpts := prometheus.GaugeOpts{
// 		Name:        "Test",
// 		Help:        "TestHelp",
// 		ConstLabels: map[string]string{},
// 	}
// 	gc := prometheus.NewGauge(gaugeOpts)
// 	promRegistry.MustRegister(gc)

// 	http.Handle("/metrics", promHandler)
// 	err := http.ListenAndServe(":4567", nil)
// 	if err != nil {
// 		log.Printf("http.ListenAndServer: %v\n", err)
// 	}
// }
