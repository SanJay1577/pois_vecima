// Package classification POIS API Documentation
//
// # The POIS API documentation is organized around REST and SOAP interfaces. REST API accepts binary and JSON request bodies and returns JSON and text responses. Our application used standard HTTP response codes and verbs
//
// BasePath: /pois/v1/channels
// Scheme: http, https
// Version: 0.0.1
// License: Vecima Networks Inc.
// Contact:
//
// Consumes:
// - application/json
// - application/xml
// Produces:
// - application/json
// - application/xml
// swagger:meta
package alias

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pois/config"
	serviceLog "pois/config/logging"
	"pois/models"
	"pois/prometheus"

	shared "pois/shared"
	"strconv"

	"git.eng.vecima.com/cloud/golib/v4/httpservice"
)

// WebService API resource path.
const (
	AliasVMIN = iota
	AliasVER1
	AliasVMAX
)

type Alias struct {
	AliasNames []string `json:"aliasNames"`
}

type serviceAPI struct{}

// applicationConfiguration
//var applicationConfig = config.GetApplicationConfig()

// This function will hold the stucture for the response type for the alias request
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
	serviceLog.AliasLog.Errorf("[%v] TID:[%v] Struct to byte convertion error", request.Method, transactionID)
	prometheus.AliasDefaultRequest.WithLabelValues().Inc()
	prometheus.AliasRequest.WithLabelValues("5xx", request.Method, request.Proto).Inc()
	return shared.InternalServerErrorResponse()
}

// BadRequest Response Types
func requestResponse(request *http.Request, responseMessage []byte, responseType int, label string) *httpservice.HttpServiceResponse {
	response := responseFormat(responseType, nil, responseMessage, "application/json")
	prometheus.AliasDefaultRequest.WithLabelValues().Inc()
	prometheus.AliasRequest.WithLabelValues(label, request.Method, request.Proto).Inc()
	return response
}

// Get handles a GET request.
// This API returns the list of alias names for the provided channels
// The Below API extract the channel name from requested url.
// if the channels have alias will gives a response as list of aliasnames
// otherwise this will return no alias found the channel
//
// swagger:route GET /alias/{channelname} ALIASAPI retrieveAliasnames
// # Retrieve alias names
//
// This API retrieves the alias names for the given channel. It accepts the channel name  from the path parameters and responds with list of alias for the given channel name
//
// Produces:
// - application/json
// Parameters:
//   - name: channelname
//     in: path
//     description: Channel name to retrieve alias names
//     type: string
//     required: true
//
// Responses:
//
//	200: GetAliasResponse (Retrieves the list of aliases for the channel)
//	400: BadRequestResponse (Channel name not provided)
//	404: NotFoundResponse (No alias found for the given channel)
func (h *serviceAPI) Get(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {

	serviceLog.AliasLog.Infof("[%v] TID:[%v] Path:[%v] ", request.Method, transactionID, request.URL.Path)
	pathVariables, versionNumber, err := shared.ExtractPathVariable(request.URL.Path, config.GetConfig().GetString("api.rootpath"), AliasVMIN, AliasVMAX, config.GetConfig().GetString("api.aliaspath"))
	serviceLog.AliasLog.Debugf("[%v] TID:[%v] Path Variables:%v", request.Method, transactionID, pathVariables)
	if err != nil {
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Unsupported version %d [%d=%d]", versionNumber, AliasVMIN, AliasVMAX)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	if len(pathVariables) >= 2 {
		serviceLog.AliasLog.Errorf("[%v] TID:[%v] Channel name missing", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please verify the request and try again"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	channelName := pathVariables[0]
	channelNameIsValid := shared.ValidateChannelName(channelName)

	if !channelNameIsValid {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid channelname:[%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the valid channel name"})

		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, responseMessage, http.StatusBadRequest, "4xx"), nil
	}

	var aliasChannelNameSlice []string

	for key, val := range shared.AliasChannelMap.GetAllChannels() {
		if val == channelName {
			aliasChannelNameSlice = append(aliasChannelNameSlice, key)
		}
	}
	if len(aliasChannelNameSlice) > 0 {

		serviceLog.AliasLog.Debugf("[%v] TID:[%v] Alias names %v for the channel [%v]", request.Method, transactionID, aliasChannelNameSlice, channelName)
		responseData, err := json.Marshal(Alias{AliasNames: aliasChannelNameSlice})

		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		serviceLog.AliasLog.Infof("[%v] TID:[%v] Alias names provided for the  given channel [%v]", request.Method, transactionID, channelName)
		return requestResponse(request, []byte(responseData), http.StatusOK, "2xx"), nil

	} else {

		serviceLog.AliasLog.Infof("[%v] TID:[%v] No Alias found for the channel [%v]", request.Method, transactionID, channelName)
		responseData := " No alias found for the channel " + channelName

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: responseData})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, []byte(responseMessage), http.StatusNotFound, "4xx"), nil
	}

}

