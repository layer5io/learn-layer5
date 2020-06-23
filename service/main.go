package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	logrus "github.com/sirupsen/logrus"
)

type MetricsResponse struct {
	URL     string
	Method  string
	Headers map[string]interface{}
}

type Metrics struct {
	ReqReceived   []string
	RespSucceeded []MetricsResponse
	RespFailed    []MetricsResponse
}

var metricsObj Metrics

var mutex sync.Mutex

func execExclusive(fn func()) {
	defer mutex.Unlock()
	mutex.Lock()
	fn()
}

const serviceID = "ServiceName"
const defaultRequestID = "Default"

var serviceName string

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execExclusive(func() {
			svcName := r.Header.Get(serviceID)
			if svcName == "" {
				svcName = "Unidentified"
			}
			metricsObj.ReqReceived = append(metricsObj.ReqReceived, svcName)
		})
		next.ServeHTTP(w, r)
	})
}

func call(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	defer r.Body.Close()

	var data map[string]interface{}
	bytes, err := ioutil.ReadAll(r.Body)
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

	url := data["url"].(string)
	method := data["method"].(string)
	headers := data["headers"]
	body := data["body"].(string)

	var req *http.Request
	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			logrus.Errorf("Error creating request %s", err.Error())
			http.Error(w, "Error creating request", http.StatusBadRequest)
			return
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			logrus.Errorf("Error creating request %s", err.Error())
			http.Error(w, "Error creating request", http.StatusBadRequest)
			return
		}
	}

	client := http.Client{}

	if headers != nil {
		headers := headers.(map[string]interface{})
		for key, val := range headers {
			req.Header.Add(key, val.(string))
		}
	}
	req.Header.Add(serviceID, serviceName)

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error completing the request %s", err.Error())
		http.Error(w, "Error completing the request", http.StatusInternalServerError)
		return
	}

	logrus.Debugf("Call response: %v", resp)
	if err != nil {
		http.Error(w, "Error making the calling", http.StatusBadRequest)
		logrus.Errorf("Error making the calling: %s", err.Error())
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		w.WriteHeader(http.StatusOK)
		execExclusive(func() {
			if headers != nil {
				headers := headers.(map[string]interface{})
				metricsObj.RespSucceeded = append(metricsObj.RespSucceeded, MetricsResponse{
					URL:     url,
					Method:  method,
					Headers: headers,
				})
			} else {
				metricsObj.RespSucceeded = append(metricsObj.RespSucceeded, MetricsResponse{
					URL:    url,
					Method: method,
				})
			}
		})
	} else {
		execExclusive(func() {
			if headers != nil {
				headers := headers.(map[string]interface{})
				metricsObj.RespFailed = append(metricsObj.RespFailed, MetricsResponse{
					URL:     url,
					Method:  method,
					Headers: headers,
				})
			} else {
				metricsObj.RespFailed = append(metricsObj.RespFailed, MetricsResponse{
					URL:    url,
					Method: method,
				})
			}
		})
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error parsing body: %s", err.Error())
		return
	}
	logrus.Debugf("Response body: %s", string(bytes))
	w.Write(bytes)
}

func echo(w http.ResponseWriter, req *http.Request) {
	req.Write(w)
}

func metrics(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		execExclusive(func() {
			if err := json.NewEncoder(w).Encode(&metricsObj); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
	} else if req.Method == http.MethodDelete {
		execExclusive(func() {
			metricsObj = Metrics{
				RespSucceeded: []MetricsResponse{},
				RespFailed:    []MetricsResponse{},
				ReqReceived:   []string{},
			}
		})
	} else {
		http.Error(w, "Method not defined", http.StatusBadRequest)
	}
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	serviceName = os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "Default"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}

	metricsObj = Metrics{
		RespSucceeded: []MetricsResponse{},
		RespFailed:    []MetricsResponse{},
		ReqReceived:   []string{},
	}

	mux := http.NewServeMux()

	mux.Handle("/call", MetricsMiddleware(http.HandlerFunc(call)))
	mux.Handle("/metrics", http.HandlerFunc(metrics))
	mux.Handle("/echo", MetricsMiddleware(http.HandlerFunc(echo)))
	logrus.Infof("Started serving at: %s", port)
	http.ListenAndServe(":"+port, mux)
}
