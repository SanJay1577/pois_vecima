package esam_test

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	esam "pois/pois/esam"
	"testing"

	"github.com/stretchr/testify/assert"
)

var UrlPath = "http://localhost:4056/esam/v1/comcast/request"

// Common http request and response setup for alias testing.
func EsamTestHttpRequest(method string, url string, data string) (int, []byte, error) {
	reqBody := bytes.NewReader([]byte(data))
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
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

// esam Replace Request testing

func TestXmlReplace(t *testing.T) {
	requestBody := (`
	<SignalProcessingEvent >
	<AcquiredSignal acquisitionPointIdentity="cnn" acquisitionSignalID="3d679213-83e2-4899-8a64-de761b6ab781">
	<UTCPoint utcPoint="2023-02-02T00:45:00Z" />
	<BinaryData signalType="SCTE35" >/DAhAAAAAAAA///wEAUAADaUf29/fgApMuAAAQAAAAA2xmw5</BinaryData>
	<StreamTimes> 
	<StreamTime timeType="PTS" timeValue="6164518823"/>
	</StreamTimes>
	</AcquiredSignal>
	</SignalProcessingEvent>
`)
	respStatus, respBody, err := EsamTestHttpRequest("POST", UrlPath, requestBody)
	var responseAction esam.SignalProcessNotification
	xml.Unmarshal(respBody, &responseAction)
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, "replace", responseAction.Responsesignal.Action)
	assert.Equal(t, nil, err)

}

// esam Delete Request testing
func TestXmlDelete(t *testing.T) {
	requestBody := (`
	<SignalProcessingEvent >
	<AcquiredSignal acquisitionPointIdentity="cnn" acquisitionSignalID="3d679213-83e2-4899-8a64-de761b6ab781">
	<UTCPoint utcPoint="2023-01-02T00:45:00Z" />
	<BinaryData signalType="SCTE35" >/DAhAAAAAAAA///wEAUAADaUf29/fgApMuAAAQAAAAA2xmw5</BinaryData>
	<StreamTimes> 
	<StreamTime timeType="PTS" timeValue="6164518823"/>
	</StreamTimes>
	</AcquiredSignal>
	</SignalProcessingEvent>
`)
	respStatus, respBody, err := EsamTestHttpRequest("POST", UrlPath, requestBody)
	var responseAction esam.SignalProcessNotification
	xml.Unmarshal(respBody, &responseAction)
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, "delete", responseAction.Responsesignal.Action)
	assert.Equal(t, nil, err)
}

// esam Noob Request testing
func TestXmlNoob(t *testing.T) {
	requestBody := (`
	<SignalProcessingEvent >
	<AcquiredSignal acquisitionPointIdentity="" acquisitionSignalID="3d679213-83e2-4899-8a64-de761b6ab781">
	<UTCPoint utcPoint="2023-02-02T00:45:00Z" />
	<BinaryData signalType="SCTE35" >/DAhAAAAAAAA///wEAUAADaUf29/fgApMuAAAQAAAAA2xmw5</BinaryData>
	<StreamTimes> 
	<StreamTime timeType="PTS" timeValue="6164518823"/>
	</StreamTimes>
	</AcquiredSignal>
	</SignalProcessingEvent>
`)
	respStatus, respBody, err := EsamTestHttpRequest("POST", UrlPath, requestBody)
	var responseAction esam.SignalProcessNotification
	xml.Unmarshal(respBody, &responseAction)
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, "noop", responseAction.Responsesignal.Action)
	assert.Equal(t, nil, err)
}

func TestSceteParser(t *testing.T) {
	requestBody := (`
	<SignalProcessingEvent >
	<AcquiredSignal acquisitionPointIdentity="cnn" acquisitionSignalID="3d679213-83e2-4899-8a64-de761b6ab781">
	<UTCPoint utcPoint="2023-01-02T00:45:00Z" />
	<BinaryData signalType="SCTE35" >/DAhAAAAAAAA///wEAUAADaUf29/fgApMuAAAQAAAAA2xm</BinaryData>
	<StreamTimes> 
	<StreamTime timeType="PTS" timeValue="6164518823"/>
	</StreamTimes>
	</AcquiredSignal>
	</SignalProcessingEvent>
`)
	respStatus, respBody, err := EsamTestHttpRequest("POST", UrlPath, requestBody)
	var responseAction esam.SignalProcessNotification
	xml.Unmarshal(respBody, &responseAction)
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, "delete", responseAction.Responsesignal.Action)
	assert.Equal(t, nil, err)
}
