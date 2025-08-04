package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 声明指标
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpDurations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request durations",
			Buckets: []float64{0.1, 0.2, 0.5, 1, 2},
		},
		[]string{"method", "path"},
	)
)

// 自定义ResponseWriter捕获状态码
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func init() {
	// 显式注册指标
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpDurations)
}

// Prometheus
func promMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: 200}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", recorder.status)

		// 记录指标
		httpDurations.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
	})
}

func callServiceB() {
	resp, err := http.Get("http://service_b:8081/api")
	if err != nil {
		log.Printf("Error calling service B: %v", err)
		return
	}
	defer resp.Body.Close()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	callServiceB()
	w.Write([]byte("Service A response"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.Handle("/metrics", promhttp.Handler())

	wrappedMux := promMiddleware(mux)

	log.Println("Service A starting at :8080...")
	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}
