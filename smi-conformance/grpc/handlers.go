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

func (s *Service) RunTest(ctx context.Context, req *conformance.Request) (*conformance.Response, error) {
	var config test_gen.ServiceMesh

	config = linkerdConfig
	switch req.Meshname {
	case "linkerd":
		config = linkerdConfig
		req.Annotations["linkerd.io/inject"] = "enabled"
	case "maesh":
		config = maeshConfig
	case "istio":
		config = istioConfig
		req.Labels["istio-injection"] = "enabled"
	case "osm":
		config = osmConfig
		req.Labels["openservicemesh.io/monitored-by"] = "osm"
	}

	result := test_gen.RunTest(config, req.Annotations, req.Labels)
	totalcases := 3
	failures := 0

	details := make([]*conformance.Detail, 0)
	for _, res := range result.Testsuite[0].Testcase {
		d := &conformance.Detail{
			Smispec:    res.Name,
			Time:       res.Time,
			Assertions: strconv.Itoa(res.Assertions),
			Capability: "Full",
			Status:     "Passing",
		}
		if len(res.Failure.Text) > 2 {
			d.Reason = res.Failure.Text
			d.Result = res.Failure.Message
			d.Status = "Failing"
			d.Capability = "None"
			failures += 1
			if (res.Assertions - failures) > (res.Assertions / 2) {
				d.Capability = "Half"
			}
		}
		details = append(details, d)
	}

	return &conformance.Response{
		Casespassed: strconv.Itoa(totalcases - failures),
		Passpercent: strconv.Itoa(((totalcases - failures) / totalcases) * 100),
		Details:     details,
	}, nil
}
