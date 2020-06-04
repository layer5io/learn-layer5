package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var requestsReceived int
var responsesSucceeded int
var responsesFailed int
var mutex sync.Mutex

func exclusive(fn func()) {
	defer mutex.Unlock()
	mutex.Lock()
	fn()
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exclusive(func() { requestsReceived++ })
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

	if string(bytes) == "" {
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	host := data["host"]
	body := data["body"]

	if host == "" {
		return
	}

	var resp *http.Response
	if body != "" {
		resp, err = http.Post(host, "application/json", strings.NewReader(body))
	} else {
		resp, err = http.Get(host)
	}
	if err != nil {
		fmt.Printf("%v", err)
		http.Error(w, "Error parsing response body", http.StatusBadRequest)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		w.WriteHeader(http.StatusOK)
		exclusive(func() { responsesSucceeded++ })
	} else {
		exclusive(func() { responsesFailed++ })
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	w.Write(bytes)
}

func getMetrics(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not defined", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"responsesSucceeded": strconv.Itoa(responsesSucceeded),
		"responsesFailed":    strconv.Itoa(responsesFailed),
		"requestsReceived":   strconv.Itoa(requestsReceived),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func refreshMetrics(w http.ResponseWriter, req *http.Request) {
	responsesSucceeded = 0
	responsesFailed = 0
	requestsReceived = 0
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/call", MetricsMiddleware(http.HandlerFunc(call)))
	mux.Handle("/metrics/get", http.HandlerFunc(getMetrics))
	mux.Handle("/metrics/refresh", http.HandlerFunc(refreshMetrics))
	http.ListenAndServe(":9091", mux)
}
