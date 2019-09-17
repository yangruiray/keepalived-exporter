package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/keepalived-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const (
	metricsPath = "/metrics"
)

var (
	listenPort int
)

// build metrics server for http request
func metricsServer(reg *prometheus.Registry) {
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		reg,
	}

	h := promhttp.HandlerFor(
		gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	http.HandleFunc(metricsPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	log.Infof("Start keepalived-exporter server at port: 9999")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", listenPort), nil); err != nil {
		log.Errorf("Error occurs on: %v", err)
		os.Exit(1)
	}
}

func init() {
	flag.IntVar(&listenPort, "port", 9999, "expose listen port")
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	// new a metric and register registry
	newMetrics := pkg.NewKeepalivedMetrics()
	reg := prometheus.NewRegistry()
	reg.MustRegister(newMetrics)

	metricsServer(reg)
}
