package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"tarantool-kv/server/handlers"
	"tarantool-kv/server/metrics"
	"tarantool-kv/server/storage"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	host := os.Getenv("TARANTOOL_HOST")
	port := os.Getenv("TARANTOOL_PORT")
	user := os.Getenv("TARANTOOL_USER")
	password := os.Getenv("TARANTOOL_PASSWORD")

	conn, err := storage.NewTarantoolConnection(ctx, host, port, user, password)
	if err != nil {
		log.Fatalf("Failed to connect to Tarantool: %s", err)
	}
	defer conn.Close()

	metrics.RegisterMetrics()

	r := mux.NewRouter()
	r.HandleFunc("/kv", handlers.PostKV(conn)).Methods("POST")
	r.HandleFunc("/kv/{id}", handlers.PutKV(conn)).Methods("PUT")
	r.HandleFunc("/kv/{id}", handlers.GetKV(conn)).Methods("GET")
	r.HandleFunc("/kv/{id}", handlers.DeleteKV(conn)).Methods("DELETE")
	r.Handle("/metrics", promhttp.Handler())

	r.Use(handlers.LoggingMiddleware)

	serverPort := ":8080"
	log.Printf("Starting server on port %s\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, r))
}
