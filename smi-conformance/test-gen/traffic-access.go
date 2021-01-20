package test_gen

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (smi *SMIConformance) TrafficAccessGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["trafficDefault"] = smi.trafficBlocked
	testHandlers["trafficAllowed"] = smi.trafficAllow
	testHandlers["trafficBlocked"] = smi.trafficBlocked

	return testHandlers
}

func (smi *SMIConformance) trafficBlocked(
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

	// This test will make SERVICE A make a request to SERVICE B
	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	var jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smi.SMObj.SvcAGetPort(), CALL)
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

	// Validates if the request failed
	if !(len(metricsSvcA.RespFailed) == 1 && len(metricsSvcA.RespSucceeded) == 0) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response count")
	if metricsSvcA.RespFailed[0].URL != svcBTestURL {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response destination")

	Logger.Log("Done")
	return nil
}

func (smi *SMIConformance) trafficAllow(
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
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)
	if err!=nil{
		t.Fail()
		Logger.Logf("Error : %s", err.Error())
		return []error{err}
	}
	ClearAllMetrics(clusterIPs, smi.SMObj)

	// This test will make SERVICE A make a request to SERVICE B
	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	var jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smi.SMObj.SvcAGetPort(), CALL)
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

	// Validates if the request succeeded
	if !(len(metricsSvcA.RespFailed) == 0 && len(metricsSvcA.RespSucceeded) == 1) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response count")

	if metricsSvcA.RespSucceeded[0].URL != svcBTestURL {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response destination")

	metricsSvcB, err := GetMetrics(clusterIPs[SvcNameB], "9091")
	if err != nil {
		t.Fail()
		Logger.Log("Error: ", err)
		return []error{err}
	}

	Logger.Log("Service B : Response Failed", metricsSvcB.RespFailed)
	Logger.Log("Service B : Response Succeeded", metricsSvcB.RespSucceeded)
	Logger.Log("Service B : Requests Received", metricsSvcB.ReqReceived)

	if !(len(metricsSvcB.ReqReceived) == 1) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Request count")
	if metricsSvcB.ReqReceived[0] != "app-a" {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Request Source")

	Logger.Log("Done")
	return nil
}
