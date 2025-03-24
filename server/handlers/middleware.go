package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"tarantool-kv/server/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(data)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}

		defer func() {
			duration := time.Since(start).Seconds()
			status := rw.status
			if status == 0 {
				status = http.StatusOK
			}

			log.Printf(
				"Method=%s URL=%s RemoteAddr=%s Status=%d Duration=%.3f",
				r.Method, r.URL.Path, r.RemoteAddr, status, duration,
			)

			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", status)).Inc()
			metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		}()

		next.ServeHTTP(rw, r)
	})
}
