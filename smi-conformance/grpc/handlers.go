package grpc

import (
	"context"
	"strconv"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/layer5io/learn-layer5/smi-conformance/conformance"
	test_gen "github.com/layer5io/learn-layer5/smi-conformance/test-gen"
	service "github.com/layer5io/service-mesh-performance/service"
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
		req.Mesh.Annotations["linkerd.io/inject"] = "enabled"
	case smp.ServiceMesh_MAESH:
		config = maeshConfig
	case smp.ServiceMesh_ISTIO:
		config = istioConfig
		req.Mesh.Labels["istio-injection"] = "enabled"
	case smp.ServiceMesh_OPEN_SERVICE_MESH:
		config = osmConfig
		req.Mesh.Labels["openservicemesh.io/monitored-by"] = "osm"

	}

	result := test_gen.RunTest(config, req.Mesh.Annotations, req.Mesh.Labels)
	totalcases := 3
	failures := 0

	details := make([]*conformance.Detail, 0)
	for _, res := range result.Testsuite[0].Testcase {
		d := &conformance.Detail{
			Smispec:     res.Name,
			Specversion: "v1alpha1",
			Assertion:   strconv.Itoa(res.Assertions),
			Duration:    res.Time,
			Capability:  conformance.Capability_FULL,
			Status:      conformance.ResultStatus_PASSED,
			Result: &conformance.Result{
				Result: &conformance.Result_Message{
					Message: "All test passed",
				},
			},
		}
		if len(res.Failure.Text) > 2 {
			d.Result = &conformance.Result{
				Result: &conformance.Result_Error{
					Error: &service.CommonError{
						Code:                 "",
						Severity:             "",
						ShortDescription:     res.Failure.Text,
						LongDescription:      res.Failure.Message,
						ProbableCause:        "",
						SuggestedRemediation: "",
					},
				},
			}
			d.Status = conformance.ResultStatus_FAILED
			d.Capability = conformance.Capability_NONE
			failures += 1
			if (res.Assertions - failures) > (res.Assertions / 2) {
				d.Capability = conformance.Capability_HALF
			}
		}
		details = append(details, d)
	}

	return &conformance.Response{
		Casespassed: strconv.Itoa(totalcases - failures),
		Passpercent: strconv.Itoa(((totalcases - failures) / totalcases) * 100),
		Mesh:        req.Mesh,
		Details:     details,
	}, nil
}

func (s *Service) Info(context.Context, *empty.Empty) (*service.ServiceInfo, error) {
	return &service.ServiceInfo{}, nil
}

func (s *Service) Health(context.Context, *empty.Empty) (*service.ServiceHealth, error) {
	return &service.ServiceHealth{}, nil
}
