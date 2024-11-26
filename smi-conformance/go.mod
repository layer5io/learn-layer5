module github.com/layer5io/learn-layer5/smi-conformance

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/kr/text v0.2.0 // indirect
	github.com/kudobuilder/kuttl v0.0.0-00010101000000-000000000000
	// github.com/layer5io/meshkit v0.2.0
	github.com/layer5io/service-mesh-performance v0.3.2-0.20210122142912-a94e0658b021
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	k8s.io/api v0.20.0
	k8s.io/client-go v0.20.0
	sigs.k8s.io/controller-runtime v0.5.1
)

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
