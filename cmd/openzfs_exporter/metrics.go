package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricLabelDataset   = "dataset"
	MetricLabelPool      = "pool"
	MetricLabelParameter = "parameter"
)

var (
	metricZfsParameter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openzfs",
		Subsystem: "zfs",
		Name:      "parameter",
		Help:      "sysctl openzfs dataset parameters",
	}, []string{
		MetricLabelPool,
		MetricLabelDataset,
		MetricLabelParameter,
	})

	metricExporterQueryDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "openzfs",
		Subsystem: "exporter",
		Name:      "query_duration_seconds",
		Buckets:   []float64{0.05, 0.07, 0.09, 0.1, 0.12, 0.14, 0.17, 0.2, 0.25, 0.3, 0.5, 0.7},
	}, []string{
		MetricLabelPool,
	})
)

func init() {
	prometheus.MustRegister(
		metricZfsParameter,
		metricExporterQueryDuration,
	)
}
