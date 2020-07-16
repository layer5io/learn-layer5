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

const (
	SERVICE_A_NAME = "app-a"
	SERVICE_B_NAME = "app-b"
	SERVICE_C_NAME = "app-c"
)

const (
	METRICS = "metrics"
	CALL    = "call"
	ECHO    = "echo"
)

type URLstruct struct {
	URL     string
	Method  string
	Headers map[string]string
}
type MetricResponse struct {
	ReqReceived   []string
	RespSucceeded []URLstruct
	RespFailed    []URLstruct
}

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

func GetHTTPClient() http.Client {
	return http.Client{
		Timeout: 30 * time.Second,
	}
}

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

func ClearAllMetrics(clusterIPs map[string]string, smObj ServiceMesh) {
	ClearMetrics(clusterIPs[SERVICE_A_NAME], smObj.SvcAGetPort())
	ClearMetrics(clusterIPs[SERVICE_B_NAME], smObj.SvcBGetPort())
	ClearMetrics(clusterIPs[SERVICE_C_NAME], smObj.SvcCGetPort())
}
