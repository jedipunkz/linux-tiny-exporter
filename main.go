package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jedipunkz/linux-tiny-exporter/internal/collector"
)

func main() {
	registry := prometheus.NewRegistry()

	collectors := []prometheus.Collector{
		collector.NewNetCollector(),
		collector.NewCPUCollector(),
		collector.NewDiskCollector(),
		collector.NewMemCollector(),
		collector.NewTempCollector(),
	}

	for _, c := range collectors {
		registry.MustRegister(c)
	}

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	fmt.Println("Starting exporter on :9101")
	if err := http.ListenAndServe(":9101", nil); err != nil {
		fmt.Printf("Error starting exporter: %v\n", err)
	}
}