// Post handles a POST request.
func (h *serviceAPI) Post(request *http.Request, transactionId uint64) (*httpservice.HttpServiceResponse, error) {

	response := &httpservice.HttpServiceResponse{
		Status:      http.StatusNotImplemented,
		Header:      nil,
		ContentType: "application/json",
		Body:        []byte("POST /channels/alias is not supported"),
	}
	prometheus.AliasDefaultRequest.WithLabelValues().Inc()
	prometheus.AliasRequest.WithLabelValues("5xx", request.Method, request.Proto).Inc()
	return response, nil
}

// Put handles a PUT request.

// The Below API extract the channel name from  the requested url.
//
//	we stored the alias names as keys, and value as channel name in memory. after storing
//
// the alias names return success response as mapped the alias for the  given channel
// swagger:route PUT /alias/{channelname} ALIASAPI addAlias
//
// Create/Update alias names.
//
// This API adds or updates the alias names for the given  channel. It accepts the channel name  from the path params and a payload with a list of alias names. This API's payload must be JSON data  with Content-Type  as application/json
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// Parameters:
//   - name: channelname
//     in: path
//     description: Channel name to add/update alias names
//     type: string
//     required: true
//
// Responses:
//
//	201: AliasCreated (Alias are created/mapped for the corresponding channel)
//	400: BadRequestResponse (Channel name or request body error parsing error)
func (h *serviceAPI) Put(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {

	var response *httpservice.HttpServiceResponse
	var err error

	serviceLog.AliasLog.Infof("[%v] TID:[%v] Path:[%v]", request.Method, transactionID, request.URL.Path)
	pathVariables, versionNumber, err := shared.ExtractPathVariable(request.URL.Path, config.GetConfig().GetString("api.rootpath"), AliasVMIN, AliasVMAX, config.GetConfig().GetString("api.aliaspath"))
	serviceLog.AliasLog.Debugf("[%v] TID:[%v] Path Variables:%v", request.Method, transactionID, pathVariables)
	if err != nil {

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Unsupported version %d [%d=%d]", versionNumber, AliasVMIN, AliasVMAX)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
	}
	if len(pathVariables) >= 2 {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Channel name missing", request.Method, transactionID)

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please verify the request and try again"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
	}

	if len(pathVariables) == 0 && request.ContentLength == 0 {

		serviceLog.AliasLog.Errorf("[%v] TID:[%v] Channel name  and request body not provided", request.Method, transactionID)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("[%d] PUT unable to read request body and path variable: %s", transactionID, err.Error())})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
	}

	channelName := pathVariables[0]
	channelNameIsValid := shared.ValidateChannelName(channelName)

	if !channelNameIsValid {
		serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid channelname:[%v]", request.Method, transactionID, channelName)
		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the valid channel name"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
	}

	if request.Body != nil && request.ContentLength != 0 {
		var requestBody Alias
		body, err := ioutil.ReadAll(request.Body)
		serviceLog.AliasLog.Debugf("[%v] TID:[%v] Got the request body", request.Method, transactionID)

		if err != nil {
			serviceLog.AliasLog.Errorf("[%v] TID:[%v] request body parsing  failed", request.Method, transactionID)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Request Body parsing failed: %v", err)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			response = &httpservice.HttpServiceResponse{
				Status:      http.StatusInternalServerError,
				Header:      nil,
				Body:        []byte(responseMessage),
				ContentType: "application/json",
			}
			prometheus.AliasDefaultRequest.WithLabelValues().Inc()
			prometheus.AliasRequest.WithLabelValues("5xx", request.Method, request.Proto).Inc()
			return response, nil
		}

		defer request.Body.Close()

		err = json.Unmarshal(body, &requestBody)

		if err != nil {
			serviceLog.AliasLog.Errorf("[%v] TID:[%v] Request body not valid", request.Method, transactionID)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Request body is not valid"})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}

			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil

		}

		//parse the alias names from request body
		aliasNamesList := requestBody.AliasNames
		if len(aliasNamesList) == 0 {
			serviceLog.AliasLog.Errorf("[%v] TID:[%v] Aliasname is empty  for the channel : [%v]", request.Method, transactionID, channelName)

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the alias name for the  channel : " + channelName})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}

			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
		}

		if len(aliasNamesList) != 0 {
			//check the key already exist in the memory
			serviceLog.AliasLog.Debugf("[%v] TID:[%v] Alias names provided for the channel [%v]", request.Method, transactionID, channelName)
			for _, aliasValue := range aliasNamesList {
				shared.AliasChannelMap.SetChannel(aliasValue, channelName)

				var aliasmodel models.Alias
				aliasmodel.Channelname = channelName
				aliasmodel.AliasName = aliasValue
				shared.CreateAlias(aliasmodel)

			}
			shared.AliasChannelMap.SetChannel(channelName, channelName)

			serviceLog.AliasLog.Infof("[%v] TID:[%v] Aliasname created for the channel[%v]", request.Method, transactionID, channelName)

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Alias names mapped for the channel %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusCreated, "2xx"), nil

		}

	}
	return response, nil
}

