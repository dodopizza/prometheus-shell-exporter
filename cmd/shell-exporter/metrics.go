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
	shellMetrics  []*shellMetrics
}

type shellMetric struct {
	Value  int               `json:"value"`
	Labels map[string]string `json:"labels"`
}

type shellMetrics struct {
	fName       string
	metrics     []shellMetric
	gaugeMetric *prometheus.GaugeVec
}

func (sm *shellMetrics) updateData() {
	sm.readFromFile(sm.fName)
}

func (sm *shellMetrics) getName() string {
	return sanitizePromLabelName(GetFileName(sm.fName))
}

func (sm *shellMetrics) getLabels() (labels []string) {
	for lk, _ := range sm.metrics[0].Labels {
		labels = append(labels, lk)
	}
	return
}

func newMetrics(namespace string, scripts []string) (result *metrics) {

	result = &metrics{}

	for _, script := range scripts {
		metric := &shellMetrics{fName: script}
		metric.updateData()
		metric.gaugeMetric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      metric.getName(),
			},
			metric.getLabels(),
		)
		result.shellMetrics = append(result.shellMetrics, metric)
	}

	result.totalScrapes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "scrapes_total",
		Help:      "Count of total scrapes",
	})

	result.failedScrapes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "failed_scrapes_total",
		Help:      "Count of failed scrapes",
	})

	return
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
	err = decoder.Decode(&pm.metrics)
	if err != nil {
		return
	}
	return
}
