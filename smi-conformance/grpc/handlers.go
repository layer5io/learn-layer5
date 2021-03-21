package grpc

import (
	"context"
	"regexp"
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
	case smp.ServiceMesh_LINKERD:
		config = linkerdConfig
		req.Mesh.Annotations["linkerd.io/inject"] = "enabled"
	case smp.ServiceMesh_APP_MESH:
		config = linkerdConfig
		req.Mesh.Labels["appmesh.k8s.aws/sidecarInjectorWebhook"] = "enabled"
	case smp.ServiceMesh_MAESH:
		config = maeshConfig
	case smp.ServiceMesh_ISTIO:
		config = istioConfig
		req.Mesh.Labels["istio-injection"] = "enabled"
	case smp.ServiceMesh_OPEN_SERVICE_MESH:
		config = osmConfig
		req.Mesh.Labels["openservicemesh.io/monitored-by"] = "osm"
	case smp.ServiceMesh_KUMA:
		req.Mesh.Annotations["kuma.io/sidecar-injection"] = "enabled"
	case smp.ServiceMesh_NGINX_SERVICE_MESH:
		req.Mesh.Annotations["njector.nsm.nginx.com/auto-inject"] = "true"

	}

	result := test_gen.RunTest(config, req.Mesh.Annotations, req.Mesh.Labels)
	totalSteps := 24
	totalFailures := 0
	stepsCount := map[string]int{
		"traffic-access": 7,
		"traffic-split":  11,
		"traffic-spec":   6,
	}
	specVersion := map[string]string{
		"traffic-access": "v0.6.0/v1alpha3",
		"traffic-split":  "v0.6.0/v1alpha4",
		"traffic-spec":   "v0.6.0/v1alpha4",
	}

	details := make([]*conformance.Detail, 0)
	for _, res := range result.Testsuite[0].Testcase {
		d := &conformance.Detail{
			Smispec:     res.Name,
			Specversion: specVersion[res.Name],
			Assertion:   strconv.Itoa(stepsCount[res.Name]),
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

			// A hacky way to see the testStep Failed, since KUDO only provides it in Failure.Message
			re := regexp.MustCompile(`[0-9]+`)
			if res.Failure.Message != "" {
				stepFailed := re.FindAllString(res.Failure.Message, 1)
				if len(stepFailed) != 0 {
					passed, _ := strconv.Atoi(stepFailed[0])
					passed = passed - 1
					failures := stepsCount[res.Name] - passed
					totalFailures += failures
					if (passed) >= (stepsCount[res.Name] / 2) {
						d.Capability = conformance.Capability_HALF
					}
				}
			}
		}
		details = append(details, d)
	}

	return &conformance.Response{
		Casespassed: strconv.Itoa(totalSteps - totalFailures),
		Passpercent: strconv.FormatFloat(float64(totalSteps-totalFailures)/float64(totalSteps)*100, 'f', 2, 64),
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
