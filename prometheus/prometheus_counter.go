package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// Declaraion and initilaization of Prometheus counter with respective request names and values
var (
	CcmsRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ccms_response_total",
			Help: "Response count by methods, i.e. 2xx, 4xx, 5xx",
		},
		[]string{"code", "method", "protocol"},
	)

	CcmsDefaultRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ccms_total_request_served",
			Help: "Total number of request served by ccms",
		},
		[]string{},
	)

	AliasRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alias_response_total",
			Help: "Response count by methods, i.e. 2xx, 4xx, 5xx",
		},
		[]string{"code", "method", "protocol"},
	)

	AliasDefaultRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alias_total_request_served",
			Help: "Total number of request served by alias",
		},
		[]string{},
	)

	EsamSucessRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_success_response_total",
			Help: "Response count esam served",
		},
		[]string{},
	)

	EsamNoopRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_noop_response_total",
			Help: "Response count esam served",
		},
		[]string{},
	)

	EsamDeleteRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_delete_response_total",
			Help: "Response count esam served",
		},
		[]string{},
	)

	EsamDefaultRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_total_request_served",
			Help: "Total number of request served by esam",
		},
		[]string{},
	)
)

// Getting the Counter Value for Required Prometheus metric with respective lables
func GetCounterValue(metric *prometheus.CounterVec, label prometheus.Labels) float64 {

	var m = &dto.Metric{}
	if err := metric.With(label).Write(m); err != nil {
		return 0
	}

	return m.Counter.GetValue()
}
