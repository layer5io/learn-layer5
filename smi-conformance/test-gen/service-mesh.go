package test_gen

import "fmt"

type ServiceMesh interface {
	SvcAGetInternalName(string) string
	SvcBGetInternalName(string) string
	SvcCGetInternalName(string) string

	SvcAGetPort() string
	SvcBGetPort() string
	SvcCGetPort() string
}

type Maesh struct {
	PortSvcA  string
	PortSvcB  string
	PortSvcC  string
}

func (sm Maesh) SvcAGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s.maesh:%s", SERVICE_A_NAME, namespace, sm.PortSvcA)
}

func (sm Maesh) SvcBGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s.maesh:%s", SERVICE_B_NAME, namespace, sm.PortSvcB)
}

func (sm Maesh) SvcCGetInternalName(namespace string) string {
	return fmt.Sprintf("http://%s.%s.maesh:%s", SERVICE_C_NAME, namespace, sm.PortSvcC)
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
