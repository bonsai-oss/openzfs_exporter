package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bonsai-oss/openzfs_exporter/internal/dataset"
)

// refreshWorker queries the given zpool for datasets and sets the dataset parameters to metrics
func (app *application) refreshWorker(ctx context.Context, done chan<- interface{}, pool string) {
	// initialize a new dataset list for determining the difference between the previous and current dataset list
	datasetCache := make(map[string]bool)
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

			datasets, err := dataset.DetectDatasets(pool)
			if err != nil {
				log.Println(err)
				break
			}

			/*
				The following code block checks if all previous detected datasets are still existing.
				If the dataset is not existing anymore, the metrics related to it is deleted.
			*/
			for item := range datasetCache {
				datasetCache[item] = false
			}
			for _, ds := range datasets {
				datasetCache[ds.Name] = true
			}
			for item := range datasetCache {
				if !datasetCache[item] {
					metricZfsParameter.DeletePartialMatch(prometheus.Labels{MetricLabelPool: pool, MetricLabelDataset: item})
					delete(datasetCache, item)
				}
			}

			wg := sync.WaitGroup{}
			for _, ds := range datasets {
				if app.poolFilter.MatchString(ds.Name) == app.reverseFilter {
					continue
				}
				wg.Add(1)
				go assignParametersToMetric(ds, pool, &wg)
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
func assignParametersToMetric(ds *dataset.Dataset, pool string, wg *sync.WaitGroup) {
	defer wg.Done()
	for key, value := range ds.Parameter {
		valueParsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println(err)
		}
		metricZfsParameter.With(
			prometheus.Labels{
				MetricLabelDataset:   ds.Name,
				MetricLabelPool:      pool,
				MetricLabelParameter: key,
			},
		).Set(valueParsed)
	}
}
