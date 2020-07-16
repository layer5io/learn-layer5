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

	testHandlers["trafficDefault"] = smi.traffic
	testHandlers["trafficAllowed"] = smi.allow
	testHandlers["trafficBlocked"] = smi.traffic

	return testHandlers
}

func (smi *SMIConformance) traffic(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)
	namespace = "kuttl-test-stage"
	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(false)
	if err != nil {
		t.Fail()
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)

	ClearAllMetrics(clusterIPs, smi.SMObj)

	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	var jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcA, err := GetMetrics(clusterIPs[SERVICE_A_NAME], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}

	Logger.Log("Service A : Response Failed", metricsSvcA.RespFailed)
	Logger.Log("Service A : Response Succeeded", metricsSvcA.RespSucceeded)
	Logger.Log("Service A : Requests Recieved", metricsSvcA.ReqReceived)

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

func (smi *SMIConformance) allow(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)
	namespace = "kuttl-test-stage"
	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(false)
	if err != nil {
		t.Fail()
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)

	ClearAllMetrics(clusterIPs, smi.SMObj)

	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	var jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcA, err := GetMetrics(clusterIPs[SERVICE_A_NAME], "9091")

	if err != nil {
		t.Fail()
		return []error{err}
	}

	Logger.Log("Service A : Response Failed", metricsSvcA.RespFailed)
	Logger.Log("Service A : Response Succeeded", metricsSvcA.RespSucceeded)
	Logger.Log("Service A : Requests Recieved", metricsSvcA.ReqReceived)

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

	metricsSvcB, err := GetMetrics(clusterIPs[SERVICE_B_NAME], "9091")
	if err != nil {
		t.Fail()
		Logger.Log("Error: ", err)
		return []error{err}
	}

	Logger.Log("Service B : Response Failed", metricsSvcB.RespFailed)
	Logger.Log("Service B : Response Succeeded", metricsSvcB.RespSucceeded)
	Logger.Log("Service B : Requests Recieved", metricsSvcB.ReqReceived)

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
