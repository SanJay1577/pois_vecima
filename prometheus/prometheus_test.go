package prometheus_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/stretchr/testify/assert"

	dto "github.com/prometheus/client_model/go"
)

// Prometheus counter for test Testrequest
var (
	CcmsTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ccms_response_classes_total",
			Help: "Response count by class, i.e. 2xx, 4xx, 5xx",
		},
		[]string{"code", "method", "protocol"},
	)

	CcmsDefaultTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ccms_total_Testrequest_served",
			Help: "Total number of Testrequest served by ccms",
		},
		[]string{},
	)

	AliasTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alias_response_classes_total",
			Help: "Response count by class, i.e. 2xx, 4xx, 5xx",
		},
		[]string{"code", "method", "protocol"},
	)

	AliasDefaultTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alias_total_Testrequest_served",
			Help: "Total number of Testrequest served by alias",
		},
		[]string{},
	)

	EsamSucessTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_success_response_total",
			Help: "Response count esam served",
		},
		[]string{},
	)

	EsamNoopTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_noop_response_total",
			Help: "Response count esam served",
		},
		[]string{},
	)

	EsamDeleteTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_delete_response_total",
			Help: "Response count esam served",
		},
		[]string{},
	)

	EsamDefaultTestRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "esam_total_Testrequest_served",
			Help: "Total number of Testrequest served by esam",
		},
		[]string{},
	)
)

