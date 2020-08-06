package test_gen

import (
	"encoding/json"
	"fmt"
	"testing"

	harness "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kudobuilder/kuttl/pkg/report"
	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
)

type Results struct {
	Tests    int    `json:"tests"`
	Failures int    `json:"failures"`
	Time     string `json:"time"`
	Name     string `json:"name"`
	Testcase []struct {
		Classname  string `json:"classname"`
		Name       string `json:"name"`
		Time       string `json:"time"`
		Assertions int    `json:"assertions"`
		Failure    struct {
			Text    string `json:"text"`
			Message string `json:"message"`
		} `json:"failure,omitempty"`
	} `json:"testcase"`
}

func RunTest(meshConfig ServiceMesh, annotations map[string]string) Results {
	manifestDirs := []string{}
	output := Results{}
	results := &report.Testsuites{}

	// Run all testCases
	testToRun := ""
	// Run only traffic-split
	// testToRun := "traffic-split"

	startKIND := false
	options := harness.TestSuite{}

	args := []string{"./test-yamls/"}

	options.TestDirs = args
	options.Timeout = 30
	options.Parallel = 1
	options.TestDirs = manifestDirs
	options.StartKIND = startKIND
	options.SkipDelete = true

	if options.KINDContext == "" {
		options.KINDContext = harness.DefaultKINDContext
	}

	if len(args) != 0 {
		options.TestDirs = args
	}

	// annotations := make(map[string]string)
	// Namespace Injection
	// annotations["linkerd.io/inject"] = "enabled"

	serviceMeshConfObj := SMIConformance{
		SMObj: meshConfig,
	}

	testHandlers := make(map[string]map[string]test.CustomTest)
	testHandlers["traffic-access"] = serviceMeshConfObj.TrafficAccessGetTests()
	testHandlers["traffic-spec"] = serviceMeshConfObj.TrafficSpecGetTests()
	testHandlers["traffic-split"] = serviceMeshConfObj.TrafficSplitGetTests()

	testutils.RunTests("kudo", testToRun, options.Parallel, func(t *testing.T) {
		harness := test.Harness{
			TestSuite:            options,
			T:                    t,
			SuiteCustomTests:     testHandlers,
			NamespaceAnnotations: annotations,
		}

		// Runs the test using the inCluster kubeConfig (runs only when the code is running inside the pod)
		harness.InCluster = true

		s, _ := json.MarshalIndent(options, "", "  ")
		fmt.Printf("Running integration tests with following options:\n%s\n", string(s))
		results = harness.Run()
		data, _ := json.Marshal(results)
		// Results of the test
		fmt.Printf("Results :\n%v\n", string(data))
		err := json.Unmarshal([]byte(data), &output)
		if err != nil {
			fmt.Printf("Unable to unmarshal results")
		}
	})

	return output
}
