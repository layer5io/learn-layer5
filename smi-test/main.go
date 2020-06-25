package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	harness "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func getCustomTestsForAccess() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)
	testHandlers["traffic"] = func(
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
		{
			ip2 := "http://" + ipMap["service-b"] + ":9091"
			var jsonStr = []byte(`{"url":"` + ip2 + `/echo", "body":"", "method": "GET", "headers": {"head": "tail"}}`)

			url := "http://" + ipMap["service-a"] + ":9091/call"

			req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				Logger.Log(err)
				return nil
			}
			// req.Header.Set("head", "tail")
			// client := &http.Client{}
			// resp, err := client.Do(req)
			if err != nil {
				Logger.Log(req)

			} else {
				// defer resp.Body.Close()
				// Logger.Log("response Status:", resp.Status)
				// Logger.Log("response Headers:", resp.Header)
				// body, _ := ioutil.ReadAll(resp.Body)
				// Logger.Log("response Body:", string(body))
			}
		}
		{
			ip1 := "http://" + ipMap["service-a"] + ":9091/metrics"

			resp, err := http.Get(ip1)
			if err != nil {
				Logger.Log(ip1)
				Logger.Log(err)
				return nil
			}
			defer resp.Body.Close()
			Logger.Log("A:response Status:", resp.Status)
			Logger.Log("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			Logger.Log("response Body:", string(body))
			var metric MetricResponse
			json.Unmarshal(body, &metric)
			if len(metric.RespSucceeded) != 0 {
				t.Fail()
				Logger.Log("Doesn't fail the request without headers")
			}

		}
		{
			ip2 := "http://" + ipMap["service-b"] + ":9091/metrics"

			resp, err := http.Get(ip2)
			if err != nil {
				Logger.Log(err)
				return nil
			}
			defer resp.Body.Close()
			Logger.Log("B:response Status:", resp.Status)
			Logger.Log("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			Logger.Log("response Body:", string(body))
		}
		{
			ip3 := "http://" + ipMap["service-c"] + ":9091/metrics"

			resp, err := http.Get(ip3)
			if err != nil {
				Logger.Log(err)
				return nil
			}
			defer resp.Body.Close()
			Logger.Log("c:response Status:", resp.Status)
			Logger.Log("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			Logger.Log("response Body:", string(body))
		}
		return nil
	}
	testHandlers["allow"] = func(
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
		{
			ip2 := "http://" + ipMap["service-b"] + ":9091"
			var jsonStr = []byte(`{"url":"` + ip2 + `/echo", "body":"", "method": "GET", "headers": {"head": "tail"}}`)

			url := "http://" + ipMap["service-a"] + ":9091/call"

			req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				Logger.Log(err)
				return nil
			}
			// req.Header.Set("head", "tail")
			// client := &http.Client{}
			// resp, err := client.Do(req)
			if err != nil {
				Logger.Log(req)

			} else {
				// defer resp.Body.Close()
				// Logger.Log("response Status:", resp.Status)
				// Logger.Log("response Headers:", resp.Header)
				// body, _ := ioutil.ReadAll(resp.Body)
				// Logger.Log("response Body:", string(body))
			}
		}
		{
			ip1 := "http://" + ipMap["service-a"] + ":9091/metrics"

			resp, err := http.Get(ip1)
			if err != nil {
				Logger.Log(ip1)
				Logger.Log(err)
				return nil
			}
			defer resp.Body.Close()
			Logger.Log("A:response Status:", resp.Status)
			Logger.Log("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			Logger.Log("response Body:", string(body))
			var metric MetricResponse
			json.Unmarshal(body, &metric)
			if len(metric.RespSucceeded) == 0 {
				t.Fail()
				Logger.Log("Doesn't fail the request without headers")
			}

		}
		{
			ip2 := "http://" + ipMap["service-b"] + ":9091/metrics"

			resp, err := http.Get(ip2)
			if err != nil {
				Logger.Log(err)
				return nil
			}
			defer resp.Body.Close()
			Logger.Log("B:response Status:", resp.Status)
			Logger.Log("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			Logger.Log("response Body:", string(body))
		}
		{
			ip3 := "http://" + ipMap["service-c"] + ":9091/metrics"

			resp, err := http.Get(ip3)
			if err != nil {
				Logger.Log(err)
				return nil
			}
			defer resp.Body.Close()
			Logger.Log("c:response Status:", resp.Status)
			Logger.Log("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			Logger.Log("response Body:", string(body))
		}
		return nil
	}
	return testHandlers
}

func main() {
	manifestDirs := []string{}
	testToRun := ""
	startKIND := false
	// skipClusterDelete := false
	options := harness.TestSuite{}
	// path, err := (os.Getwd())
	// if err != nil {
	// 	return
	// }
	args := []string{"./"}

	options.TestDirs = args
	// If a config is not set and kudo-test.yaml exists, set configPath to kudo-test.yaml.
	options.TestDirs = manifestDirs
	options.StartKIND = startKIND
	options.SkipDelete = false
	if options.KINDContext == "" {
		options.KINDContext = harness.DefaultKINDContext
	}
	// options.SkipDelete = skipClusterDelete
	if len(args) != 0 {
		options.TestDirs = args
	}
	testHandlers := make(map[string]map[string]test.CustomTest)
	testHandlers["trafficAccess"] = getCustomTestsForAccess()
	// testHandlers["trafficAccess"]["install"] = func(
	// 	t *testing.T,
	// 	namespace string,
	// 	client func(forceNew bool) (client.Client, error),
	// 	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	// 	Logger testutils.Logger,
	// ) []error {
	// 	Logger.Logf("This is a sample test log")
	// 	// time.Sleep(2 * time.Minute)
	// 	cl2, err := DiscoveryClient()
	// 	var runtimeInfo runtime.Object
	// 	err = cl2.RESTClient().Get().Do().Into(runtimeInfo)
	// 	if err != nil {
	// 		Logger.Log(err)
	// 		return nil
	// 	}
	// 	Logger.Log(runtimeInfo)
	// 	if err != nil {
	// 		Logger.Log(err)

	// 		t.Fail()
	// 	}
	// 	return nil
	// }
	testutils.RunTests("kudo", testToRun, options.Parallel, func(t *testing.T) {
		harness := test.Harness{
			TestSuite:        options,
			T:                t,
			SuiteCustomTests: testHandlers,
		}
		s, _ := json.MarshalIndent(options, "", "  ")
		fmt.Printf("Running integration tests with following options:\n%s\n", string(s))
		harness.Run()
	})
}
