package esam

import "encoding/xml"

/// Structure used for xml request and xml response.

// Struct is used to handle the request variable in xml format.
type SignalProcessingEvents struct {
	XMLName        xml.Name       `xml:"SignalProcessingEvent" json:"-"`
	AcquiredSignal AcquiredSignal `xml:"AcquiredSignal" json:"AcquiredSignal"`
}

// Nested Struct that is used to hold the values of request parameter
type AcquiredSignal struct {
	AcquisitionPointIdentity string      `xml:"acquisitionPointIdentity,attr"`
	AcquisitionSignalID      string      `xml:"acquisitionSignalID,attr"`
	UTCPoint                 UTCPoint    `xml:"UTCPoint"`
	BinaryData               BinaryData  `xml:"BinaryData"`
	StreamTimes              StreamTimes `xml:"StreamTimes"`
}

// The below struct hold the values of scete35 payload
type BinaryData struct {
	BinaryData string `xml:",chardata"`
	SignalType string `xml:"signalType,attr"`
}

// The below struct holds the value of utc point data in the request fields
type UTCPoint struct {
	Utcpoint string `xml:"utcPoint,attr"`
}

// The below struct holds the value of streamtime data in the request fields
type StreamTimes struct {
	StreamTime StreamTime `xml:"StreamTime"`
}

// The below struct holds the value of streamtime attributes in the request fields
type StreamTime struct {
	TimeType  string `xml:"timeType,attr"`
	TimeValue string `xml:"timeValue,attr"`
}

// The below struct is used to store the values of response filed for replace response
type SignalProcessingNotification struct {
	XMLName        xml.Name       `xml:"SignalProcessingNotification" json:"-"`
	StatusCode     StatusCode     `xml:"Statuscode"`
	ResponseSignal ResponseSignal `xml:"ResponseSignal"`
}

// The below struct is used to store the values of response filed for delete and noop repsonse
type SignalProcessNotification struct {
	XMLName        xml.Name       `xml:"SignalProcessingNotification" json:"-"`
	StatusCode     StatusCode     `xml:"Statuscode"`
	Responsesignal Responsesignal `xml:"ResponseSignal"`
}

// The below strcut hold the nested values of the response field
type Responsesignal struct {
	UTCPoint                 UTCPoint `xml:"UTCPoint"`
	Action                   string   `xml:"action,attr"`
	AcquisitionPointIdentity string   `xml:"acquisitionPointIdentity,attr"`
	AcquisitionSignalID      string   `xml:"acquisitionSignalID,attr"`
}

// the response statuscode value and note attribute is stored in this struct
type StatusCode struct {
	Classcode string `xml:"classCode,attr"`
	Note      Note   `xml:"Note"`
}

type Note struct {
	Note string `xml:",chardata"`
}

// replace response values are stred in this struct
type ResponseSignal struct {
	UTCPoint                 UTCPoint              `xml:"UTCPoint"`
	Action                   string                `xml:"action,attr"`
	AcquisitionPointIdentity string                `xml:"acquisitionPointIdentity,attr"`
	AcquisitionSignalID      string                `xml:"acquisitionSignalID,attr"`
	SCTE35PointDescriptor    SCTE35PointDescriptor `xml:"SCTE35PointDescriptor,omitempty"`
	StreamTimes              StreamTimes           `xml:"StreamTimes,omitempty"`
}

// secete decoded values will be stored and passed in the response fields in this strcut
type SCTE35PointDescriptor struct {
	SpliceCommandType float64      `xml:"spliceCommandType,attr"`
	SpliceInsert      SpliceInsert `xml:"SpliceInsert"`
}

// Spliceinsert values from the secte35 decoded values will be stores in this struct
type SpliceInsert struct {
	SpliceEventId         float64 `xml:"spliceEventID,attr"`
	OutOfNetworkIndicator bool    `xml:"outOfNetworkIndicator,attr"`
	UniqueProgramID       float64 `xml:"uniqueProgramId,attr"`
	Duration              string  `xml:"duration,attr"`
}

// scete35 decode values will be stored inside the struct
type Scete35Data struct {
	SpliceCommand SpliceCommand `json:"spliceCommand"`
}

// seceteCommand type values
type SpliceCommand struct {
	Type                       int           `json:"type"`
	BreakDuration              BreakDuration `json:"breakDuration"`
	SpliceEventId              int           `json:"spliceEventId"`
	SpliceEventCancelIndicator bool          `json:"spliceEventCancelIndicator"`
	SpliceImmediateFlag        bool          `json:"spliceImmediateFlag"`
	OutOfNetworkIndicator      bool          `json:"outOfNetworkIndicator"`
	UniqueProgramId            int           `json:"uniqueProgramId"`
}

// Breakduration values of secte35 data will be stores in this struct.
type BreakDuration struct {
	AutoReturn bool    `json:"autoReturn"`
	Duration   float64 `json:"duration"`
}
