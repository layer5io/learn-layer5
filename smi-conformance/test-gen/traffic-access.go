package test_gen

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
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
	testHandlers["blockTraffic"] = smi.traffic

	return testHandlers
}

func (smi *SMIConformance) traffic(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	namespace  = "kuttl-test-stage"
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

	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	var jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {"head": "tail"}}`)

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
	Logger.Log(metricsSvcA.RespFailed)
	Logger.Log(metricsSvcA.RespSucceeded)
	if !(len(metricsSvcA.RespFailed) == 1 && len(metricsSvcA.RespSucceeded) == 0) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Response count")
	if metricsSvcA.RespFailed[0].URL != svcBTestURL {
		// t.Fail()
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
	namespace  = "kuttl-test-stage"
	httpClient := GetHTTPClient()
	kubeClient, err := clientFn(false)
	if err != nil {
		// t.Fail()
		return []error{err}
	}
	clusterIPs, err := GetClusterIPs(kubeClient, namespace)

	ClearMetrics(clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort())
	ClearMetrics(clusterIPs[SERVICE_B_NAME], smi.SMObj.SvcBGetPort())
	ClearMetrics(clusterIPs[SERVICE_C_NAME], smi.SMObj.SvcCGetPort())

	svcBTestURL := fmt.Sprintf("%s/%s", smi.SMObj.SvcBGetInternalName(namespace), ECHO)
	var jsonStr = []byte(`{"url":"` + svcBTestURL + `", "body":"", "method": "GET", "headers": {"head": "tail"}}`)

	Logger.Log(string(jsonStr))
	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SERVICE_A_NAME], smi.SMObj.SvcAGetPort(), CALL)
	_, err = httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		// t.Fail()
		return []error{err}
	}

	metricsSvcA, err := GetMetrics(clusterIPs[SERVICE_A_NAME], "9091")
	if err != nil {
		// t.Fail()
		return []error{err}
	}
	if !(len(metricsSvcA.RespFailed) == 0 && len(metricsSvcA.RespSucceeded) == 1) {
		// t.Fail()
		return nil
	}
	Logger.Log("Validated: Response count")
	
	if metricsSvcA.RespSucceeded[0].URL != svcBTestURL {
		// t.Fail()
		return nil
	}
	Logger.Log("Validated: Response destination")

	metricsSvcB, err := GetMetrics(clusterIPs[SERVICE_B_NAME], "9091")
	if err != nil {
		// t.Fail()
		return []error{err}
	}
	Logger.Log(metricsSvcB)
	if !(len(metricsSvcB.ReqReceived) == 1) {
		// t.Fail()
		return nil
	}
	Logger.Log("Validated: Request count")
	if metricsSvcB.ReqReceived[0] != "app-a" {
		// t.Fail()
		return nil
	}
	Logger.Log("Validated: Request Source")

	Logger.Log("Done")
	return nil
}
