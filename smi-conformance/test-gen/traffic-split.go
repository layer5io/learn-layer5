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

const reqNo = 50

// TrafficSplitGetTests return type of map[string]test.CustomTest
func (smi *SMIConformance) TrafficSplitGetTests() map[string]test.CustomTest {
	testHandlers := make(map[string]test.CustomTest)

	testHandlers["trafficDefault"] = smi.trafficSplitDefault
	testHandlers["trafficOnlyB"] = smi.trafficSplitOnlyB
	testHandlers["trafficOnlyC"] = smi.trafficSplitOnlyC
	testHandlers["trafficBGrtC"] = smi.trafficSplitBGrtC
	testHandlers["trafficCGrtB"] = smi.trafficSplitCGrtB

	return testHandlers
}

// Executes a Test Scenario where SERVICE A makes a request to the service named `app-svc`
func trafficSplitTestScenario(clusterIPs map[string]string, namespace string, smObj ServiceMesh) error {
	svcTrafficSplit := fmt.Sprintf("http://app-svc.%s.svc.cluster.local.:9091/%s", namespace, ECHO)
	jsonStr := []byte(`{"url":"` + svcTrafficSplit + `", "body":"", "method": "GET", "headers": {}}`)

	url := fmt.Sprintf("http://%s:%s/%s", clusterIPs[SvcNameA], smObj.SvcAGetPort(), CALL)

	return generatePOSTLoad(reqNo, url, jsonStr)
}

func (smi *SMIConformance) trafficSplitDefault(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

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

	if err = trafficSplitTestScenario(clusterIPs, namespace, smi.SMObj); err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcB, err := GetMetrics(clusterIPs[SvcNameB], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service B : Requests Received", metricsSvcB.ReqReceived)

	metricsSvcC, err := GetMetrics(clusterIPs[SvcNameC], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service C : Requests Received", metricsSvcC.ReqReceived)

	if len(metricsSvcB.ReqReceived) == 0 || len(metricsSvcC.ReqReceived) == 0 {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: Random Request count")

	Logger.Log("Done")
	return nil
}

func (smi *SMIConformance) trafficSplitOnlyB(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

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

	if err = trafficSplitTestScenario(clusterIPs, namespace, smi.SMObj); err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcB, err := GetMetrics(clusterIPs[SvcNameB], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service B : Requests Received", metricsSvcB.ReqReceived)

	metricsSvcC, err := GetMetrics(clusterIPs[SvcNameC], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service C : Requests Received", metricsSvcC.ReqReceived)

	if !(len(metricsSvcB.ReqReceived) == reqNo && len(metricsSvcC.ReqReceived) == 0) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: C Request count zero")

	Logger.Log("Done")
	return nil
}

func (smi *SMIConformance) trafficSplitOnlyC(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

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

	if err = trafficSplitTestScenario(clusterIPs, namespace, smi.SMObj); err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcB, err := GetMetrics(clusterIPs[SvcNameB], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service B : Requests Received", metricsSvcB.ReqReceived)

	metricsSvcC, err := GetMetrics(clusterIPs[SvcNameC], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service C : Requests Received", metricsSvcC.ReqReceived)

	if !(len(metricsSvcB.ReqReceived) == 0 && len(metricsSvcC.ReqReceived) == reqNo) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: B Request count zero")

	Logger.Log("Done")
	return nil
}

func (smi *SMIConformance) trafficSplitBGrtC(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

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

	if err = trafficSplitTestScenario(clusterIPs, namespace, smi.SMObj); err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcB, err := GetMetrics(clusterIPs[SvcNameB], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service B : Requests Received", metricsSvcB.ReqReceived)

	metricsSvcC, err := GetMetrics(clusterIPs[SvcNameC], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service C : Requests Received", metricsSvcC.ReqReceived)

	if !(len(metricsSvcB.ReqReceived) > len(metricsSvcC.ReqReceived)) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: B Request count greater than C")

	Logger.Log("Done")
	return nil
}

func (smi *SMIConformance) trafficSplitCGrtB(
	t *testing.T,
	namespace string,
	clientFn func(forceNew bool) (client.Client, error),
	DiscoveryClient func() (discovery.DiscoveryInterface, error),
	Logger testutils.Logger,
) []error {
	time.Sleep(5 * time.Second)

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

	if err = trafficSplitTestScenario(clusterIPs, namespace, smi.SMObj); err != nil {
		t.Fail()
		return []error{err}
	}

	metricsSvcB, err := GetMetrics(clusterIPs[SvcNameB], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service B : Requests Received", metricsSvcB.ReqReceived)

	metricsSvcC, err := GetMetrics(clusterIPs[SvcNameC], "9091")
	if err != nil {
		t.Fail()
		return []error{err}
	}
	Logger.Log("Service C : Requests Received", metricsSvcC.ReqReceived)

	if !(len(metricsSvcB.ReqReceived) < len(metricsSvcC.ReqReceived)) {
		t.Fail()
		return nil
	}
	Logger.Log("Validated: C Request count greater than B")

	Logger.Log("Done")
	return nil
}
