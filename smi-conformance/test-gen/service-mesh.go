package test_gen

import "fmt"

type ServiceMesh interface {
	SvcAGetInternalName() string
	SvcBGetInternalName() string
	SvcCGetInternalName() string

	SvcAGetExternalName() string
	SvcBGetExternalName() string
	SvcCGetExternalName() string
}

type Maesh struct {
	Namespace string
	PortSvcA  string
	PortSvcB  string
	PortSvcC  string
}

func (sm Maesh) SvcAGetInternalName() string {
	return fmt.Sprintf("http://%s.%s.maesh:%s/", SERVICE_A_NAME, sm.Namespace, sm.PortSvcA)
}

func (sm Maesh) SvcBGetInternalName() string {
	return fmt.Sprintf("http://%s.%s.maesh:%s/", SERVICE_B_NAME, sm.Namespace, sm.PortSvcB)
}

func (sm Maesh) SvcCGetInternalName() string {
	return fmt.Sprintf("http://%s.%s.maesh:%s/", SERVICE_C_NAME, sm.Namespace, sm.PortSvcC)
}

func (sm Maesh) SvcAGetExternalName() string {
	return fmt.Sprintf("http://%s:%s/", SERVICE_A_NAME, sm.PortSvcA)
}

func (sm Maesh) SvcBGetExternalName() string {
	return fmt.Sprintf("http://%s:%s/", SERVICE_B_NAME, sm.PortSvcB)
}

func (sm Maesh) SvcCGetExternalName() string {
	return fmt.Sprintf("http://%s:%s/", SERVICE_C_NAME, sm.PortSvcC)
}
