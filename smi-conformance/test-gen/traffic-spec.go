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

func (smi *SMIConformance) TrafficSpecGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["trafficPath"] = smi.traffic1
	// testHandlers["trafficMethod"] = smi.allow
	// testHandlers["trafficHeader"] = smi.traffic

	return testHandlers
}

func (smi *SMIConformance) traffic1(
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

	ClearMetrics(clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort())
	ClearMetrics(clusterIPs[SERVICE_B_NAME], smi.SMObj.SvcBGetPort())
	ClearMetrics(clusterIPs[SERVICE_C_NAME], smi.SMObj.SvcCGetPort())

	// call to metrics (allowed)
	svcBTestURLMetrics := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), METRICS)
	jsonStr := []byte(`{"url":"` + svcBTestURLMetrics + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		t.Fail()
		return []error{err}
	}

	// call to echo (blocked)
	svcBTestURLEcho := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	jsonStr = []byte(`{"url":"` + svcBTestURLEcho + `", "body":"", "method": "GET", "headers": {}}`)

	url = fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), CALL)
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
