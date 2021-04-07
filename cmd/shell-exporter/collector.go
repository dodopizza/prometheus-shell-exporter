package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Collector type for prometheus.Collector interface implementation
type Collector struct {
	metrics *metrics
}

// NewCollector is Collector constructor
func NewCollector() *Collector {
	return &Collector{
		metrics: newMetrics("shell_exporter"),
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
}
