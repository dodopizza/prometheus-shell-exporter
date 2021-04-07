package main

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	totalScrapes  prometheus.Counter
	failedScrapes prometheus.Counter
	shellMetrics  []*prometheus.GaugeVec
}

type shellMetric struct {
	Value  int               `json:"value"`
	Labels map[string]string `json:"labels"`
}

type shellMetrics struct {
	Name    string
	Metrics []shellMetric
}

func getMetrics(fname string) *shellMetrics {
	s := &shellMetrics{
		Name: sanitizePromLabelName(GetFileName(fname)),
	}
	s.readFromFile(fname)

	return s
}

func newMetrics(namespace string, scripts []string) *metrics {
	return &metrics{
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "scrapes_total",
			Help:      "Count of total scrapes",
		}),

		failedScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "failed_scrapes_total",
			Help:      "Count of failed scrapes",
		}),
	}
}

func sanitizePromLabelName(str string) string {
	re := regexp.MustCompile(`[\.\-]`)
	result := re.ReplaceAllString(str, "_")
	re = regexp.MustCompile(`^\d`)
	result = re.ReplaceAllString(result, "_$0")
	return result
}

func (pm *shellMetrics) readFromFile(fname string) (err error) {
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
