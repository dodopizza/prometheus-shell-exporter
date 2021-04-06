package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetric struct {
	value  int               `json:"value"`
	labels map[string]string `json:"labels"`
}

type PromMetrics struct {
	metrics []PromMetric `json:"metrics"`
}

func (pm *PromMetrics) ReadFromFile(fname string) (err error) {
	file, err := os.Open(fname)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&pm.metrics)
	if err != nil {
		return
	}
	return
}

func (pm *PromMetrics) SaveConfig(fname string) (err error) {
	file, err := json.MarshalIndent(pm, "", " ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fname, file, 0644)
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
