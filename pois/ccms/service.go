package ccms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	config "pois/config"
	serviceLog "pois/config/logging"
	core "pois/pois/core"
	"pois/prometheus"
	"pois/shared"
	sh "pois/shared"
	"strconv"
	"strings"

	"git.eng.vecima.com/cloud/golib/v4/httpservice"
)

// Get the application config file from the config file
// var applicationConfig = config.GetApplicationConfig()

// WebService API resource path.
const (
	CCMSMIN = iota
	CCMSVER1
	CCMSMAX
)

// Now define the RESTfulWebService API.
// For restful api services we the interface of service api which holds the (GET,PUT,POST,Delete) Methods
type serviceAPI struct{}

// swagger:route GET /{channelname}/{date} CCMSAPI GetChannelSchedule
//
// # Retrieve Channel Schedule Information
//
// This API retrieves the channel schedule information for the provided date from the database. It accepts the channel name and date from the path parameters to retrieve the channel scheduler information
//
// Produces:
// - application/json
// Parameters:
//   - name: channelname
//     in: path
//     description: Channel name to retrieve schedule information (channel accepts alphanumeric and  special characters (.-|!@#$%&_><+=))
//     type: string
//     required: true
//   - name: date
//     in: path
//     description: Retrieve schedule information for the date(DDMMYYYY)
//     type: string
//     required: true
//
// Responses:
//
//	200:scheduleRetrivalSuccessResponse (Schedule found for the channel for the given date {{DDMMYYYY}})
//	400:BadRequestResponse (Channel name or date not provided )
//	404:ChannelNotFoundResponse (No schedule found for the channel for the date {{DDMMYYYY}})
//
// Get handles a GET request.
// It parse's the request URL, grap the channelname and date from the request and it will validate and check
// the schedule for corresponding channel and date is available in the database and
// returns the schedule if it is present or else returns a 404 status as Not found.
func (h *serviceAPI) Get(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {

	var response *httpservice.HttpServiceResponse
	acceptTypeJson, acceptTypePlainText, acceptTypeDefault := false, false, true
	serviceLog.CcmsLog.Infof("[%v] TID:[%v] Path : [%v]", request.Method, transactionID, request.URL.Path)
	pathVariables, versionNumber, err := shared.ExtractPathVariable(request.URL.Path, config.GetConfig().GetString("api.rootpath"), CCMSMIN, CCMSMAX, config.GetConfig().GetString("api.channelpath"))
	if err != nil {
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Unsupported version %d [%d=%d]", versionNumber, CCMSMIN, CCMSMAX)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	if len(pathVariables) != 2 {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Channel or date not found", request.Method, transactionID)

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please verify the request and try again"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	channelName := pathVariables[0]
	date := pathVariables[1]
	//name and date validation
	channelNameIsValid := shared.ValidateChannelName(channelName)

	if channelNameIsValid {
		channelName, channelNameIsValid = core.GetChannelNameByAliasName(channelName)
	}

	dayVal, monthVal, yearVal, dateIsValid := shared.ValidateDate(date)
	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] Channelname:[%v] and Date:[%v] and year:[%v]", request.Method, transactionID, channelName, date, yearVal)

	if !channelNameIsValid {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid channelname:[%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide a valid channel name"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	if !dateIsValid {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid Date:[%v]", request.Method, transactionID, date)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the valid date"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	serviceLog.CcmsLog.Infof("[%v] TID:[%v] ChannelName:[%v] Date:[%v]", request.Method, transactionID, channelName, date)
	// Requested the Schedule in the file the Header content-type as text/plain
	_, acceptTypeExist := request.Header["Accept"]
	if acceptTypeExist {
		//check the Accept type as text/plain or application/json
		// The Initated request  Accept header is text/plain give response as the scheduled file
		// The Initiated request Accept header is application/json  give response as preprocessed scheduled file.
		//the  request.Header is a slice of string
		requestHeaderSlice := request.Header["Accept"]
		for _, contentTypeValue := range requestHeaderSlice {
			if strings.ToLower(contentTypeValue) == "text/plain" {
				acceptTypePlainText = true
			}
			if strings.ToLower(contentTypeValue) == "application/json" {
				acceptTypeJson = true
			}

		}
		//return the file response
		if acceptTypePlainText {
			ch := sh.ChannelSchedules{Channel: channelName, Date: monthVal + dayVal}

			if schedules, ok := sh.ChannelScheduleData.GetSchedule(ch); ok {
				scheduleLine := "REM Created On " + monthVal + "/" + dayVal + "/" + yearVal + " " + channelName + "\n"
				scheduleLine = scheduleLine + sh.ConcatinateScheduleInformation(schedules)
				scheduleLine = scheduleLine + "END"
				return requestResponse(request, []byte(scheduleLine), http.StatusOK, "2xx"), nil
			} else {
				responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("No Schedule found for the channel %v in the given date %v", channelName, date)})
				if err != nil {
					return internalServerErrorResponse(request, transactionID), nil
				}
				return requestResponse(request, responseMessage, http.StatusNotFound, "4xx"), nil
			}
		}

		if acceptTypeJson || acceptTypeDefault {
			ch := sh.ChannelSchedules{Channel: channelName, Date: monthVal + dayVal}
			if schedules, ok := sh.ChannelScheduleData.GetSchedule(ch); ok {
				jsonStr, err := json.Marshal(schedules)
				if err != nil {
					//fmt.Printf("Error: %s", err.Error())
					serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Schedule JSON is empty  for the channel [%v]", request.Method, transactionID, channelName)
					responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("No data for the channel %v in a given date %v", channelName, date)})
					if err != nil {
						return internalServerErrorResponse(request, transactionID), nil
					}
					return requestResponse(request, responseMessage, http.StatusNoContent, "4xx"), nil

				} else {

					serviceLog.CcmsLog.Infof("[%v] TID:[%v] Schedule found for channel [%v]", request.Method, transactionID, channelName)
					response = responseFormat(http.StatusOK, nil, jsonStr, "application/json")
					prometheus.CcmsDefaultRequest.WithLabelValues().Inc()
					prometheus.CcmsRequest.WithLabelValues("2xx", request.Method, request.Proto).Inc()
					return response, nil
				}
			} else {

				serviceLog.CcmsLog.Errorf("[%v] TID:[%v] No schedule found for the channel [%v]", request.Method, transactionID, channelName)
				responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("No schedule found for the channel %v in a date %v", channelName, date)})
				if err != nil {
					return internalServerErrorResponse(request, transactionID), nil
				}

				return requestResponse(request, responseMessage, http.StatusNotFound, "4xx"), nil
			}
		}

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Content-Type not Accepted", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the content-type as application/json or text/plain"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusNotAcceptable, "4xx"), nil

	} else {
		serviceLog.CcmsLog.Debugf("[%v] TID:[%v] Accept Type not Provided", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the Accept-Type Header"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusNotAcceptable, "4xx"), nil
	}
}

