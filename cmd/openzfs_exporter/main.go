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

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// webListener - web server to access the metrics
func (app *application) webListener() {
	log.Printf("listening on %+q", app.listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	if err := app.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

// application is the settings and state holding structure of the exporter
type application struct {
	listenAddress    string
	exportedPools    arrayFlags
	server           *http.Server
	interval         time.Duration
	useAutodiscovery bool
}

// generate the main application instance
var app = application{}

func init() {
	// parse command line flags
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

	// autodiscovery mode; append auto discovered hosts to `app.exportedPools`
	if app.useAutodiscovery {
		pools, err := pool.Discover()
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range pools {
			app.exportedPools = append(app.exportedPools, p.Name)
		}
	}

	// only start exporter if pools are set
	if len(app.exportedPools) == 0 {
		log.Fatalln("no pools to check")
	}
}

func main() {
	// capture input signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ctx, cancelWorkers := context.WithCancel(context.Background())
	done := make(chan interface{})

	// start http listener
	go app.webListener()

	startedWorkers := len(app.exportedPools) - 1
	stoppedWorkers := 0

	// start one worker per pool
	for _, p := range app.exportedPools {
		log.Printf("monitoring p %+q\n", p)
		go app.refreshWorker(ctx, done, p)
	}

	// wait for incoming interrupts
	<-sigChan

	cancelWorkers()

	// wait for all workers to be stopped
	for range done {
		if stoppedWorkers == startedWorkers {
			break
		}
		stoppedWorkers++
	}

	// stop the http service
	if err := app.server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}
