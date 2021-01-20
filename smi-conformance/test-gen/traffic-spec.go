package test_gen

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (smi *SMIConformance) TrafficSpecGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["trafficPath"] = smi.trafficPath
	testHandlers["trafficMethod"] = smi.trafficMethod

	return testHandlers
}

func (smi *SMIConformance) trafficPath(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(false)
	if err != nil {
		t.Fail()
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)
	if err!=nil{
		t.Fail()
		return []error{err}
	}
	ClearAllMetrics(clusterIPs, smi.SMObj)

	// call to SERVICE B metrics (allowed)
	svcBTestURLMetrics := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), METRICS)
	jsonStr := []byte(`{"url":"` + svcBTestURLMetrics + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	Logger.Logf("URL : \n", url)
	Logger.Logf("Body : \n", string(jsonStr))
	if err != nil {
		t.Fail()
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}

	// call to SERVICE B echo (blocked)
	svcBTestURLEcho := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	jsonStr = []byte(`{"url":"` + svcBTestURLEcho + `", "body":"", "method": "GET", "headers": {}}`)

	url = fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcA, err := GetMetrics(clusterIPs[SvcNameA], "9091")
	if err != nil {
		t.Fail()
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}

	Logger.Log("Service A : Response Failed", metricsSvcA.RespFailed)
	Logger.Log("Service A : Response Succeeded", metricsSvcA.RespSucceeded)
	Logger.Log("Service A : Requests Received", metricsSvcA.ReqReceived)

	// validates the requests that failed and the ones that succeeded
	if !(len(metricsSvcA.RespFailed) == 1 && len(metricsSvcA.RespSucceeded) == 1) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response count")
	if metricsSvcA.RespSucceeded[0].URL != svcBTestURLMetrics {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Allowed Path")
	if metricsSvcA.RespFailed[0].URL != svcBTestURLEcho {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Disallowed Path")

	Logger.Log("Done")
	return nil
}

func (smi *SMIConformance) trafficMethod(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(false)
	if err != nil {
		t.Fail()
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)
	if err != nil {
		t.Fail()
		return []error{err}
	}
	ClearAllMetrics(clusterIPs, smi.SMObj)

	// GET to echo (allowed)
	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	jsonStr := []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	Logger.Logf("URL : \n", url)
	Logger.Logf("Body : \n", string(jsonStr))
	if err != nil {
		t.Fail()
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}

	// POST to echo (blocked)
	jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "POST", "headers": {}}`)

	url = fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	Logger.Logf("URL : \n", url)
	Logger.Logf("Body : \n", string(jsonStr))
	if err != nil {
		t.Fail()
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}

	metricsSvcA, err := GetMetrics(clusterIPs[SvcNameA], "9091")
	if err != nil {
		t.Fail()
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}

	Logger.Log("Service A : Response Failed", metricsSvcA.RespFailed)
	Logger.Log("Service A : Response Succeeded", metricsSvcA.RespSucceeded)
	Logger.Log("Service A : Requests Received", metricsSvcA.ReqReceived)

	// validates the requests that failed and the ones that succeeded
	if !(len(metricsSvcA.RespFailed) == 1 && len(metricsSvcA.RespSucceeded) == 1) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response count")
	if metricsSvcA.RespSucceeded[0].Method != http.MethodGet {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Allowed Method")
	if metricsSvcA.RespFailed[0].Method != http.MethodPost {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Disallowed Method")

	Logger.Log("Done")
	return nil
}