// store the schedular file in .SCH format
// Post handles a POST request.
func (h *serviceAPI) Post(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {
	response := &httpservice.HttpServiceResponse{
		Status:      http.StatusNotImplemented,
		Header:      nil,
		ContentType: "application/json",
		Body:        []byte("POST API not Implemented"),
	}

	return response, nil
}

// swagger:route PUT /{channelname}/{date} CCMSAPI
//
// # Add/Update Channel Schedule Information
//
// This API is responsible for adding or updating the channel schedule information. It accepts the channel name and date from the path parameters and updates the payload schedule data  into local database. This API's payload must be a binary file with Content-Type  as application/octet-stream
//
// Consumes:
// -  application/json
// Produces:
// - application/json
// Parameters:
//   - name: channelname
//     in: path
//     description: Channel name to add schedule information (channel name  accept  alphanumeric and  special characters (.-|!@#$%&_><+=))
//     type: string
//     required: true
//   - name: date
//     in: path
//     description: Add schedule information for the date(DDMMYYYY)
//     type: string
//     required: true
//
// Responses:
//
//		201:ScheduleCreated (Schedule created for the channel {{channel name}})
//		400:BadRequestResponse (Channel name, date parsing error, or request body parsing error)
//	    500:preprocessingFails (unable to process the file )
//
// Put handles a PUT request.
// create and update the schedules for corresponding channel with a given date
func (h *serviceAPI) Put(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {

	var response *httpservice.HttpServiceResponse
	var contentTypeOctetStream bool
	contentTypeOctetStream = false

	serviceLog.CcmsLog.Infof("[%v] TID:[%v] Path:[%v]", request.Method, transactionID, request.URL.Path)
	//Extract the  path params returns the params in slice
	//Example : /pois/v1/channels/cnn/06122022 -> returns the params [cnn 06122022]
	pathVariables, versionNumber, err := shared.ExtractPathVariable(request.URL.Path, config.GetConfig().GetString("api.rootpath"), CCMSMIN, CCMSMAX, config.GetConfig().GetString("api.channelpath"))

	if err != nil {
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Unsupported version %d [%d=%d]", versionNumber, CCMSMIN, CCMSMAX)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	if len(pathVariables) != 2 {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Channel or date not found", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please verify the request and try again"})
		if err != nil {

			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	// validate channel name
	channelName := pathVariables[0]
	channelNameIsValid := shared.ValidateChannelName(channelName)

	if channelNameIsValid {
		channelName, channelNameIsValid = core.GetChannelNameByAliasName(channelName)
	}

	if !channelNameIsValid {

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v]  Invalid Channelname:[%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide a valid channel name"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	// validate date
	date := pathVariables[1]
	dayVal, monthVal, yearVal, dateIsValid := shared.ValidateDate(date)
	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] Day:[%v] Month:[%v] Year:[%v]", request.Method, transactionID, dayVal, monthVal, yearVal)
	if !dateIsValid {

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid date:[%v]", request.Method, transactionID, date)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the valid date"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] Channelname:[%v] and Date:[%v]", request.Method, transactionID, channelName, date)

	//this function won't allow to create a schedules for previous dates
	schduledDateValid := shared.CompareDate(dayVal, monthVal, yearVal)

	if !(schduledDateValid) {

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid date:[%v]", request.Method, transactionID, date)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide a valid date for the given channel name"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	//Header type validation
	_, contentTypeExists := request.Header["Content-Type"]
	if !contentTypeExists {

		serviceLog.CcmsLog.Debugf("[%v] TID:[%v] Content-Type not Provided", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the content type"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusNotAcceptable, "4xx"), nil
	}

	requestHeaderSlice := request.Header["Content-Type"]
	for _, contentTypeValue := range requestHeaderSlice {
		if strings.ToLower(contentTypeValue) == "application/octet-stream" {
			contentTypeOctetStream = true
		}
	}

	if !contentTypeOctetStream {

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the content-type as application/octet-stream"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusNotAcceptable, "4xx"), nil
	}

	//parse the request body, processing the scheduler file
	// store the channels schedules in memory and file
	if request.Body != nil && request.ContentLength != 0 {
		body, err := ioutil.ReadAll(request.Body)

		if err != nil {

			serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Request body parsing:[%v]", request.Method, transactionID, request.Body)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Request Body parsing error: %v", err)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
		}

		defer request.Body.Close()
		isPreprocess, validScheduleCount, invalidSchdeuleCount := shared.PreprocessSchedulerFile(string(body), channelName, monthVal+dayVal)

		// file preprocessing error response
		if !(isPreprocess) {

			serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Schedular file preprocessing error", request.Method, transactionID)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Unable to preprocess file"})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			prometheus.CcmsDefaultRequest.WithLabelValues().Inc()
			prometheus.CcmsRequest.WithLabelValues("5xx", request.Method, request.Proto).Inc()
			response = responseFormat(http.StatusInternalServerError, nil, responseMessage, "application/json")
			return response, nil
		}

		serviceLog.CcmsLog.Infof("[%v] TID:[%v] Schedule created for the channel [%v] in a date [%v] valid schedule count [%v] invalid schedule count [%v]", request.Method, transactionID, channelName, date, validScheduleCount, invalidSchdeuleCount)
		defer serviceLog.CcmsLog.Sync()

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Schedule created for the channel %s to the date:%s valid schedule count: %v and invalid schedule count: %v", channelName, date, validScheduleCount, invalidSchdeuleCount)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusCreated, "2xx"), nil
	} else {

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Request body:[%v]", request.Method, transactionID, request.Body)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Request body is required"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil

	}
}

// swagger:route DELETE /{channelname}/{date} CCMSAPI RemoveScheduleForTheChannel
//
// # Delete Channel Schedule Information
//
// This API deletes the channel schedule information for the provided date from the local database. It accepts the channel name and date from the path parameters to delete the specified channel schedule information
//
// Produces:
// - application/json
// Parameters:
//   - name: channelname
//     in: path
//     description: Channel name to delete schedule information (channel name accepts alphanumeric and special characters (.-|!@#$%&_><+=))
//     type: string
//     required: true
//   - name: date
//     in: path
//     description: Delete schedule information for the date(DDMMYYYY)
//     type: string
//     required: true
//
// Responses:
//
//	200:DeleteSchdeuleResponse (Schedule deleted for the  given channel and  provided date {{DDMMYYYY}})
//	400:BadRequestResponse (Channel or date  not provided)
//	404:ChannelNotFoundResponse (No Schedule found for the channel for the   date {{DDMMYYYY}})
//
// Parse the request URL, grap the channelname and date from the request.
// this will check the schedule for corresponding channel in a provided date and  delete the schedule if it exists.
func (h *serviceAPI) Delete(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {

	serviceLog.CcmsLog.Infof("[%v] TID:[%v] Path:[%v]", request.Method, transactionID, request.URL.Path)
	pathVariables, versionNumber, err := shared.ExtractPathVariable(request.URL.Path, config.GetConfig().GetString("api.rootpath"), CCMSMIN, CCMSMAX, config.GetConfig().GetString("api.channelpath"))
	if err != nil {
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Unsupported version %d [%d=%d]", versionNumber, CCMSMIN, CCMSMAX)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	if len(pathVariables) != 2 {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Channel or date not found", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please verify the request and try again"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil

	}

	channelName := pathVariables[0]
	date := pathVariables[1]
	//validate the channel name
	channelNameIsValid := shared.ValidateChannelName(channelName)

	if channelNameIsValid {
		channelName, channelNameIsValid = core.GetChannelNameByAliasName(channelName)
	}

	if !channelNameIsValid {

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid channelname:[%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide a valid name"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	dayVal, monthVal, yearVal, dateIsValid := shared.ValidateDate(date)
	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] Channelname [%v] and Date:[%v] Year:[%v]", request.Method, transactionID, channelName, date, yearVal)

	if !dateIsValid {

		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid date:[%v]", request.Method, transactionID, date)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide a valid date"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	err = sh.DeleteSchedule(channelName, monthVal+dayVal)

	if err != nil {
		return internalServerErrorResponse(request, transactionID), nil
	}

	//delete the schedule from in-memory
	isDeleteFromMemory := deleteScheduleFromInMemory(channelName, monthVal+dayVal, request, transactionID)

	if isDeleteFromMemory {
		serviceLog.CcmsLog.Infof("[%v] TID:[%v] Schedule deleted for the channel [%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Schedule deleted for the channel %v", channelName)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusOK, "2xx"), nil
	} else {
		serviceLog.CcmsLog.Infof("[%v] TID:[%v] No Schedule found for channel [%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "No schedule found for the given date " + date})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusNotFound, "4xx"), nil

	}

}

// getACL returns the access control list for the ccms Service.
// func getACL() *acl.Config {
// 	acl := ccmsConfig.ACL
// 	if acl == nil {
// 		acl = sharedPoisConfig.GetACL()
// 	}
// 	return acl
// }

var ccmsService httpservice.RESTfulWebService = &serviceAPI{}

var serviceRoot string

// GetServiceAPIs returns the list of service APIs implemented
// by the CCMS service.
func GetServiceAPIs(handler httpservice.HttpServiceHandler, baseResource string) []*httpservice.HttpServiceAPI {
	apis := make([]*httpservice.HttpServiceAPI, 0, 0)
	serviceRoot = baseResource
	// CcmsConfig := getACL()
	for i := CCMSMIN + 1; i < CCMSMAX; i++ {
		apis = append(apis,
			&httpservice.HttpServiceAPI{
				ApiRoot:      config.GetConfig().GetString("api.rootpath") + "v" + strconv.Itoa(i) + serviceRoot + "/",
				ReadACL:      nil,
				WriteACL:     nil,
				CheckHeaders: false,
				Handler:      nil,
				WebService:   ccmsService,
			})
	}

	return apis
}

// This function will hold the stucture for the response type for the ccms request
// with the status, header , body and content-type values.
func responseFormat(Status int, header http.Header, body []byte, ContentType string) *httpservice.HttpServiceResponse {

	response := &httpservice.HttpServiceResponse{
		Status:      Status,
		Header:      header,
		Body:        body,
		ContentType: ContentType,
	}
	return response
}

// internal server response
func internalServerErrorResponse(request *http.Request, transactionID uint64) *httpservice.HttpServiceResponse {
	serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Struct to byte convertion error", request.Method, transactionID)
	prometheus.CcmsDefaultRequest.WithLabelValues().Inc()
	prometheus.CcmsRequest.WithLabelValues("5xx", request.Method, request.Proto).Inc()
	return shared.InternalServerErrorResponse()
}

// BadRequest Response Types
func requestResponse(request *http.Request, responseMessage []byte, responseType int, label string) *httpservice.HttpServiceResponse {
	response := responseFormat(responseType, nil, responseMessage, "application/json")
	prometheus.CcmsDefaultRequest.WithLabelValues().Inc()
	prometheus.CcmsRequest.WithLabelValues("4xx", request.Method, request.Proto).Inc()
	return response
}

// delete the schedule from the in-memory
func deleteScheduleFromInMemory(channelName string, date string, request *http.Request, transactionId uint64) bool {
	ch := sh.ChannelSchedules{Channel: channelName, Date: date}

	if _, ok := sh.ChannelScheduleData.GetSchedule(ch); ok {
		sh.ChannelScheduleData.DelSchedule(ch)
		serviceLog.CcmsLog.Infof("[%v] TID:[%v] Schedule deleted for the channel [%v]", request.Method, transactionId, channelName)
		return true
	}

	return false
}
