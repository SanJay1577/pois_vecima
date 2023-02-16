package esam

// Package that are required are imported in the below set of lines of code
import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"pois/config"
	serviceLog "pois/config/logging"
	core "pois/pois/core"
	esamCore "pois/pois/esam/core"
	"pois/prometheus"
	"pois/schedule"
	"strings"

	"github.com/Comcast/scte35-go/pkg/scte35"
)

// The below function will be responsible for the post xml request
// from which the request body data will be parsed, Unmarshalled and validated
// to send a proper response such as (delete, passthrough, replace),
// the below set of function will also decode and validate theh scete35 payload from the request body
// further will validate the channel name from ccms and get the duration value to pass in though the replace response
func EsamRequest(w http.ResponseWriter, request *http.Request) {
	serviceLog.EsamLog.Infof("[%v] request is initiated ", request.Method)

	serviceLog.EsamLog.Infof("Method:[%v] Path : [%v]", request.Method, request.URL.Path)

	pathProvider := strings.Split(request.URL.Path, "/")[3]

	serviceLog.EsamLog.Infof("Path provider name : [%v]", pathProvider)

	// Request body will have a xml data that are set to be validated to give the proper response
	// to acheive that the request body is read in the following line of code.
	reponseByte, err := ioutil.ReadAll(request.Body)
	// error handling is maintained in case of any error in reading the request body
	if err != nil {
		serviceLog.EsamLog.Errorf("[%v] xml is read from the request body Error: %v", request.Method, err)
		prometheus.EsamDefaultRequest.WithLabelValues().Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serviceLog.EsamLog.Infof("[%v] xml request body is read from the request body", request.Method)

	//  Request body is closed at the end after reading the data from the request xml file.
	defer request.Body.Close()

	// declaring a variable to unmarhall the request struct
	var signalProcessingEvent SignalProcessingEvents

	xml.Unmarshal(reponseByte, &signalProcessingEvent)
	serviceLog.EsamLog.Debugf("[%v] request body data is read and structured in signalprocessing event : %v", request.Method, signalProcessingEvent)

	// seting the headers for thre request type
	w.Header().Set("Content-Type", "application/xml")

	ResponseType, ResponseValidation = signalEventDataValidations(signalProcessingEvent)

	if !ResponseValidation {
		esamResponse(w, signalProcessingEvent)
		serviceLog.EsamLog.Infof("[%v] response is given for the request", request.Method)
		return
	}

	//Decoding the secte35 payload to get the streamtime, duration and spliceid and also to validate the fields.
	serviceLog.EsamLog.Infof("[%v] scete35 payload is set to decode the values", request.Method)
	secete35Decoder, err := scte35.DecodeBase64(signalProcessingEvent.AcquiredSignal.BinaryData.BinaryData)
	if err != nil {
		serviceLog.EsamLog.Errorf("[%v] incorrect scete35 payload is found in the request body while decoding Error: ", request.Method, err)

		prometheus.EsamDefaultRequest.WithLabelValues().Inc()
		prometheus.EsamDeleteRequest.WithLabelValues().Inc()

		w.Write(deleteResponse(signalProcessingEvent))
		serviceLog.EsamLog.Infof("[%v] response is given for the request", request.Method)
		return
	} else {

		// The below set of code will execute if the error statement from decoding the value is nil
		// defining a varible and ato unmarhall the secete35 deocoded contents.
		decodedJson, _ := json.MarshalIndent(secete35Decoder, "", "\t")

		// scete35 data is decoded in the struct value of SceteData
		var scete35Payload Scete35Data
		json.Unmarshal([]byte(string(decodedJson)), &scete35Payload)
		serviceLog.EsamLog.Debugf("[%v] scete35 payload is decoded and scete35 spliceCommand is parsed in json format  %v", request.Method, scete35Payload)
		serviceLog.EsamLog.Infof("[%v] scete35 payload is decode success ", request.Method)

		// The below function is called to validate the secete35 fields, requestbody fields and also to check
		// whether the channel and duration value is available from the ccms
		ResponseType, ResponseValidation = scete35DataValidations(scete35Payload)
		if !ResponseValidation {
			esamResponse(w, signalProcessingEvent)
			serviceLog.EsamLog.Infof("[%v] response is given for the request", request.Method)
			return
		}

		ResponseType, ResponseValidation = scheduleDatavalidations(signalProcessingEvent, pathProvider)

		if !ResponseValidation {
			esamResponse(w, signalProcessingEvent)
			serviceLog.EsamLog.Infof("[%v] response is given for the request", request.Method)
			return
		}

		prometheus.EsamDefaultRequest.WithLabelValues().Inc()
		prometheus.EsamSucessRequest.WithLabelValues().Inc()
		w.Write(replaceResponse(signalProcessingEvent, scete35Payload))
		serviceLog.EsamLog.Infof("[%v] response is given for the request", request.Method)
		return

	}
}

func esamResponse(w http.ResponseWriter, signalProcessingEvent SignalProcessingEvents) {
	switch ResponseType {
	case "delete":
		prometheus.EsamDefaultRequest.WithLabelValues().Inc()
		prometheus.EsamDeleteRequest.WithLabelValues().Inc()
		w.Write(deleteResponse(signalProcessingEvent))
	case "noop":
		prometheus.EsamDefaultRequest.WithLabelValues().Inc()
		prometheus.EsamNoopRequest.WithLabelValues().Inc()
		w.Write(noopResponse(signalProcessingEvent))
	default:
		prometheus.EsamDefaultRequest.WithLabelValues().Inc()
		prometheus.EsamDeleteRequest.WithLabelValues().Inc()
		w.Write(deleteResponse(signalProcessingEvent))
	}
}

// The below function is set to validate the Request fileds, secte35 fields, duration and channel informations
// based on the validations a response type will be returned from the config file that return value will be set to
// Responsetype variable from which we can send a required response.
func signalEventDataValidations(signal SignalProcessingEvents) (string, bool) {
	serviceLog.EsamLog.Infof("[%v] request and scete35 fields are set for validations", "POST")
	// Request body field validations are executed in the below lines of code.
	if signal.AcquiredSignal.AcquisitionPointIdentity == "" {
		serviceLog.EsamLog.Errorf("[POST] acquisitionpointidentity is not available in request field")
		return config.GetConfig().GetString("esam.response.acquisitionPointIdentity"), false

	}
	if signal.AcquiredSignal.UTCPoint.Utcpoint == "" {
		serviceLog.EsamLog.Errorf("[POST] utc point value is not available in request field")
		return config.GetConfig().GetString("esam.response.utcPoint"), false
	}
	if signal.AcquiredSignal.BinaryData.BinaryData == "" {
		serviceLog.EsamLog.Errorf("[POST] binarydata value is not available in request field")
		return config.GetConfig().GetString("esam.response.scete35Payload"), false

	}
	if signal.AcquiredSignal.BinaryData.SignalType == "" {
		serviceLog.EsamLog.Errorf("[POST] signal type value is not available in request field")
		return config.GetConfig().GetString("esam.response.signalType"), false

	}
	if signal.AcquiredSignal.StreamTimes.StreamTime.TimeType == "" {
		serviceLog.EsamLog.Errorf("[POST] streamtime type value is not available in request field")
		return config.GetConfig().GetString("esam.response.streamTimeType"), false

	}
	serviceLog.EsamLog.Infof("[POST] request and scete35 fields passed validations")
	// If all the above validations are passed then the function will return a valid response type.
	return config.GetConfig().GetString("esam.response.validRecord"), true

}

func scete35DataValidations(scete35Payload Scete35Data) (string, bool) {

	// secte35 Decoded values will be validated in below line of code to return a required response
	if scete35Payload.SpliceCommand.Type == 0 {
		serviceLog.EsamLog.Errorf("[POST] splice Command type is not available in decode secte35 playload")
		return config.GetConfig().GetString("esam.response.spliceCommandType"), false
	}
	if scete35Payload.SpliceCommand.SpliceEventId == 0 {
		serviceLog.EsamLog.Errorf("[POST] splice event type is not available in decode secte35 playload")
		return config.GetConfig().GetString("esam.response.spliceEventId"), false

	}
	if scete35Payload.SpliceCommand.BreakDuration.Duration == 0 {
		serviceLog.EsamLog.Errorf("[POST] breakDuration Duartion is not available in decode secte35 playload")
		return config.GetConfig().GetString("esam.response.spliceDuration"), false

	}
	serviceLog.EsamLog.Infof("[POST] request and scete35 fields passed validations")
	// If all the above validations are passed then the function will return a valid response type.
	return config.GetConfig().GetString("esam.response.validRecord"), true

}

func scheduleDatavalidations(signal SignalProcessingEvents, pathProvider string) (string, bool) {
	// The below set of conditions will call a function which will return a channel name in string and and channel Validation
	// as a boolean value, By validating the channelValidation Value , it will return a reqiuired response type.
	channelName, channelValidation := core.GetChannelNameByAliasName(signal.AcquiredSignal.AcquisitionPointIdentity)
	if !channelValidation {
		serviceLog.EsamLog.Errorf("[POST] acquisitionpointidentity is not available in request field")
		return config.GetConfig().GetString("esam.response.acquisitionPointIdentity"), false
	}

	// The below line of code will call a function that will convert the utc point value in the
	// desired date and time format data - DDMM time - HHMMSS ...
	date, time := esamCore.UtcFormater(signal.AcquiredSignal.UTCPoint.Utcpoint)
	//The Below lines of code will call a function that will return the duration value and boolean value
	// Based on the return response the varible will be assigned and validated based on those validations
	// it will return a reponse type.

	DurationValue, DurationValidation = schedule.GetSchedule.GetProviderSchedule(schedule.Provider(pathProvider), pathProvider, channelName, date, time)
	//DurationValue = core.GetScheduleByChannelAndTime(channelName, date, time)

	serviceLog.EsamLog.Infof("[POST] Return DurationValue %v DurationValidation %v", DurationValue, DurationValidation)
	if DurationValue == "" {
		return config.GetConfig().GetString("esam.response.spliceDuration"), false

	}
	serviceLog.EsamLog.Infof("[POST] request and scete35 fields passed validations")
	// If all the above validations are passed then the function will return a valid response type.
	return config.GetConfig().GetString("esam.response.validRecord"), true

}

// The below function will allow us to give delete response in a xml format
// The fields and values will be declared and marshalled into xml format.
func deleteResponse(signalEvent SignalProcessingEvents) []byte {
	delete := &SignalProcessNotification{}
	delete.StatusCode.Classcode = "0"
	delete.Responsesignal.Action = "delete"
	delete.Responsesignal.AcquisitionPointIdentity = signalEvent.AcquiredSignal.AcquisitionPointIdentity
	delete.Responsesignal.AcquisitionSignalID = signalEvent.AcquiredSignal.AcquisitionSignalID
	delete.StatusCode.Note.Note = "Discarding AcquiredSignal acquisitionSignalID "
	delete.Responsesignal.UTCPoint.Utcpoint = signalEvent.AcquiredSignal.UTCPoint.Utcpoint
	responseDelete, _ := xml.MarshalIndent(delete, " ", " ")
	return responseDelete
}

// The below function will allow us to give noop response in a xml format
// The fields and values will be declared and marshalled into xml format.
func noopResponse(signalEvent SignalProcessingEvents) []byte {
	noop := &SignalProcessNotification{}
	noop.StatusCode.Classcode = "0"
	noop.Responsesignal.Action = "noop"
	noop.Responsesignal.AcquisitionPointIdentity = signalEvent.AcquiredSignal.AcquisitionPointIdentity
	noop.Responsesignal.AcquisitionSignalID = signalEvent.AcquiredSignal.AcquisitionSignalID
	noop.StatusCode.Note.Note = "Discarding AcquiredSignal acquisitionSignalID "
	noop.Responsesignal.UTCPoint.Utcpoint = signalEvent.AcquiredSignal.UTCPoint.Utcpoint
	responseNoop, _ := xml.MarshalIndent(noop, " ", " ")
	return responseNoop
}

// The below function will allow us to give replace response in a xml format
// The fields and values will be declared and marshalled into xml format.
func replaceResponse(signalEvent SignalProcessingEvents, scete35Payload Scete35Data) []byte {
	replace := &SignalProcessingNotification{}
	replace.StatusCode.Classcode = "0"
	replace.StatusCode.Note.Note = "insert normalized signal"
	replace.ResponseSignal.Action = "replace"
	replace.ResponseSignal.AcquisitionPointIdentity = signalEvent.AcquiredSignal.AcquisitionPointIdentity
	replace.ResponseSignal.AcquisitionSignalID = signalEvent.AcquiredSignal.AcquisitionSignalID
	replace.ResponseSignal.UTCPoint.Utcpoint = signalEvent.AcquiredSignal.UTCPoint.Utcpoint
	replace.ResponseSignal.SCTE35PointDescriptor.SpliceCommandType = float64(scete35Payload.SpliceCommand.Type)
	replace.ResponseSignal.SCTE35PointDescriptor.SpliceInsert.SpliceEventId = float64(scete35Payload.SpliceCommand.SpliceEventId)
	replace.ResponseSignal.SCTE35PointDescriptor.SpliceInsert.OutOfNetworkIndicator = scete35Payload.SpliceCommand.OutOfNetworkIndicator
	replace.ResponseSignal.SCTE35PointDescriptor.SpliceInsert.UniqueProgramID = float64(scete35Payload.SpliceCommand.UniqueProgramId)
	replace.ResponseSignal.SCTE35PointDescriptor.SpliceInsert.Duration = esamCore.DurationTypeConverstion(DurationValue)
	replace.ResponseSignal.StreamTimes.StreamTime.TimeType = signalEvent.AcquiredSignal.StreamTimes.StreamTime.TimeType
	replace.ResponseSignal.StreamTimes.StreamTime.TimeValue = signalEvent.AcquiredSignal.StreamTimes.StreamTime.TimeValue
	responseReplace, _ := xml.MarshalIndent(replace, " ", " ")
	return responseReplace
}
