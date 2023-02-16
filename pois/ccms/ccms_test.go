package ccms_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CCMSResponse struct {
	Message string `json:"message"`
}

var currentDate = time.Now()
var Date = currentDate.Format("01022006")

// default test Schedule Content
var DefaultTestScheduleContent = []string{"LOI 1201 001500 0003 0057 001 001 000030 000000 00000000 000 00000039108 0000 PAPA",
	"LOI 1201 001500 0003 0057 001 002 000030 000000 00000000 000 00000018709 0000 ESPN", "LOI 1201 003000 0003 0057 002 001 000030 000000 00000000 000 00000021902 0000 AIR",
	"LOI 1201 003000 0003 0057 002 002 000030 000000 00000000 000 00000018708 0000 ESPN", "LOI 1201 004500 0003 0057 003 001 000030 000000 00000000 000 00000039804 0000 CITY OF"}

func GetCurrentMonthAndDate() (string, string) {
	var monthStr, dayStr string
	currentDate := time.Now()
	monthInt := int(currentDate.Month())
	if monthInt < 10 {
		monthStr = "0" + strconv.Itoa(monthInt)
	} else {
		monthStr = strconv.Itoa(monthInt)
	}
	day := currentDate.Day()
	if day < 10 {
		dayStr = "0" + strconv.Itoa(day)
	} else {
		dayStr = strconv.Itoa(day)
	}
	return monthStr, dayStr
}

// Generate the schedule content for the current date
func GenerateScheduleContent() string {
	//get the current date
	var scheduleContent string
	month, day := GetCurrentMonthAndDate()
	for _, scheCont := range DefaultTestScheduleContent {
		scheContComponents := strings.Split(scheCont, " ")
		scheContComponents[1] = month + day
		scheduleContent += strings.Join(scheContComponents, " ") + "\n"
	}
	return scheduleContent
}

// create a test file for positive case// returns file name
func CreateScheduleFile() string {
	//Get the scheduleConten
	scheduleContent := GenerateScheduleContent()
	//create file with current day and create file with future date
	_, day := GetCurrentMonthAndDate()
	hexaMonth := strconv.FormatInt(int64(time.Now().Month()), 16)
	fileName := hexaMonth + day + "01001" + ".SCH"
	scheduleFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal("File creation failed")
	}
	defer scheduleFile.Close()
	//write the schedule file content to store into the file
	_, err = scheduleFile.WriteString(scheduleContent)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return fileName
}

func DeleteFile() {
	currentDiretory, _ := os.Getwd()
	files, err := ioutil.ReadDir(currentDiretory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".sch" {
			_, err := os.Stat(file.Name())
			if !errors.Is(err, os.ErrNotExist) {
				os.Remove(file.Name())
			}
		}
	}
}

// Common http request and response setup for alias testing.
func CCMsTestHttpRequest(method string, url string, data string) (int, []byte, error) {
	file, _ := os.ReadFile(data)
	reqBody := bytes.NewReader(file)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, nil, err
	}

	if method == "GET" {
		req.Header.Set("Accept", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
	}
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

// Test function for ccms Put request
// Below function will check for the response status code and response body
func TestCCMSPutRequest(t *testing.T) {
	//Positve Scenario testing
	requestUrl := "http://localhost:8130/pois/v1/channels/cnn/" + Date
	fmt.Println("requesturl ", requestUrl)
	responseMessage := "Schedule created for the channel cnn to the date:02062023 valid schedule count: 0 and invalid schedule count: 5"
	respStatus, respBody, err := CCMsTestHttpRequest("PUT", requestUrl, CreateScheduleFile())
	var posResult CCMSResponse
	json.Unmarshal(respBody, &posResult)
	assert.Equal(t, http.StatusCreated, respStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, responseMessage, posResult.Message)

	//Negative Scenario
	negRequestUrl := "http://localhost:8130/pois/v1/channels/cnn/01022023"
	negResponseMessage := "Please provide a valid date for the given channel name"
	negRespStatus, negRespBody, err := CCMsTestHttpRequest("PUT", negRequestUrl, "cnn2311.SCH")
	var negResult CCMSResponse
	json.Unmarshal(negRespBody, &negResult)
	assert.Equal(t, http.StatusBadRequest, negRespStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, negResponseMessage, negResult.Message)
}

// Test function for alias Put request
// Below function will check for the response status code
func TestCCMSGetRequest(t *testing.T) {

	//Positve Scenario testing
	requestUrl := "http://localhost:8130/pois/v1/channels/cnn/" + Date
	respStatus, _, err := CCMsTestHttpRequest("GET", requestUrl, "")
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, respStatus)

	//Negative scenario testing
	responseMessage := "No schedule found for the channel cnn in a date 30012023"
	negrequestUrl := "http://localhost:8130/pois/v1/channels/cnn/30012023"
	negrespStatus, respBody, err := CCMsTestHttpRequest("GET", negrequestUrl, "")
	var result CCMSResponse
	json.Unmarshal(respBody, &result)
	assert.Equal(t, http.StatusNotFound, negrespStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, responseMessage, result.Message)

}

// Test function for CCMS Delete request
// Below function will check for the response status code and response body
func TestCCMSDelRequest(t *testing.T) {
	// Positive Scenario
	requestUrl := "http://localhost:8130/pois/v1/channels/cnn/" + Date
	respMessage := "Schedule deleted for the channel cnn"
	respStatus, respBody, err := CCMsTestHttpRequest("DELETE", requestUrl, "")
	var posResult CCMSResponse
	json.Unmarshal(respBody, &posResult)
	assert.Equal(t, http.StatusOK, respStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, respMessage, posResult.Message)

	//Negative scenario testing
	responseMessage := "No schedule found for the given date 30012023"
	negrequestUrl := "http://localhost:8130/pois/v1/channels/cnn/30012023"
	negrespStatus, respBody, err := CCMsTestHttpRequest("DELETE", negrequestUrl, "")
	var result CCMSResponse
	json.Unmarshal(respBody, &result)
	assert.Equal(t, http.StatusNotFound, negrespStatus)
	assert.Equal(t, nil, err)
	assert.Equal(t, responseMessage, result.Message)
}
