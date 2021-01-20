package grpc

import (
	"context"
	"strconv"

	"github.com/layer5io/learn-layer5/smi-conformance/conformance"
	test_gen "github.com/layer5io/learn-layer5/smi-conformance/test-gen"
	"github.com/layer5io/service-mesh-performance/common"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	maeshConfig = &test_gen.Maesh{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}

	linkerdConfig = &test_gen.Linkerd{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}
	istioConfig = &test_gen.Istio{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}
	osmConfig = &test_gen.OSM{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}
)

// RunTest return conformance response
func (s *Service) RunTest(ctx context.Context, req *conformance.Request) (*conformance.Response, error) {
	var config test_gen.ServiceMesh

	config = linkerdConfig
	switch req.Mesh.Type {
	case smp.ServiceMesh_APP_MESH:
		config = linkerdConfig
		req.Annotations["linkerd.io/inject"] = "enabled"
	case smp.ServiceMesh_MAESH:
		config = maeshConfig
	case smp.ServiceMesh_ISTIO:
		config = istioConfig
		req.Labels["istio-injection"] = "enabled"
	case smp.ServiceMesh_OPEN_SERVICE_MESH:
		config = osmConfig
		req.Labels["openservicemesh.io/monitored-by"] = "osm"
	}

	result := test_gen.RunTest(config, req.Annotations, req.Labels)
	totalcases := 3
	failures := 0

	details := make([]*conformance.Detail, 0)
	for _, res := range result.Testsuite[0].Testcase {
		d := &conformance.Detail{
			Smispec:     res.Name,
			Specversion: "v1alpha1",
			Duration:    res.Time,
			Assertion:   strconv.Itoa(res.Assertions),
			// Capability: conformance.Capability(conformance.Capability_FULL),
			Status: conformance.ResultStatus(conformance.ResultStatus_PASSED),
			Result: &conformance.Result{
				Result: &conformance.Result_Message{
					Message: "",
				},
			},
		}
		if len(res.Failure.Text) > 2 {
			d.Result = &conformance.Result{
				Result: &conformance.Result_Error{
					Error: &common.CommonError{
						Code:                 "",
						Severity:             "",
						ShortDescription:     res.Failure.Text,
						LongDescription:      res.Failure.Message,
						ProbableCause:        "",
						SuggestedRemediation: "",
					},
				},
			}
			d.Status = conformance.ResultStatus(conformance.ResultStatus_FAILED)
			// d.Capability = "None"
			failures += 1
			if (res.Assertions - failures) > (res.Assertions / 2) {
				// d.Capability = "Half"
			}
		}
		details = append(details, d)
	}
	capability := conformance.Capability_NONE
	if totalcases-failures > totalcases/2 {
		capability = conformance.Capability_HALF
	} else if failures == 0 {
		capability = conformance.Capability_FULL
	}

	return &conformance.Response{
		Casespassed: strconv.Itoa(totalcases - failures),
		Passpercent: strconv.Itoa(((totalcases - failures) / totalcases) * 100),
		Details:     details,
		Mesh:        req.Mesh,
		Capability:  capability,
	}, nil
}
