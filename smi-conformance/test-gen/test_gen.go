package test_gen

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	harness "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kudobuilder/kuttl/pkg/report"
	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
)

type Results struct {
	Name      string `json:"name"`
	Tests     int    `json:"tests"`
	Failures  int    `json:"failures"`
	Time      string `json:"time"`
	Testsuite []struct {
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
			} `json:"failure"`
		} `json:"testcase"`
	} `json:"testsuite"`
}

func RunTest(meshConfig ServiceMesh, annotations, labels map[string]string) Results {

	c := make(chan Results)
	go func() {
		manifestDirs := []string{}
		results := &report.Testsuites{}
		output := Results{}

		// Run all testCases
		testToRun := ""
		// Run only traffic-split
		// testToRun := "traffic-split"

		startKIND := false
		options := harness.TestSuite{}

		args := []string{"./test-yamls/"}

		options.TestDirs = args
		options.Timeout = 300
		options.Parallel = 1
		options.TestDirs = manifestDirs
		options.StartKIND = startKIND
		options.SkipDelete = false

		if options.KINDContext == "" {
			options.KINDContext = harness.DefaultKINDContext
		}

		serviceMeshConfObj := SMIConformance{
			SMObj: meshConfig,
		}

		if len(args) != 0 {
			options.TestDirs = args
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
				NamespaceLabels:      labels,
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
			c <- output
			time.Sleep(30 * time.Second)
		})
	}()
	select {
	case x := <-c:
		return x
	}
}
