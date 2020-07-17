package main

import (
	"encoding/json"
	"fmt"
	"testing"

	harness "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
	test_gen "github.com/layer5/learn-layer5/smi-conformance/test-gen"
)

func main() {
	manifestDirs := []string{}
	testToRun := "traffic-split"
	startKIND := false
	options := harness.TestSuite{}

	args := []string{"./test-gen/test-yamls/"}

	options.TestDirs = args
	options.Timeout = 120
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

	// maeshConfig := test_gen.Maesh{
	// 	PortSvcA: "9091",
	// 	PortSvcB: "9091",
	// 	PortSvcC: "9091",
	// }

	linkerdConfig := test_gen.Linkerd{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}

	annotations := make(map[string]string)
	annotations["linkerd.io/inject"] = "enabled"

	serviceMeshConfObj := test_gen.SMIConformance{
		SMObj: linkerdConfig,
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
		s, _ := json.MarshalIndent(options, "", "  ")
		fmt.Printf("Running integration tests with following options:\n%s\n", string(s))
		results := harness.Run()
		data, _ := json.Marshal(results)
		fmt.Printf("Results :\n%v\n", string(data))
	})
}
