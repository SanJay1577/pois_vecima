package prometheus

import (
	"encoding/json"
	"io/ioutil"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
	"github.com/prometheus/client_golang/prometheus"
)

// prometheus_metric struct which contains a request, redirect, default and failure routes
type Prometheus_metric struct {
	Value    int    `json:"value"`
	Code     string `json:"code"`
	Method   string `json:"method"`
	Protocol string `json:"protocol"`
}

// prometheus_metric struct which contains a request, redirect, default and failure routes
type Metric_data struct {
	Name string              `json:"name"`
	Data []Prometheus_metric `json:"data"`
}

// initialize gloabal prometheus struct
var metric []Metric_data
var prometheusFile string = "prometheus.json"

// Fetch the value from JSON file and update prometheus counter
func ReadPrometheusStatsFromJson(log *zaplogger.Logger) error {
	byteValue, err := ioutil.ReadFile(prometheusFile)
	if err != nil {
		log.Errorf("Error reading prometheus configuration file %v", err)
		return err
	}
	// unmarshal byteArray which contains jsonFile's content into struct
	// which we defined above
	err = json.Unmarshal(byteValue, &metric)
	if err != nil {
		log.Errorf("Error unmarshalling prometheus configuration file %v", err)
		return err
	}
	for i := 0; i < len(metric); i++ {
		UpdatePrometheusCounter(metric[i].Name, metric[i].Data)
	}
	return err
}

// Update prometheus counter for respective counter name
func UpdatePrometheusCounter(name string, data []Prometheus_metric) {
	for i := 0; i < len(data); i++ {
		if name == "ccms_response_classes_total" {
			CcmsRequest.WithLabelValues(data[i].Code, data[i].Method, data[i].Protocol).Add(float64(data[i].Value))
		} else if name == "ccms_total_request_served" {
			CcmsDefaultRequest.WithLabelValues().Add(float64(data[i].Value))
		} else if name == "alias_response_classes_total" {
			AliasRequest.WithLabelValues(data[i].Code, data[i].Method, data[i].Protocol).Add(float64(data[i].Value))
		} else if name == "alias_total_request_served" {
			AliasDefaultRequest.WithLabelValues().Add(float64(data[i].Value))
		} else if name == "esam_success_response_total" {
			EsamSucessRequest.WithLabelValues().Add(float64(data[i].Value))
		} else if name == "esam_noop_response_total" {
			EsamNoopRequest.WithLabelValues().Add(float64(data[i].Value))
		} else if name == "esam_delete_response_total" {
			EsamDeleteRequest.WithLabelValues().Add(float64(data[i].Value))
		} else if name == "esam_total_request_served" {
			EsamDefaultRequest.WithLabelValues().Add(float64(data[i].Value))
		}
	}
}

// Collect the prometheus counter value and persist the data by metric data struct format
func FormPrometheusData(name string, data []Prometheus_metric) Metric_data {
	var prometheusMetric []Prometheus_metric
	for i := 0; i < len(data); i++ {
		var prometheusStruct Prometheus_metric
		if name == "ccms_response_classes_total" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(CcmsRequest, prometheus.Labels{"code": data[i].Code, "method": data[i].Method, "protocol": data[i].Protocol})), Code: data[i].Code, Method: data[i].Method, Protocol: data[i].Protocol}
		} else if name == "ccms_total_request_served" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(CcmsDefaultRequest, prometheus.Labels{}))}
		} else if name == "alias_response_classes_total" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(AliasRequest, prometheus.Labels{"code": data[i].Code, "method": data[i].Method, "protocol": data[i].Protocol})), Code: data[i].Code, Method: data[i].Method, Protocol: data[i].Protocol}
		} else if name == "alias_total_request_served" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(AliasDefaultRequest, prometheus.Labels{}))}
		} else if name == "esam_success_response_total" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(EsamSucessRequest, prometheus.Labels{}))}
		} else if name == "esam_noop_response_total" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(EsamNoopRequest, prometheus.Labels{}))}
		} else if name == "esam_delete_response_total" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(EsamDeleteRequest, prometheus.Labels{}))}
		} else if name == "esam_total_request_served" {
			prometheusStruct = Prometheus_metric{Value: int(GetCounterValue(EsamDefaultRequest, prometheus.Labels{}))}
		}
		prometheusMetric = append(prometheusMetric, prometheusStruct)
	}
	metricData := Metric_data{Name: name, Data: prometheusMetric}
	return metricData
}

// Write prometheus metric in the file to persist the data
func WritePrometheusToJsonFile(log *zaplogger.Logger) error {
	var prometheusFiedlData []Metric_data
	for i := 0; i < len(metric); i++ {
		data := FormPrometheusData(metric[i].Name, metric[i].Data)
		prometheusFiedlData = append(prometheusFiedlData, data)
	}
	content, err := json.Marshal(prometheusFiedlData)
	if err != nil {
		log.Errorf("Error marshalling prometheus data %v", err)
		return err
	}
	err = ioutil.WriteFile(prometheusFile, content, 0644)
	if err != nil {
		log.Errorf("Error Writing prometheus stats to file %v", err)
		return err
	}
	return err
}
