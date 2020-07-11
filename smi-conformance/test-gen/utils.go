package test_gen

const SERVICE_A_NAME = "app-a"
const SERVICE_B_NAME = "app-b"
const SERVICE_C_NAME = "app-c"

type URLstruct struct {
	URL     string
	Method  string
	Headers map[string]string
}
type MetricResponse struct {
	ReqReceived   []string
	RespSucceeded []URLstruct
	RespFailed    []URLstruct
}
