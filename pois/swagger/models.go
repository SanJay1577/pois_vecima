package swagger

/* Alias Request Response Models */

//swagger:model deleteSuccessResponse
type DeleteAliasResponse struct {
	//example:Alias name deleted for the channel {{channel name}}
	Message string `json:"message"`
}

//swagger:model BadRequestResponse
type BadRequestErrorResponse struct {
	//example:Please verify the request and try again
	Message string `json:"message"`
	//example:nil
}

//swagger:model NotFoundResponse
type AliastNotFoundResponse struct {
	//example:No alias found for the channel {{channel name}}
	Message string `json:"message"`
}

//swagger:model GetAliasResponse
type AliasGetResponse struct {
	//example:["cnnlive1","cnnlive2"]
	AliasNames []string `json:"aliasNames"`
	//example:nil
}

//swagger:model AliasCreated
type AliasCreatedResonse struct {
	//example:Alias names mapped for the channel {{channel name}}
	Message string `json:"message"`
}

//swagger:parameters addAlias
type AliasRequest struct {
	//alias names for the channels
	//in:body
	Body struct {
		//example:["cnnlive1","cnnlive2"]
		AliasNames []string `json:"aliasNames"`
	}
}

//swagger:model DeleteSchdeuleResponse
type DeleteScheduleResponse struct {
	//example:Schedule removed for the channel in given date {{DDMMYYYY}}
	Message string `json:"message"`
}

//swagger:model BadRequestResponse
type BadRequestResponse struct {

	//example:Please verify the request and try again
	Message string `json:"message"`
}

//swagger:model ChannelNotFoundResponse
type ScheduleNotFoundResponse struct {
	//example:No Schedule found for the channel in a given date {{DDMMYYYY}}
	Message string `json:"message"`
}

//swagger:model ScheduleCreated
type ScheduleCreatedResponse struct {
	//example:Schedule created for the channels in a provided date {{DDMMYYYY}}
	Message string `json:"message"`
}

//swagger:model preprocessingFails
type PreprocessingFailureResponse struct {
	//example:Unable to preprocess the file
	Message string `json:"message"`
}

//swagger:parameters addUpateSchedule
type ChannelScheduleRequest struct {
	//Channel Schedular Information
	//in:body
	Body struct {
		//example:REM Created on 02/09/99 13:45 TNT COLUMBUS LOI 0210 001500 0003 0057 001 001 000030 000000 00000000 000 00000039108 0000 PAPA LOI 0210 001500 0003 0057 001 002 000030 000000 00000000 000 00000018709 0000 ESPN END
		ChannelSchdedule string
	}
}

//swagger:model scheduleRetrivalSuccessResponse
type ScheduleRetrivalResponse struct {
	TimeDurationMap
}

type TimeDurationMap struct {
	TimeDuration Schedules `json:"001500"`
}
type Schedules struct {
	ScheduleSlices []Schedule `json:"Schedule"`
}

type Schedule struct {
	//example: LOI
	EventType string `json:"EventType"`
	//example:0110
	ScheduledDate string `json:"ScheduledDate"`
	//example:003000
	ScheduledTime string `json:"ScheduledTime"`
	//example:0003
	WindowStartTime string `json:"WindowStartTime"`
	//example:0057
	WindowDurationTime string `json:"WindowDurationTime"`
	//example:002
	BreakWithinWindow string `json:"BreakWithinWindow"`
	//example:001
	PositionWithInBreak string `json:"PositionWithInBreak"`
	//example:000030
	ScheduledLength string `json:"ScheduledLength"`
	//example:000000
	ActualAiredTime string `json:"ActualAiredTime"`
	//example:00000000
	ActualAiredLength string `json:"ActualAiredLength"`
	//example:000
	ActualAiredPosition string `json:"ActualAiredPosition"`
	//example:00000021902
	SpotIdentification string `json:"SpotIdentification"`
	//example:0000
	Statuscode string `json:"Statuscode"`
	//example:AIR
	UserDefined string `json:"UserDefined"`
}
