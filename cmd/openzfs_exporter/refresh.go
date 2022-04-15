package main

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dsi "golang.fsrv.services/openzfs_exporter/internal/dataset"
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

			wg := sync.WaitGroup{}
			for _, dataset := range dsi.DetectDatasets(pool) { // loop through datasets
				// apply query filter
				if app.poolFilter.MatchString(dataset.Name) == app.reverseFilter {
					continue
				}
				// start one assignment process per dataset
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
	dataset.ParseParameters() // read parameter values
	for key, value := range dataset.Parameter {
		// assign parameters and values to corresponding metrics
		metricZfsParameter.With(
			prometheus.Labels{
				MetricLabelDataset:   dataset.Name,
				MetricLabelPool:      pool,
				MetricLabelParameter: key,
			},
		).Set(float64(value))
	}
}
