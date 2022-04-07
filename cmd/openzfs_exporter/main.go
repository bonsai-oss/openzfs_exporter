package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"openzfs_exporter/internal/pool"
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
		Subsystem: "zfs",
		Name:      "zfs_parameters",
		Help:      "sysctl openzfs dataset parameters",
	}, []string{
		MetricLabelPool,
		MetricLabelName,
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

type application struct {
	listenAddress    string
	exportedPools    arrayFlags
	server           *http.Server
	interval         time.Duration
	useAutodiscovery bool
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
	flag.BoolVar(&app.useAutodiscovery, "discover-pools", false, "use autodiscovery for zfs pools")
	flag.Var(&app.exportedPools, "exported-pools", "address listening on")
	flag.Parse()

	// create http server and assign address
	app.server = &http.Server{
		Handler: nil,
		Addr:    app.listenAddress,
	}

	// autodiscovery mode
	if app.useAutodiscovery {
		pools, err := pool.Discover()
		if err != nil {
			log.Fatal(err)
		}
		for _, pool := range pools {
			app.exportedPools = append(app.exportedPools, pool.Name)
		}
	}

	if len(app.exportedPools) == 0 {
		log.Fatalln("no pools to check")
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
		log.Printf("monitoring pool %+q\n", pool)
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
