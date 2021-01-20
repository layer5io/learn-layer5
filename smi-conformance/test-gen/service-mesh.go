package test_gen

import "fmt"

// SMIConformance holds the SMI conformance tests
type SMIConformance struct {
	SMObj ServiceMesh
}

// ServiceMesh provides an abstract interface for different service meshes.
// This is required as each service mesh has different ways to expose their internals.
type ServiceMesh interface {
	SvcAGetInternalName(string) string
	SvcBGetInternalName(string) string
	SvcCGetInternalName(string) string

	SvcAGetPort() string
	SvcBGetPort() string
	SvcCGetPort() string
}

type Maesh struct {
	PortSvcA string
	PortSvcB string
	PortSvcC string
}

func (sm Maesh) SvcAGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s.maesh:%s", SvcNameA, namespace, sm.PortSvcA)
}

func (sm Maesh) SvcBGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s.maesh:%s", SvcNameB, namespace, sm.PortSvcB)
}

func (sm Maesh) SvcCGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s.maesh:%s", SvcNameC, namespace, sm.PortSvcC)
}

func (sm Maesh) SvcAGetPort() string {
	return sm.PortSvcA
}

func (sm Maesh) SvcBGetPort() string {
	return sm.PortSvcB
}

func (sm Maesh) SvcCGetPort() string {
	return sm.PortSvcC
}

type Linkerd struct {
	PortSvcA string
	PortSvcB string
	PortSvcC string
}

func (sm Linkerd) SvcAGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameA, namespace, sm.PortSvcA)
}

func (sm Linkerd) SvcBGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameB, namespace, sm.PortSvcB)
}

func (sm Linkerd) SvcCGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameC, namespace, sm.PortSvcC)
}

func (sm Linkerd) SvcAGetPort() string {
	return sm.PortSvcA
}

func (sm Linkerd) SvcBGetPort() string {
	return sm.PortSvcB
}

func (sm Linkerd) SvcCGetPort() string {
	return sm.PortSvcC
}

type Istio struct {
	PortSvcA string
	PortSvcB string
	PortSvcC string
}

func (sm Istio) SvcAGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameA, namespace, sm.PortSvcA)
}

func (sm Istio) SvcBGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameB, namespace, sm.PortSvcB)
}

func (sm Istio) SvcCGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameC, namespace, sm.PortSvcC)
}

func (sm Istio) SvcAGetPort() string {
	return sm.PortSvcA
}

func (sm Istio) SvcBGetPort() string {
	return sm.PortSvcB
}

func (sm Istio) SvcCGetPort() string {
	return sm.PortSvcC
}

type OSM struct {
	PortSvcA string
	PortSvcB string
	PortSvcC string
}

func (sm OSM) SvcAGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameA, namespace, sm.PortSvcA)
}

func (sm OSM) SvcBGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameB, namespace, sm.PortSvcB)
}

func (sm OSM) SvcCGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s..svc.cluster.local.:%s", SvcNameC, namespace, sm.PortSvcC)
}

func (sm OSM) SvcAGetPort() string {
	return sm.PortSvcA
}

func (sm OSM) SvcBGetPort() string {
	return sm.PortSvcB
}

func (sm OSM) SvcCGetPort() string {
	return sm.PortSvcC
}
