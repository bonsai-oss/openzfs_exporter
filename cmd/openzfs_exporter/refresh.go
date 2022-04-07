package main

import (
	"context"
	"fmt"
	"openzfs_exporter/internal/dataset"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (app *application) refreshWorker(ctx context.Context, done chan<- interface{}, pool string) {
	var sleepCounter int
	for {
		select {
		case <-ctx.Done():
			done <- nil
			return
		default:
			if sleepCounter/int(app.interval.Seconds()) == 0 {
				break
			}
			datasets := dataset.DetectDatasets(pool)
			st := time.Now()
			for _, ds := range datasets {
				ds.ParseValues()
				for key, value := range ds.Parameter {
					zpoolStats.With(
						prometheus.Labels{
							MetricLabelName:      ds.Name,
							MetricLabelPool:      pool,
							MetricLabelParameter: key,
						},
					).Set(float64(value))
				}
			}
			fmt.Println(pool, time.Since(st))
			sleepCounter = 0
		}
		sleepCounter++
		time.Sleep(1 * time.Second)
	}
}
