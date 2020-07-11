package test_gen

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

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

	testHandlers["traffic"] = smi.traffic
	testHandlers["allow"] = smi.allow

	return testHandlers
}

func (smi *SMIConformance) traffic(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
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
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	resp, err = http.Get("http://" + ipMap["app-a"] + ":9091/metrics")
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Logger.Log("Body: ", string(data))
	return nil
}

func (smi *SMIConformance) allow(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
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
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	resp, err = http.Get("http://" + ipMap["app-a"] + ":9091/metrics")
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Logger.Log("Body: ", string(data))
	return nil
}
