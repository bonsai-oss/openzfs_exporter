package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	MetricLabelName      = "name"
	MetricLabelPool      = "pool"
	MetricLabelParameter = "parameter"
)

// webListener - web server to access the metrics
func (app *application) webListener() {
	log.Printf("listening on %+q", app.listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	if err := app.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

var (
	app = application{}

	zpoolStats = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openzfs",
		Name:      "zpool_parameters",
		Help:      "sysctl openzfs parameters",
	}, []string{
		MetricLabelPool,
		MetricLabelName,
		MetricLabelParameter,
	})
)

type application struct {
	listenAddress string
	exportedPools arrayFlags
	server        *http.Server
	interval      time.Duration
}

type arrayFlags []string

func (af *arrayFlags) String() string {
	return "my string representation"
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

func init() {
	flag.DurationVar(&app.interval, "interval", 5*time.Second, "refresh interval for metrics")
	flag.StringVar(&app.listenAddress, "web.listen-address", ":8080", "address listening on")
	flag.Var(&app.exportedPools, "exported-pools", "address listening on")
	flag.Parse()

	app.server = &http.Server{
		Handler: nil,
		Addr:    app.listenAddress,
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ctx, cancelWorkers := context.WithCancel(context.Background())
	done := make(chan interface{})

	go app.webListener()

	startedWorkers := len(app.exportedPools) - 1
	stoppedWorkers := 0
	for _, pool := range app.exportedPools {
		go app.refreshWorker(ctx, done, pool)
	}

	<-sigChan

	cancelWorkers()

	for range done {
		if stoppedWorkers == startedWorkers {
			break
		}
		stoppedWorkers++
	}

	if err := app.server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}
