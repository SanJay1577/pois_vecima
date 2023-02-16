package shared_test

import (
	"pois/shared"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// The Below Test function is to test the function which extract the path variables of the given url
func TestExtractPathVariabels(t *testing.T) {
	// CCMS Url Scenario
	ccmspath := "/pois/v1/channels/cnn/02022023"
	ccmsbasePath := "/pois/"
	ccmsmin := 0
	ccmsmax := 2
	ccmsrootPath := "channels"
	ccmstestLength := 2

	ccmsPathVariables, ccmsVersion, ccmsErr := shared.ExtractPathVariable(ccmspath, ccmsbasePath, ccmsmin, ccmsmax, ccmsrootPath)

	assert.Equal(t, len(ccmsPathVariables), ccmstestLength)
	assert.Equal(t, ccmsVersion, 1)
	assert.Equal(t, ccmsErr, nil)

	//ALias Url Scenario
	aliaspath := "/pois/v1/channels/alias/cnn"
	aliasbasePath := "/pois/"
	aliasmin := 0
	alliasmax := 2
	aliasrootPath := "channels/alias"
	aliastestLength := 1

	aliasPathVariables, aliasVersion, aliasErr := shared.ExtractPathVariable(aliaspath, aliasbasePath, aliasmin, alliasmax, aliasrootPath)

	assert.Equal(t, len(aliasPathVariables), aliastestLength)
	assert.Equal(t, aliasVersion, 1)
	assert.Equal(t, aliasErr, nil)
}

// Test to Check the Valid Channel Names
func Test_nameValidation(t *testing.T) {
	var channelName string
	validName := true
	notValidName := false

	// Positive Test Cases
	channelName = "cnn123"
	test1 := shared.ValidateChannelName(channelName)
	assert.Equal(t, test1, validName)

	channelName = "cnn&&"
	test2 := shared.ValidateChannelName(channelName)
	assert.Equal(t, test2, validName)

	// Negative Test Cases
	channelName = "cnn??"
	test3 := shared.ValidateChannelName(channelName)
	assert.Equal(t, test3, notValidName)

	channelName = "cnn*"
	test4 := shared.ValidateChannelName(channelName)
	assert.Equal(t, test4, notValidName)

}

// Test function to Comapre the Valid date functionality.
func TestComapreDate(t *testing.T) {
	validDate := true
	notValidDate := false

	test1 := shared.CompareDate("01", "12", "2023")
	assert.Equal(t, test1, validDate)

	test2 := shared.CompareDate("05", "11", "2024")
	assert.Equal(t, test2, validDate)

	test3 := shared.CompareDate("01", "11", "2021")
	assert.Equal(t, test3, notValidDate)

	test4 := shared.CompareDate("05", "13", "2023")
	assert.Equal(t, test4, notValidDate)

}

// Validating the Validdate functionality
func TestValidateDate(t *testing.T) {
	expectedDate := "02"
	expectedMonth := "02"
	expectedYear := "2023"
	validdate := true

	actualDate, actualMonth, actualYear, actualTest := shared.ValidateDate("02022023")

	assert.Equal(t, expectedDate, actualDate)
	assert.Equal(t, expectedMonth, actualMonth)
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, validdate, actualTest)

}

// Deleting inmemory alias name function test cases
func TestTTL(t *testing.T) {
	shared.InitailizeCleanUp(1 * time.Second)
	key := shared.ChannelSchedules{Channel: "00", Date: "1500"}

	scheduleinfo := shared.ScheduleInfo{
		EventType:           "LOI",
		ScheduledDate:       "0117",
		ScheduledTime:       "001500",
		WindowStartTime:     "0003",
		WindowDurationTime:  "0057",
		BreakWithinWindow:   "001",
		PositionWithinBreak: "001",
		ScheduledLength:     "000030",
		ActualAiredTime:     "000000",
		ActualAiredLength:   "00000000",
		ActualAiredPosition: "000",
		SpotIdentification:  "00000039108",
		StatusCode:          "0000",
		UserDefined:         "PAPA\r",
	}

	schedules := shared.Schedules{
		Schedules: []shared.ScheduleInfo{
			scheduleinfo,
		},
	}

	var scheduleMap = make(map[string]shared.Schedules)
	scheduleMap["001500"] = schedules
	shared.ChannelScheduleData.SetSchedule(key, scheduleMap)
	_, ok := shared.ChannelScheduleData.GetSchedule(key)
	assert.Equal(t, true, ok)
	t.Log("Mocking waiting time to delete schedule information")
	time.Sleep(35 * time.Second)
	_, after := shared.ChannelScheduleData.GetSchedule(key)
	assert.Equal(t, false, after)

}