// Delete handles a DELETE request.
// This below API extract the channel name and aliasnames from requested url.
// First will check the alias is present in the channel. if the channel have an alias name. will remove the alias name from memory return the succcessful response
// if the channel doesn't have  an alias returns as not found reponse no alias is found for the channe

// swagger:route DELETE /alias/{channelname}/{aliasname} ALIASAPI RemoveAlias
//
// # Delete alias names
//
// This API deletes the alias name for the given channel. It accepts the channel name and alias name from the path parameters to delete the channel's alias name(s)
//
// Produces:
// - application/json
//
// Parameters:
//
//   - name: channelname
//     in: path
//     description: Channel name to delete alias names
//     type: string
//     required: true
//
//   - name: aliasname
//     in: path
//     description: Alias name of the channel
//     type: string
//     required: true
//
// Responses:
//
//	200:deleteSuccessResponse (Alias names are deleted for the given channel {{channel name}})
//	400:BadRequestResponse (Channel name or alias name  not provided)
//	404:NotFoundResponse (No alias found for the channel {{channel name}})
func (h *serviceAPI) Delete(request *http.Request, transactionID uint64) (*httpservice.HttpServiceResponse, error) {
	var response *httpservice.HttpServiceResponse

	serviceLog.AliasLog.Infof("[%v] TID:[%v] Path:[%v]", request.Method, transactionID, request.URL.Path)

	pathVariables, versionNumber, err := shared.ExtractPathVariable(request.URL.Path, config.GetConfig().GetString("api.rootpath"), AliasVMIN, AliasVMAX, config.GetConfig().GetString("api.aliaspath"))

	if err != nil {

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Unsupported version %d [%d=%d]", versionNumber, AliasVMIN, AliasVMAX)})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}
		prometheus.AliasDefaultRequest.WithLabelValues().Inc()
		return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil

	}
	length := len(pathVariables)
	switch length {

	case 0:
		serviceLog.AliasLog.Errorf("[%v] TID:[%v] Channelname not provided", request.Method, transactionID)

		responseMessage, err := json.Marshal(shared.ResponseMessage{Message: "Please provide the valid request and try again"})
		if err != nil {
			return internalServerErrorResponse(request, transactionID), nil
		}

		prometheus.AliasRequest.WithLabelValues("4xx", request.Method, request.Proto).Inc()
		return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
	case 1:
		//delete all the alias for the channel name
		channelName := pathVariables[0]

		channelNameIsValid := shared.ValidateChannelName(channelName)

		if !channelNameIsValid {
			serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid channelname:[%v]", request.Method, transactionID, channelName)

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Please provide valid channel name %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
		}

		if channelName == "" {

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Please provide valid channel name %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
		}
		serviceLog.AliasLog.Infof("[%v] TID:[%v] Deleted all alias from the channel[%v]", request.Method, transactionID, channelName)

		if shared.AliasChannelMap.GetChannel(channelName) != "" {
			for alias, channel := range shared.AliasChannelMap.GetAllChannels() {
				//remove all the alias name, keep the default channel name as alias name
				if channel == channelName && alias != channelName {
					shared.AliasChannelMap.DelChannel(alias)
				}
			}
			serviceLog.AliasLog.Infof("[%v] TID:[%v] all alias is deleted for the channel [%v]", request.Method, transactionID, channelName)

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Alias names deleted for the channel %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusOK, "2xx"), nil
		} else {

			serviceLog.AliasLog.Infof("[%v] TID:[%v] No Alias found for the channel [%v]", request.Method, transactionID, channelName)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("No alias found for the  channel  %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusNotFound, "4xx"), nil
		}
	case 2:

		channelName := pathVariables[0]
		aliasName := pathVariables[1]
		channelNameIsValid := shared.ValidateChannelName(channelName)
		if !channelNameIsValid {
			serviceLog.CcmsLog.Errorf("[%v] TID:[%v] Invalid channelname:[%v]", request.Method, transactionID, channelName)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Please provide valid channel name %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
		}

		if aliasName == "" {
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Please provide alias name for the channel  %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
		}
		if channelName == "" {

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("Please provide valid channel name %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil

		}

		shared.DeleteAlias(channelName, aliasName)

		if shared.AliasChannelMap.GetChannel(aliasName) != "" {
			shared.AliasChannelMap.DelChannel(aliasName)

			serviceLog.AliasLog.Infof("[%v] TID:[%v] Alias is deleted for the channel [%v]", request.Method, transactionID, channelName)

			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("%v alias name deleted for the channel %v", aliasName, channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}

			response = &httpservice.HttpServiceResponse{
				Status:      http.StatusOK,
				Header:      nil,
				Body:        []byte(responseMessage),
				ContentType: "application/json",
			}

		} else {
			serviceLog.AliasLog.Infof("[%v] TID:[%v] [%v] No Alias found for the channel [%v]", request.Method, transactionID, aliasName, channelName)
			responseMessage, err := json.Marshal(shared.ResponseMessage{Message: fmt.Sprintf("No alias found for the  channel  %v", channelName)})
			if err != nil {
				return internalServerErrorResponse(request, transactionID), nil
			}
			return requestResponse(request, []byte(responseMessage), http.StatusBadRequest, "4xx"), nil
		}
	}
	prometheus.AliasDefaultRequest.WithLabelValues().Inc()
	prometheus.AliasRequest.WithLabelValues("2xx", request.Method, request.Proto).Inc()
	return response, nil
}

// func getACL() *acl.Config {
// 	acl := aliasConfig.ACL
// 	if acl == nil {
// 		acl = sharedPoisConfig.GetACL()
// 	}
// 	return acl
// }

var aliasService httpservice.RESTfulWebService = &serviceAPI{}
var serviceRoot string

// GetServiceAPIs returns the list of service APIs implemented
// by the Alias service.
func GetServiceAPIs(handler httpservice.HttpServiceHandler, baseResource string) []*httpservice.HttpServiceAPI {

	apis := make([]*httpservice.HttpServiceAPI, 0, 0)

	serviceRoot = baseResource
	// AliasConfig := getACL()
	for i := AliasVMIN + 1; i < AliasVMAX; i++ {
		apis = append(apis,
			&httpservice.HttpServiceAPI{
				ApiRoot:      config.GetConfig().GetString("api.rootpath") + "v" + strconv.Itoa(i) + serviceRoot + "/",
				ReadACL:      nil,
				WriteACL:     nil,
				CheckHeaders: false,
				Handler:      nil,
				WebService:   aliasService,
			})
	}
	return apis
}
