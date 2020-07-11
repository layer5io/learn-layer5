package main

import (
	"encoding/json"
	"fmt"
	"testing"

	harness "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kudobuilder/kuttl/pkg/test"
	testutils "github.com/kudobuilder/kuttl/pkg/test/utils"
)

func main() {
	manifestDirs := []string{}
	testToRun := ""
	startKIND := false
	options := harness.TestSuite{}

	args := []string{"./"}

	options.TestDirs = args
	options.TestDirs = manifestDirs
	options.StartKIND = startKIND
	options.SkipDelete = false

	if options.KINDContext == "" {
		options.KINDContext = harness.DefaultKINDContext
	}

	if len(args) != 0 {
		options.TestDirs = args
	}

	testHandlers := make(map[string]map[string]test.CustomTest)
	// testHandlers["trafficAccess"] = getCustomTestsForAccess()

	testutils.RunTests("kudo", testToRun, options.Parallel, func(t *testing.T) {
		harness := test.Harness{
			TestSuite:        options,
			T:                t,
			SuiteCustomTests: testHandlers,
		}
		s, _ := json.MarshalIndent(options, "", "  ")
		fmt.Printf("Running integration tests with following options:\n%s\n", string(s))
		harness.Run()
	})
}
