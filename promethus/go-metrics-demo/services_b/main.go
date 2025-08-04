package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	serviceBCalls = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_b_calls_total",
			Help: "Total number of calls to service B",
		},
	)
)

func init() {
	// 显式注册指标
	prometheus.MustRegister(serviceBCalls)
}

func handler(w http.ResponseWriter, r *http.Request) {
	serviceBCalls.Inc()
	w.Write([]byte("Service B response"))
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Service B starting at :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
