package test_gen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// App service names
const (
	SvcNameA = "app-a"
	SvcNameB = "app-b"
	SvcNameC = "app-c"
)

// App API endpoints
const (
	METRICS = "metrics"
	CALL    = "call"
	ECHO    = "echo"
)

// URLstruct is a part of the metrics exposed by the app
type URLstruct struct {
	URL     string
	Method  string
	Headers map[string]string
}

// MetricResponse is a part of the metrics exposed by the app
type MetricResponse struct {
	ReqReceived   []string
	RespSucceeded []URLstruct
	RespFailed    []URLstruct
}

// GetClusterIPs returns the ClusterIPs of various services exposed in the namespace
func GetClusterIPs(kubeClient client.Client, namespace string) (map[string]string, error) {
	deps := &v1.ServiceList{}
	err := kubeClient.List(context.TODO(), deps, client.InNamespace(namespace))
	if err != nil {
		return nil, err
	}
	ipMap := make(map[string]string)
	for _, svc := range deps.Items {
		ipMap[svc.Name] = svc.Spec.ClusterIP
	}
	return ipMap, nil
}

// GetHTTPClient returns a configured HTTP client
func GetHTTPClient() http.Client {
	return http.Client{
		Timeout: 30 * time.Second,
	}
}

// ClearMetrics remove all the resources
func ClearMetrics(hostname string, port string) error {
	url := fmt.Sprintf("http://%s:%s/%s", hostname, port, METRICS)
	httpClient := GetHTTPClient()
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Request failed")
	}
	return nil
}

// GetMetrics return a type MetricResponse
func GetMetrics(hostname string, port string) (*MetricResponse, error) {
	url := fmt.Sprintf("http://%s:%s/%s", hostname, port, METRICS)
	httpClient := GetHTTPClient()
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Request failed")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	metrics := MetricResponse{}
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}
	return &metrics, nil
}

func generatePOSTLoad(no int, url string, body []byte) error {
	hclient := GetHTTPClient()
	for i := 0; i < no; i++ {
		if _, err := hclient.Post(url, "application/json", bytes.NewReader(body)); err != nil {
			return err
		}
	}
	return nil
}

// ClearAllMetrics aggregate all the svc metrics
func ClearAllMetrics(clusterIPs map[string]string, smObj ServiceMesh) {
	ClearMetrics(clusterIPs[SvcNameA], smObj.SvcAGetPort())
	ClearMetrics(clusterIPs[SvcNameB], smObj.SvcBGetPort())
	ClearMetrics(clusterIPs[SvcNameC], smObj.SvcCGetPort())
}