// Perfom Test for prometheus counter for ccms
func TestCCMSPrometheusCounter(t *testing.T) {

	assert := assert.New(t)

	// Test counter have to be registered
	prometheus.MustRegister(CcmsTestRequest)

	// Increase prometheus counter for all labels
	CcmsTestRequest.WithLabelValues("2xx", "GET", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("2xx", "GET", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("3xx", "GET", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("4xx", "GET", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("5xx", "GET", "HTTP/1.1").Inc()

	CcmsTestRequest.WithLabelValues("2xx", "PUT", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("2xx", "PUT", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("3xx", "PUT", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("4xx", "PUT", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("5xx", "PUT", "HTTP/1.1").Inc()

	CcmsTestRequest.WithLabelValues("2xx", "DELETE", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("2xx", "DELETE", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("3xx", "DELETE", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("4xx", "DELETE", "HTTP/1.1").Inc()
	CcmsTestRequest.WithLabelValues("5xx", "DELETE", "HTTP/1.1").Inc()

	// Collected three metrics
	assert.Equal(12, testutil.CollectAndCount(CcmsTestRequest))

	// Check the expected values with actual values
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("3xx", "GET", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("5xx", "GET", "HTTP/1.1")))
	assert.Equal(float64(2), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("2xx", "GET", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("4xx", "GET", "HTTP/1.1")))

	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("4xx", "PUT", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("3xx", "PUT", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("5xx", "PUT", "HTTP/1.1")))
	assert.Equal(float64(2), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("2xx", "PUT", "HTTP/1.1")))

	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("4xx", "DELETE", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("3xx", "DELETE", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("5xx", "DELETE", "HTTP/1.1")))
	assert.Equal(float64(2), testutil.ToFloat64(CcmsTestRequest.WithLabelValues("2xx", "DELETE", "HTTP/1.1")))

	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "3xx", "method": "GET", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "5xx", "method": "GET", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(2), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "2xx", "method": "GET", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "4xx", "method": "GET", "protocol": "HTTP/1.1"}))

	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "3xx", "method": "PUT", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "5xx", "method": "PUT", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(2), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "2xx", "method": "PUT", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "4xx", "method": "PUT", "protocol": "HTTP/1.1"}))

	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "3xx", "method": "DELETE", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "5xx", "method": "DELETE", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "4xx", "method": "DELETE", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(2), GetCounterValues(CcmsTestRequest, prometheus.Labels{"code": "2xx", "method": "DELETE", "protocol": "HTTP/1.1"}))

}

// prometheus counter testing for alias
func TestAliasPrometheusCounter(t *testing.T) {

	assert := assert.New(t)

	// Test counter have to be registered
	prometheus.MustRegister(AliasTestRequest)

	// Increase prometheus counter for all labels
	AliasTestRequest.WithLabelValues("2xx", "GET", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("2xx", "GET", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("3xx", "GET", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("4xx", "GET", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("5xx", "GET", "HTTP/1.1").Inc()

	AliasTestRequest.WithLabelValues("2xx", "PUT", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("2xx", "PUT", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("3xx", "PUT", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("4xx", "PUT", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("5xx", "PUT", "HTTP/1.1").Inc()

	AliasTestRequest.WithLabelValues("2xx", "DELETE", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("2xx", "DELETE", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("3xx", "DELETE", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("4xx", "DELETE", "HTTP/1.1").Inc()
	AliasTestRequest.WithLabelValues("5xx", "DELETE", "HTTP/1.1").Inc()

	// Collected three metrics
	assert.Equal(12, testutil.CollectAndCount(AliasTestRequest))

	// Check the expected values with actual values
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("3xx", "GET", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("5xx", "GET", "HTTP/1.1")))
	assert.Equal(float64(2), testutil.ToFloat64(AliasTestRequest.WithLabelValues("2xx", "GET", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("4xx", "GET", "HTTP/1.1")))

	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("4xx", "PUT", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("3xx", "PUT", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("5xx", "PUT", "HTTP/1.1")))
	assert.Equal(float64(2), testutil.ToFloat64(AliasTestRequest.WithLabelValues("2xx", "PUT", "HTTP/1.1")))

	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("4xx", "DELETE", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("3xx", "DELETE", "HTTP/1.1")))
	assert.Equal(float64(1), testutil.ToFloat64(AliasTestRequest.WithLabelValues("5xx", "DELETE", "HTTP/1.1")))
	assert.Equal(float64(2), testutil.ToFloat64(AliasTestRequest.WithLabelValues("2xx", "DELETE", "HTTP/1.1")))

	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "3xx", "method": "GET", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "5xx", "method": "GET", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(2), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "2xx", "method": "GET", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "4xx", "method": "GET", "protocol": "HTTP/1.1"}))

	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "3xx", "method": "PUT", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "5xx", "method": "PUT", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(2), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "2xx", "method": "PUT", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "4xx", "method": "PUT", "protocol": "HTTP/1.1"}))

	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "3xx", "method": "DELETE", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "5xx", "method": "DELETE", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(1), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "4xx", "method": "DELETE", "protocol": "HTTP/1.1"}))
	assert.Equal(float64(2), GetCounterValues(AliasTestRequest, prometheus.Labels{"code": "2xx", "method": "DELETE", "protocol": "HTTP/1.1"}))

}

func TestCCmsDefaultCounter(t *testing.T) {

	assert := assert.New(t)

	// Test counter have to be registered
	prometheus.MustRegister(CcmsDefaultTestRequest)
	//incrementing the counter for default
	CcmsDefaultTestRequest.WithLabelValues().Inc()
	//asserting the counter value
	assert.Equal(1, testutil.CollectAndCount(CcmsDefaultTestRequest))
	assert.Equal(float64(1), testutil.ToFloat64(CcmsDefaultTestRequest.WithLabelValues()))
	assert.Equal(float64(1), GetCounterValues(CcmsDefaultTestRequest, prometheus.Labels{}))

}

// testing alias default request
func TestAliasDefaultCounter(t *testing.T) {
	assert := assert.New(t)
	// Test counter have to be registered
	prometheus.MustRegister(AliasDefaultTestRequest)
	//incrementing the counter for default
	AliasDefaultTestRequest.WithLabelValues().Inc()
	//asserting the counter value
	assert.Equal(1, testutil.CollectAndCount(AliasDefaultTestRequest))
	assert.Equal(float64(1), testutil.ToFloat64(AliasDefaultTestRequest.WithLabelValues()))
	assert.Equal(float64(1), GetCounterValues(AliasDefaultTestRequest, prometheus.Labels{}))

}

// testing esam default request

func TestEsamDefaultCounter(t *testing.T) {
	assert := assert.New(t)
	// Test counter have to be registered
	prometheus.MustRegister(EsamDefaultTestRequest)
	//incrementing the counter for default
	EsamDefaultTestRequest.WithLabelValues().Inc()
	//asserting the counter value
	assert.Equal(1, testutil.CollectAndCount(EsamDefaultTestRequest))
	assert.Equal(float64(1), testutil.ToFloat64(EsamDefaultTestRequest.WithLabelValues()))
	assert.Equal(float64(1), GetCounterValues(EsamDefaultTestRequest, prometheus.Labels{}))

}

// Testing Esam Success Response counter
func TestEsamSuccessCounter(t *testing.T) {

	assert := assert.New(t)
	// Test counter have to be registered
	prometheus.MustRegister(EsamSucessTestRequest)
	//incrementing the counter for Success
	EsamSucessTestRequest.WithLabelValues().Inc()
	//asserting the counter value
	assert.Equal(1, testutil.CollectAndCount(EsamSucessTestRequest))
	assert.Equal(float64(1), testutil.ToFloat64(EsamSucessTestRequest.WithLabelValues()))
	assert.Equal(float64(1), GetCounterValues(EsamSucessTestRequest, prometheus.Labels{}))

}

// Testing Esam Noop Response counter

func TestEsamNoopCounter(t *testing.T) {
	assert := assert.New(t)
	// Test counter have to be registered
	prometheus.MustRegister(EsamNoopTestRequest)
	//incrementing the counter for Success
	EsamNoopTestRequest.WithLabelValues().Inc()
	//asserting the counter value
	assert.Equal(1, testutil.CollectAndCount(EsamNoopTestRequest))
	assert.Equal(float64(1), testutil.ToFloat64(EsamNoopTestRequest.WithLabelValues()))
	assert.Equal(float64(1), GetCounterValues(EsamNoopTestRequest, prometheus.Labels{}))

}

// Testing Esam Delete Response counter

func TestEsamDeleteCounter(t *testing.T) {
	assert := assert.New(t)
	// Test counter have to be registered
	prometheus.MustRegister(EsamDeleteTestRequest)
	//incrementing the counter for Success
	EsamDeleteTestRequest.WithLabelValues().Inc()
	//asserting the counter value
	assert.Equal(1, testutil.CollectAndCount(EsamDeleteTestRequest))
	assert.Equal(float64(1), testutil.ToFloat64(EsamDeleteTestRequest.WithLabelValues()))
	assert.Equal(float64(1), GetCounterValues(EsamDeleteTestRequest, prometheus.Labels{}))

}

// Get prometheus counter value by labels
func GetCounterValues(metric *prometheus.CounterVec, label prometheus.Labels) float64 {

	var m = &dto.Metric{}
	if err := metric.With(label).Write(m); err != nil {
		return 0
	}

	return m.Counter.GetValue()
}
