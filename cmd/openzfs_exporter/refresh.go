package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"openzfs_exporter/internal/dataset"
	"time"
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
			for _, ds := range datasets {

				// TODO: make loop for dataset values
				zpool_stats.With(prometheus.Labels{"name": ds.Name, "node": ds.Name})
			}
			fmt.Println(pool)
			sleepCounter = 0
		}
		sleepCounter++
		time.Sleep(1 * time.Second)
	}
}
