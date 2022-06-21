package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	dsi "github.com/fsrv-xyz/openzfs_exporter/internal/dataset"
	"github.com/prometheus/client_golang/prometheus"
)

// refreshWorker queries the given zpool for datasets and sets the dataset parameters to metrics
func (app *application) refreshWorker(ctx context.Context, done chan<- interface{}, pool string) {
	var sleepCounter int
	for {
		select {
		case <-ctx.Done():
			// write to `done` interface to declare worker as finished
			done <- nil
			return
		default:
			if sleepCounter/int(app.interval.Seconds()) == 0 {
				break
			}
			startTime := time.Now() // begin time measurement

			datasets, err := dsi.DetectDatasets(pool)
			if err != nil {
				log.Println(err)
				break
			}

			wg := sync.WaitGroup{}
			for _, dataset := range datasets {
				wg.Add(1)
				go assignParametersToMetric(dataset, pool, &wg)
			}
			wg.Wait()

			// add spent query time to corresponding metric
			metricExporterQueryDuration.With(prometheus.Labels{MetricLabelPool: pool}).Observe(time.Since(startTime).Seconds())
			sleepCounter = 0
		}
		sleepCounter++
		time.Sleep(time.Second)
	}
}

// assignParametersToMetric - assign read parameter values to `metricZfsParameter` metric
func assignParametersToMetric(dataset *dsi.Dataset, pool string, wg *sync.WaitGroup) {
	defer wg.Done()
	for key, value := range dataset.Parameter {
		valueParsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println(err)
		}
		metricZfsParameter.With(
			prometheus.Labels{
				MetricLabelDataset:   dataset.Name,
				MetricLabelPool:      pool,
				MetricLabelParameter: key,
			},
		).Set(valueParsed)
	}
}
