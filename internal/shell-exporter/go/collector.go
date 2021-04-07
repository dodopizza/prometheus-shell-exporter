package shellexporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Collector type for prometheus.Collector interface implementation
type Collector struct {
	metrics     *metrics
	scripts     []string
	getDataFunc func(string) ([]shellMetric, error)
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
	for lk := range sm.Labels {
		labels = append(labels, lk)
	}
	return
}

// NewCollector is Collector constructor
func NewCollector(scripts []string, getDataFunc func(string) ([]shellMetric, error)) *Collector {
	return &Collector{
		metrics: &metrics{
			totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "scrapes_total",
				Help: "Count of total scrapes",
			}),
			failedScrapes: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "scrapes_failed_total",
				Help: "Count of total failed scrapes",
			}),
			gaugeMetric: map[string]*prometheus.GaugeVec{},
		},
		scripts:     scripts,
		getDataFunc: getDataFunc,
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

	for _, script := range c.scripts {
		scriptName := sanitizePromLabelName(GetFileName(script))
		metrics, err := c.getDataFunc(script)

		if len(metrics) <= 0 {
			continue
		}

		if err != nil {
			c.metrics.failedScrapes.Inc()
			ch <- c.metrics.failedScrapes
		}

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
