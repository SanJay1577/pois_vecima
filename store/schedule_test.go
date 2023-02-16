package store_test

import (
	"pois/models"
	"pois/store"
	"testing"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type ScheduleSuite struct {
	suite.Suite
	scheduleCount uint
}

var scheduleStore *store.SQLSchedule

func TestScheduleStore(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(ScheduleSuite))
}

func (suite *ScheduleSuite) SetupSuite() {
	logger := zaplogger.MainLog
	scheduleStore = &store.SQLSchedule{DB: DB, Log: logger}
}

func (suite *ScheduleSuite) SetupTest() {
	scheduleStore.DB.Exec("DELETE FROM ccms_schedule")
	suite.scheduleCount = 0
}

func (suite *ScheduleSuite) insertSchedule(channelname string, schedule string, eventtype string, scheduleddate string,
	scheduledtime string, windowstarttime string, windowdurationtime string, breakwithinwindow string, positionwithinbreak string,
	scheduledlength string, actualairedtime string, actualairedlength string, actualairedposition string, spotidentification string,
	statuscode string, userdefined string) {

	scheduleStore.DB.MustExec(
		`INSERT INTO ccms_schedule (
			channel_name, schedule, event_type, scheduled_date, scheduled_time, window_start_time, 
			window_duration_time, break_within_window, position_within_break, scheduled_length, 
			actual_aired_time, actual_aired_length, actual_aired_position, spot_identification, 
			status_code, user_defined
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		channelname, schedule, eventtype, scheduleddate, scheduledtime, windowstarttime, windowdurationtime,
		breakwithinwindow, positionwithinbreak, scheduledlength, actualairedtime, actualairedlength, actualairedposition,
		spotidentification, statuscode, userdefined,
	)
}

func (suite *ScheduleSuite) TestDeleteSchedule() {
	suite.insertSchedule("testChannel", "0101", "LOI", "0101", "001500", "0003", "0057", "001", "001",
		"000030", "000000", "00000000", "000", "00000039108", "0000", "PAPA")

	err := scheduleStore.DeleteSchedule("testChannel", "0101")
	suite.NoError(err)

	_, err = scheduleStore.FindSchedule("testChannel", "0101")
	suite.Nil(err)

	_, err = scheduleStore.FindScheduleByDate("0101")
	suite.Nil(err)

	err = scheduleStore.DeleteSchedule("testChannel", "0101")
	suite.Error(err, "shouldn't be able to delete a deleted schedule")
}

func (suite *ScheduleSuite) TestFindSchedule() {
	suite.insertSchedule("testChannel", "0101", "LOI", "0101", "001500", "0003", "0057", "001", "001",
		"000030", "000000", "00000000", "000", "00000039108", "0000", "PAPA")

	schedule, err := scheduleStore.FindSchedule("testChannel", "0101")

	if suite.NoError(err) {
		for _, element := range schedule {
			suite.Equal("testChannel", element.Channelname)
			suite.Equal("0101", element.Schedule)
			suite.Equal("LOI", element.Eventtype)
			suite.Equal("0101", element.Scheduleddate)
			suite.Equal("001500", element.Scheduledtime)
			suite.Equal("0003", element.Windowstarttime)
			suite.Equal("0057", element.Windowdurationtime)
			suite.Equal("001", element.Breakwithinwindow)
			suite.Equal("001", element.Positionwithinbreak)
			suite.Equal("000030", element.Scheduledlength)
			suite.Equal("000000", element.Actualairedtime)
			suite.Equal("00000000", element.Actualairedlength)
			suite.Equal("000", element.Actualairedposition)
			suite.Equal("00000039108", element.Spotidentification)
			suite.Equal("0000", element.Statuscode)
			suite.Equal("PAPA", element.Userdefined)

			suite.NotEqual("testChannelZ", element.Channelname)
			suite.NotEqual("0102", element.Schedule)
			suite.NotEqual("REM", element.Eventtype)
			suite.NotEqual("0102", element.Scheduleddate)
			suite.NotEqual("001502", element.Scheduledtime)
			suite.NotEqual("0002", element.Windowstarttime)
			suite.NotEqual("0052", element.Windowdurationtime)
			suite.NotEqual("002", element.Breakwithinwindow)
			suite.NotEqual("002", element.Positionwithinbreak)
			suite.NotEqual("000032", element.Scheduledlength)
			suite.NotEqual("000002", element.Actualairedtime)
			suite.NotEqual("00000002", element.Actualairedlength)
			suite.NotEqual("002", element.Actualairedposition)
			suite.NotEqual("00000039102", element.Spotidentification)
			suite.NotEqual("0002", element.Statuscode)
			suite.NotEqual("PAPZ", element.Userdefined)
		}
	}

	schedule, err = scheduleStore.FindScheduleByDate("0101")

	if suite.NoError(err) {
		for _, element := range schedule {
			suite.Equal("testChannel", element.Channelname)
			suite.Equal("0101", element.Schedule)
			suite.Equal("LOI", element.Eventtype)
			suite.Equal("0101", element.Scheduleddate)
			suite.Equal("001500", element.Scheduledtime)
			suite.Equal("0003", element.Windowstarttime)
			suite.Equal("0057", element.Windowdurationtime)
			suite.Equal("001", element.Breakwithinwindow)
			suite.Equal("001", element.Positionwithinbreak)
			suite.Equal("000030", element.Scheduledlength)
			suite.Equal("000000", element.Actualairedtime)
			suite.Equal("00000000", element.Actualairedlength)
			suite.Equal("000", element.Actualairedposition)
			suite.Equal("00000039108", element.Spotidentification)
			suite.Equal("0000", element.Statuscode)
			suite.Equal("PAPA", element.Userdefined)

			suite.NotEqual("testChannelZ", element.Channelname)
			suite.NotEqual("0102", element.Schedule)
			suite.NotEqual("REM", element.Eventtype)
			suite.NotEqual("0102", element.Scheduleddate)
			suite.NotEqual("001502", element.Scheduledtime)
			suite.NotEqual("0002", element.Windowstarttime)
			suite.NotEqual("0052", element.Windowdurationtime)
			suite.NotEqual("002", element.Breakwithinwindow)
			suite.NotEqual("002", element.Positionwithinbreak)
			suite.NotEqual("000032", element.Scheduledlength)
			suite.NotEqual("000002", element.Actualairedtime)
			suite.NotEqual("00000002", element.Actualairedlength)
			suite.NotEqual("002", element.Actualairedposition)
			suite.NotEqual("00000039102", element.Spotidentification)
			suite.NotEqual("0002", element.Statuscode)
			suite.NotEqual("PAPZ", element.Userdefined)
		}
	}

	_, err = scheduleStore.FindSchedule("testChannelZ", "0102")
	suite.Nil(err)

	err = scheduleStore.DeleteSchedule("testChannel", "0101")
	suite.NoError(err)
}

func (suite *ScheduleSuite) TestCreateSchedule() {
	request := models.Schedule{Channelname: "testChannel", Schedule: "0101", Eventtype: "LOI", Scheduleddate: "0101", Scheduledtime: "001500",
		Windowstarttime: "0003", Windowdurationtime: "0057", Breakwithinwindow: "001", Positionwithinbreak: "001",
		Scheduledlength: "000030", Actualairedtime: "000000", Actualairedlength: "00000000", Actualairedposition: "000", Spotidentification: "00000039108",
		Statuscode: "0000", Userdefined: "PAPA"}

	schedules, err := scheduleStore.CreateSchedule(request)

	if suite.NoError(err) {
		for _, element := range schedules {
			suite.Equal("testChannel", element.Channelname)
			suite.Equal("0101", element.Schedule)
			suite.Equal("LOI", element.Eventtype)
			suite.Equal("0101", element.Scheduleddate)
			suite.Equal("001500", element.Scheduledtime)
			suite.Equal("0003", element.Windowstarttime)
			suite.Equal("0057", element.Windowdurationtime)
			suite.Equal("001", element.Breakwithinwindow)
			suite.Equal("001", element.Positionwithinbreak)
			suite.Equal("000030", element.Scheduledlength)
			suite.Equal("000000", element.Actualairedtime)
			suite.Equal("00000000", element.Actualairedlength)
			suite.Equal("000", element.Actualairedposition)
			suite.Equal("00000039108", element.Spotidentification)
			suite.Equal("0000", element.Statuscode)
			suite.Equal("PAPA", element.Userdefined)

			suite.NotEqual("testChannelZ", element.Channelname)
			suite.NotEqual("0102", element.Schedule)
			suite.NotEqual("REM", element.Eventtype)
			suite.NotEqual("0102", element.Scheduleddate)
			suite.NotEqual("001502", element.Scheduledtime)
			suite.NotEqual("0002", element.Windowstarttime)
			suite.NotEqual("0052", element.Windowdurationtime)
			suite.NotEqual("002", element.Breakwithinwindow)
			suite.NotEqual("002", element.Positionwithinbreak)
			suite.NotEqual("000032", element.Scheduledlength)
			suite.NotEqual("000002", element.Actualairedtime)
			suite.NotEqual("00000002", element.Actualairedlength)
			suite.NotEqual("002", element.Actualairedposition)
			suite.NotEqual("00000039102", element.Spotidentification)
			suite.NotEqual("0002", element.Statuscode)
			suite.NotEqual("PAPZ", element.Userdefined)
		}
	}

	err = scheduleStore.DeleteSchedule("testChannel", "0101")
	suite.NoError(err)
}
