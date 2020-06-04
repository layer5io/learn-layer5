package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func call(w http.ResponseWriter, req *http.Request) {}
func getMetrics(w http.ResponseWriter, req *http.Request) {}
func refreshMetrics(w http.ResponseWriter, req *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/call", MetricsMiddleware(http.HandlerFunc(hello)))
	mux.Handle("/metrics/get", http.HandlerFunc(getMetrics))
	mux.Handle("/metrics/refresh", http.HandlerFunc(refreshMetrics))
	http.ListenAndServe(":9091", mux)
}