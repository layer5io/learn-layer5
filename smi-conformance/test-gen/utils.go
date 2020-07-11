package test_gen

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