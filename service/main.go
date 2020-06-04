package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func call(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not defined", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	var data map[string]string
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	host := data["host"]
	body := data["body"]

	resp, err := http.Post(host, "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Printf("%v", err)
		http.Error(w, "Error parsing response body", http.StatusBadRequest)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		w.WriteHeader(http.StatusOK)
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	w.Write(bytes)
}

func getMetrics(w http.ResponseWriter, req *http.Request)     {}
func refreshMetrics(w http.ResponseWriter, req *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/call", MetricsMiddleware(http.HandlerFunc(call)))
	mux.Handle("/metrics/get", http.HandlerFunc(getMetrics))
	mux.Handle("/metrics/refresh", http.HandlerFunc(refreshMetrics))
	http.ListenAndServe(":9091", mux)
}
