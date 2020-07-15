package test_gen

import (
	"fmt"
	"testing"
	"time"
	"bytes"

	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (smi *SMIConformance) TrafficSplitGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["trafficDefault"] = smi.traffics

	return testHandlers
}

func (smi *SMIConformance) traffics(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	Logger.Log("Service sdfsdfsdf")
	time.Sleep(5 * time.Second)
	namespace = "kuttl-test-stage"
	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(false)
	if err != nil {
		t.Fail()
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)

	ClearMetrics(clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort())
	ClearMetrics(clusterIPs[SERVICE_B_NAME], smi.SMObj.SvcBGetPort())
	ClearMetrics(clusterIPs[SERVICE_C_NAME], smi.SMObj.SvcCGetPort())

	// call to metrics (allowed)
	svcBTestURLMetrics := "http://app-svc.kuttl-test-stage.svc.cluster.local.:9091/echo"
	jsonStr := []byte(`{"url":"` + svcBTestURLMetrics + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), CALL)
	for i := 0; i < 10; i++ {
		if _, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr)); err != nil {
			Logger.Log(err)
			break;
		}
	}
	if err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcA, err := GetMetrics(clusterIPs[SERVICE_A_NAME], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}

	Logger.Log("Service A : Response Falied", metricsSvcA.RespFailed)
	Logger.Log("Service A : Response Succeeded", metricsSvcA.RespSucceeded)
	Logger.Log("Service A : Requests Recieved", metricsSvcA.ReqReceived)

	metricsSvcB, err := GetMetrics(clusterIPs[SERVICE_B_NAME], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service B : Requests Recieved", metricsSvcB.ReqReceived)

	metricsSvcC, err := GetMetrics(clusterIPs[SERVICE_C_NAME], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service C : Requests Recieved", metricsSvcC.ReqReceived)

	Logger.Log("Done")
	return nil
}
