package main

import (
	"context"
	"openzfs_exporter/internal/dataset"
	"time"

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
			datasets := dataset.DetectDatasets(pool)
			st := time.Now()
			for _, ds := range datasets {
				// apply query filter
				if app.poolFilter.MatchString(ds.Name) == app.reverseFilter {
					continue
				}

				// read parameter values
				ds.ParseParameters()
				for key, value := range ds.Parameter {
					// assign parameters and values to corresponding metrics
					zpoolStats.With(
						prometheus.Labels{
							MetricLabelDataset:   ds.Name,
							MetricLabelPool:      pool,
							MetricLabelParameter: key,
						},
					).Set(float64(value))
				}
			}
			// add spent query time to corresponding metric
			queryTime.With(prometheus.Labels{MetricLabelPool: pool}).Add(time.Since(st).Seconds())
			sleepCounter = 0
		}
		sleepCounter++
		time.Sleep(1 * time.Second)
	}
}
