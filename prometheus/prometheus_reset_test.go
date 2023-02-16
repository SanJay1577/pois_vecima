package prometheus_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Struct contains a test request, redirect, default and failure routes
type Test_Prometheus_metric struct {
	Value    int    `json:"value"`
	Code     string `json:"code,omitempty"`
	Method   string `json:"method,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// Struct contains collection of test request router metrics
type Test_Metric_data struct {
	Name string                   `json:"name"`
	Data []Test_Prometheus_metric `json:"data"`
}

// Initialize gloabal struct
var Test_metric []Test_Metric_data

func TestReadFromJsonFile(t *testing.T) {
	//reading the json file.
	err := ReadJsonFile()
	if err != nil {
		t.Logf("Error while read the value from JSON file and unmarshal it to go object %s\n", err)
	}
	lengthCheck := len(Test_metric)
	// Checking  the return object have a data in it
	assert.Greater(t, lengthCheck, 0)
}

// Read the value from JSON file and unmarshal it to go object
func ReadJsonFile() error {
	var err error
	err = nil
	byteValue, err := ioutil.ReadFile("/root/workspace/pois/prometheus.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &Test_metric)
	if err != nil {
		return err
	}

	return err
}

// Testing Data Marshalling in Prometheus
func TestJsonMarshal(t *testing.T) {
	got := demoMarshallContent()
	want := "string"
	assert.IsType(t, got, want)
}

// Test the unmarshalled go object having data
func demoMarshallContent() string {
	ReadJsonFile()
	stringValue, _ := json.Marshal(Test_metric)
	stringValues := string(stringValue)

	return stringValues
}
