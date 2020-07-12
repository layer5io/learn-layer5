package test_gen

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SMIConformance struct {
	SMObj ServiceMesh
}

func (smi *SMIConformance) TrafficAccessGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["defaultTraffic"] = smi.traffic
	testHandlers["allowTraffic"] = smi.allow

	return testHandlers
}

func (smi *SMIConformance) traffic(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(true)
	if err != nil {
		Logger.Log(err)
		return nil
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)
	var jsonStr = []byte(`{"url":` + fmt.Sprintf(`"%s/%s"`, smi.SMObj.SvcBGetInternalName(), ECHO) + `, "body":"", "method": "GET", "headers": {"head": "tail"}}`)

	Logger.Log(jsonStr)
	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs["app-a"], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		Logger.Log("ERror: ", err)
		return nil
	}

	metrics, err := GetMetrics(clusterIPs["app-a"], "9091")

	if err != nil {
		Logger.Log("ERror: ", err)
		return nil
	}
	Logger.Log("metrics: ", metrics)
	return nil
}

func (smi *SMIConformance) allow(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	hclient := http.Client{
		Timeout: 30 * time.Second,
	}
	cl2, err := clientFn(true)
	if err != nil {
		Logger.Log(err)
		return nil
	}
	deps := &v1.ServiceList{}
	err = cl2.List(context.TODO(), deps, client.InNamespace(namespace))
	if err != nil {
		Logger.Log(err)
		return nil
	}
	ipMap := make(map[string]string)
	for _, svc := range deps.Items {
		ipMap[svc.Name] = svc.Spec.ClusterIP
	}
	ip2 := "http://" + "app-b." + namespace + ".maesh" + ":9091"
	var jsonStr = []byte(`{"url":"` + ip2 + `/echo", "body":"", "method": "GET", "headers": {"head": "tail"}}`)
	Logger.Log(string(jsonStr))
	url := "http://" + ipMap["app-a"] + ":9091/call"
	resp, err := hclient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	resp, err = hclient.Get("http://" + ipMap["app-a"] + ":9091/metrics")
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Logger.Log("Body: ", string(data))
	return nil
}
