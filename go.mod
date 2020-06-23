module github.com/layer5/learn-layer5

go 1.13

replace github.com/kudobuilder/kuttl => /home/naveen/go/src/github.com/kanishkarj/kuttl

require (
	github.com/kudobuilder/kuttl v0.0.0-00010101000000-000000000000
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.5.1
)
