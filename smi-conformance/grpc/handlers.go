package grpc

import (
	"context"
	"strconv"

	"github.com/layer5io/learn-layer5/smi-conformance/conformance"
	test_gen "github.com/layer5io/learn-layer5/smi-conformance/test-gen"
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
)

func (s *Service) RunTest(ctx context.Context, req *conformance.Request) (*conformance.Response, error) {
	results := make([]*conformance.SingleTestResult, 0)
	var config test_gen.ServiceMesh

	switch req.Meshname {
	case "linkerd":
		config = linkerdConfig
	case "maesh":
		config = maeshConfig
	}

	result := test_gen.RunTest(config, req.Annotations)
	for _, res := range result {
		results = append(results, &conformance.SingleTestResult{
			Name:            res.Name,
			TestCasesPassed: strconv.Itoa(res.Passed),
			TotalCases:      strconv.Itoa(res.Total),
			Message:         res.Message,
		})
	}

	return &conformance.Response{
		SingleTestResult: results,
	}, nil
}
