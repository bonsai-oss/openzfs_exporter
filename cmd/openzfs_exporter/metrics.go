package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetricLabelDataset   = "dataset"
	MetricLabelPool      = "pool"
	MetricLabelParameter = "parameter"
)

var (
	zpoolStats = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openzfs",
		Subsystem: "zfs",
		Name:      "parameter",
		Help:      "sysctl openzfs dataset parameters",
	}, []string{
		MetricLabelPool,
		MetricLabelDataset,
		MetricLabelParameter,
	})

	queryTime = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "openzfs",
		Subsystem: "exporter",
		Name:      "query_seconds_total",
		Help:      "time spent to gather parameters",
	}, []string{
		MetricLabelPool,
	})
)
