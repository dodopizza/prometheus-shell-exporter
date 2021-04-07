package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// Collector type for prometheus.Collector interface implementation
type Collector struct {
	metrics *metrics
}

// NewCollector is Collector constructor
func NewCollector() *Collector {

	scripts, err := WalkMatch("/workspaces/prometheus-shell-exporter/metrics_examples", "*.json")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return &Collector{
		metrics: newMetrics("shell_exporter", scripts),
	}
}

// Describe for prometheus.Collector interface implementation
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect for prometheus.Collector interface implementation
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.metrics.totalScrapes.Inc()
	ch <- c.metrics.totalScrapes
	c.metrics.failedScrapes.Inc()
	ch <- c.metrics.failedScrapes

	for _, mv := range c.metrics.shellMetrics {
		mv.updateData()

		for _, mvv := range mv.metrics {
			mv.gaugeMetric.With(mvv.Labels).Set(float64(mvv.Value))
		}

		mv.gaugeMetric.Collect(ch)
	}
}
