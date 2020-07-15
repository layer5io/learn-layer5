package test_gen

import (
	"fmt"
	"testing"
	"time"

	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (smi *SMIConformance) TrafficSplitGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["trafficPath"] = smi.trafficPath

	return testHandlers
}

func (smi *SMIConformance) traffics(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)
	namespace = "kuttl-test-stage"
	// httpClient := GetHTTPClient()
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
	// svcBTestURLMetrics := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), METRICS)
	// jsonStr := []byte(`{"url":"` + svcBTestURLMetrics + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), ECHO)
	err = generateLoad(10, url)
	if err != nil {
		t.Fail()
		return []error{err}
	}

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
