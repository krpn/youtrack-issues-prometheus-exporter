package main

import (
	"fmt"
	"github.com/krpn/youtrack-issues-prometheus-exporter/config"
	"github.com/krpn/youtrack-issues-prometheus-exporter/httpwrap"
	"github.com/krpn/youtrack-issues-prometheus-exporter/monitoring"
	"github.com/krpn/youtrack-issues-prometheus-exporter/prometheus"
	"github.com/krpn/youtrack-issues-prometheus-exporter/youtrack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	b, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}

	c, err := config.New(b)
	if err != nil {
		panic(err)
	}

	var (
		client       = httpwrap.New(&http.Client{Timeout: time.Duration(c.RequestTimeoutSeconds) * time.Second})
		refreshDelay = time.Duration(c.RefreshDelaySeconds) * time.Second
	)

	yt, err := youtrack.New(c.Endpoint, c.Token, client)
	if err != nil {
		panic(err)
	}

	monitor := monitoring.New(yt, prometheus.New(), c.Queries)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		panic(http.ListenAndServe(fmt.Sprintf(":%v", c.ListenPort), nil))
	}()

	for {
		monitor.RefreshMetrics()
		time.Sleep(refreshDelay)
	}
}
