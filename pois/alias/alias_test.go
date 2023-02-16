package alias_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AliasResponse struct {
	Message string `json:"message"`
}

// Common http request and response setup for alias testing.
func AliasTestHttpRequest(method string, url string, data string) (int, []byte, error) {
	reqBody := bytes.NewReader([]byte(data))
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, respBody, nil
}

// Test function for alias Put request
// Below function will check for the response status code and response body
func TestAliasPutRequest(t *testing.T) {
	// Positive scenario testing
	requestUrl := "http://localhost:8130/pois/v1/channels/alias/cnn"
	requestBody := `{
    "aliasNames" : [
        "cnnlive4",
        "cnnlive5"
    ]
}`
	responseMessage := "Alias names mapped for the channel cnn"
	respStatus, respBody, err := AliasTestHttpRequest("PUT", requestUrl, requestBody)
	var result AliasResponse
	json.Unmarshal(respBody, &result)
	assert.Equal(t, http.StatusCreated, respStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, responseMessage, result.Message)

	//Negative scenarion testing
	negRequestUrl := "http://localhost:8130/pois/v1/channels/alias/cnn/saka"
	negRequestBody := `{
    "aliasNames" : [
        "cnnlive4",
        "cnnlive5"
    ]
}`
	negResponseMessage := "Please verify the request and try again"
	negRespStatus, negRespBody, err := AliasTestHttpRequest("PUT", negRequestUrl, negRequestBody)
	var negResult AliasResponse
	json.Unmarshal(negRespBody, &negResult)
	assert.Equal(t, http.StatusBadRequest, negRespStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, negResponseMessage, negResult.Message)

}

// Test function for alias GET request
// Below function will check for the response status code
func TestAliasGetRequest(t *testing.T) {
	//Positve Scenario testing
	requestUrl := "http://localhost:8130/pois/v1/channels/alias/cnn"
	respStatus, _, err := AliasTestHttpRequest("GET", requestUrl, "")
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, nil, err)

	//Negative scenario testing
	responseMessage := " No alias found for the channel c"
	negrequestUrl := "http://localhost:8130/pois/v1/channels/alias/c"
	negrespStatus, respBody, err := AliasTestHttpRequest("GET", negrequestUrl, "")
	var result AliasResponse
	json.Unmarshal(respBody, &result)
	assert.Equal(t, http.StatusNotFound, negrespStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, responseMessage, result.Message)

}

// Test function for alias Delete request
// Below function will check for the response status code and response body
func TestAliasDeleteRequest(t *testing.T) {

	// Positive scenario Testing

	requestUrl := "http://localhost:8130/pois/v1/channels/alias/cnn/cnnlive5"

	responseMessage := "cnnlive5 alias name deleted for the channel cnn"
	respStatus, respBody, err := AliasTestHttpRequest("DELETE", requestUrl, "")
	var result AliasResponse
	json.Unmarshal(respBody, &result)
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, responseMessage, result.Message)

	//Negative scenarion testing
	negRequestUrl := "http://localhost:8130/pois/v1/channels/alias/cnn/cnnlive"
	negResponseMessage := "No alias found for the  channel  cnn"
	negRespStatus, negRespBody, err := AliasTestHttpRequest("DELETE", negRequestUrl, "")
	var negResult AliasResponse
	json.Unmarshal(negRespBody, &negResult)
	assert.Equal(t, http.StatusBadRequest, negRespStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, negResponseMessage, negResult.Message)

}
