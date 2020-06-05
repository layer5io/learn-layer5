package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	logrus "github.com/sirupsen/logrus"
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
		logrus.Errorf("Method not defined")
		return
	}

	defer req.Body.Close()

	var data map[string]string
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		logrus.Errorf("Error reading body: %s", err.Error())
		return
	}
	logrus.Debugf("request body: %v", string(bytes))

	if string(bytes) == "" {
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		logrus.Errorf("Error parsing body: %s", err.Error())
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

	logrus.Debugf("Call response: %v", resp)
	if err != nil {
		http.Error(w, "Error making the calling", http.StatusBadRequest)
		logrus.Errorf("Error making the calling: %s", err.Error())
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
		logrus.Errorf("Error parsing body: %s", err.Error())
		return
	}
	logrus.Debugf("Response body: %s", string(bytes))
	w.Write(bytes)
}

func metrics(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"responsesSucceeded": strconv.Itoa(responsesSucceeded),
			"responsesFailed":    strconv.Itoa(responsesFailed),
			"requestsReceived":   strconv.Itoa(requestsReceived),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if req.Method == http.MethodDelete {
		responsesSucceeded = 0
		responsesFailed = 0
		requestsReceived = 0
	} else {
		http.Error(w, "Method not defined", http.StatusBadRequest)
	}
}

func main() {
	logrus.SetOutput(os.Stdout)

	mux := http.NewServeMux()
	mux.Handle("/call", MetricsMiddleware(http.HandlerFunc(call)))
	mux.Handle("/metrics", http.HandlerFunc(metrics))
	http.ListenAndServe(":9091", mux)
}
