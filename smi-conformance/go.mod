module github.com/layer5io/learn-layer5/smi-conformance

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/kudobuilder/kuttl v0.0.0-00010101000000-000000000000
	github.com/layer5io/gokit v0.1.12
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.23.0
	k8s.io/api v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.5.1
)
