package main

import (
	"encoding/json"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

// Collector type for prometheus.Collector interface implementation
type Collector struct {
	metrics *metrics
	scripts []string
}

//
type metrics struct {
	totalScrapes  prometheus.Counter
	failedScrapes prometheus.Counter
	gaugeMetric   map[string]*prometheus.GaugeVec
}

type shellMetric struct {
	Value  int               `json:"value"`
	Labels map[string]string `json:"labels"`
}

func (sm *shellMetric) getLabels() (labels []string) {
	for lk, _ := range sm.Labels {
		labels = append(labels, lk)
	}
	return
}

func getData(fname string) (metricsData []shellMetric) {
	file, err := os.Open(fname)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&metricsData)
	if err != nil {
		return
	}
	return
}

// NewCollector is Collector constructor
func NewCollector(scripts []string) *Collector {
	return &Collector{
		metrics: &metrics{
			totalScrapes:  nil,
			failedScrapes: nil,
			gaugeMetric:   map[string]*prometheus.GaugeVec{},
		},
		scripts: scripts,
	}
}

// Describe for prometheus.Collector interface implementation
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect for prometheus.Collector interface implementation
func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	if c.metrics.totalScrapes == nil {
		c.metrics.totalScrapes = prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scrapes_total",
			Help: "Count of total scrapes",
		})
	}
	c.metrics.totalScrapes.Inc()
	ch <- c.metrics.totalScrapes

	if c.metrics.failedScrapes == nil {
		c.metrics.failedScrapes = prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scrapes_failed_total",
			Help: "Count of total failed scrapes",
		})
	}
	c.metrics.failedScrapes.Inc()
	ch <- c.metrics.failedScrapes

	for _, script := range c.scripts {
		scriptName := sanitizePromLabelName(GetFileName(script))
		metrics := getData(script)

		if _, ok := c.metrics.gaugeMetric[scriptName]; !ok {
			m := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{Name: scriptName},
				metrics[0].getLabels(),
			)
			c.metrics.gaugeMetric[scriptName] = m
		}

		for _, mv := range metrics {
			c.metrics.gaugeMetric[scriptName].With(mv.Labels).Set(float64(mv.Value))
		}

		c.metrics.gaugeMetric[scriptName].Collect(ch)
	}
}
