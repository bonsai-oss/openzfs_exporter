package main

import (
	"context"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// webListener - web server to access the metrics
func (app application) webListener() {
	http.Handle("/metrics", promhttp.Handler())
	if err := app.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

var (
	app = application{}

	zpool_stats = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fsrv_bsd_userland_version",
		Help: "version of the FreeBSD userland",
	})
)

type application struct {
	listenAddress string
	exportedPools arrayFlags
	server        *http.Server
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
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

	go app.webListener()

	<-sigChan

	if err := app.server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}
