package store

import (
	"pois/models"
	"strings"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
	"github.com/jmoiron/sqlx"
)

type Schedule interface {
	FindSchedule(string, string) ([]models.Schedule, error)
	CreateSchedule(models.Schedule) (models.Schedule, error)
	DeleteSchedule(string, string) error
}

// SQLSchedule is a DB-backed concrete store.
type SQLSchedule struct {
	DB  *sqlx.DB
	Log *zaplogger.Logger
}

// FindSchedule fetch all schedules for the channel on a given date from DB.
func (s *SQLSchedule) FindSchedule(channelname string, scheduledate string) ([]models.Schedule, error) {
	schedule := []models.Schedule{}
	err := s.DB.Select(&schedule, "SELECT channel_name, schedule, event_type, scheduled_date, scheduled_time, window_start_time, window_duration_time, break_within_window, position_within_break, scheduled_length, actual_aired_time, actual_aired_length, actual_aired_position, spot_identification, status_code, user_defined FROM ccms_schedule where channel_name = $1 and schedule = $2",
		channelname, scheduledate)
	return schedule, err
}

// FindScheduleByDate fetch all schedules for the given date from DB.
func (s *SQLSchedule) FindScheduleByDate(scheduledate string) ([]models.Schedule, error) {
	schedule := []models.Schedule{}
	err := s.DB.Select(&schedule, "SELECT channel_name, schedule, event_type, scheduled_date, scheduled_time, window_start_time, window_duration_time, break_within_window, position_within_break, scheduled_length, actual_aired_time, actual_aired_length, actual_aired_position, spot_identification, status_code, user_defined FROM ccms_schedule where schedule >= $1",
		scheduledate)
	return schedule, err
}

// CreateScheduler Creates a schedule for a channel for the date in the DB
func (s *SQLSchedule) CreateSchedule(request models.Schedule) ([]models.Schedule, error) {

	schedules := []models.Schedule{}

	channelname := strings.TrimSpace(request.Channelname)
	schedule := strings.TrimSpace(request.Schedule)
	eventtype := strings.TrimSpace(request.Eventtype)
	scheduleddate := strings.TrimSpace(request.Scheduleddate)
	scheduledtime := strings.TrimSpace(request.Scheduledtime)
	windowstarttime := strings.TrimSpace(request.Windowstarttime)
	windowdurationtime := strings.TrimSpace(request.Windowdurationtime)
	breakwithinwindow := strings.TrimSpace(request.Breakwithinwindow)
	positionwithinbreak := strings.TrimSpace(request.Positionwithinbreak)
	scheduledlength := strings.TrimSpace(request.Scheduledlength)
	actualairedtime := strings.TrimSpace(request.Actualairedtime)
	actualairedlength := strings.TrimSpace(request.Actualairedlength)
	actualairedposition := strings.TrimSpace(request.Actualairedposition)
	spotidentification := strings.TrimSpace(request.Spotidentification)
	statuscode := strings.TrimSpace(request.Statuscode)
	userdefined := strings.TrimSpace(request.Userdefined)

	err := s.DB.Select(
		&schedules,
		`INSERT INTO ccms_schedule (
			channel_name, schedule, event_type, scheduled_date, scheduled_time, window_start_time, window_duration_time, break_within_window, position_within_break, scheduled_length, actual_aired_time, actual_aired_length, actual_aired_position, spot_identification, status_code, user_defined
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING channel_name, schedule, event_type, scheduled_date, scheduled_time, window_start_time, window_duration_time, break_within_window, position_within_break, scheduled_length, actual_aired_time, actual_aired_length, actual_aired_position, spot_identification, status_code, user_defined`,
		channelname, schedule, eventtype, scheduleddate, scheduledtime, windowstarttime, windowdurationtime, breakwithinwindow, positionwithinbreak, scheduledlength, actualairedtime, actualairedlength, actualairedposition, spotidentification, statuscode, userdefined,
	)

	return schedules, err
}

// DeleteSchedule deletes a schedules based on channel name and date.
func (s *SQLSchedule) DeleteSchedule(channelname string, scheduledate string) error {

	_, err := s.DB.Exec(`DELETE FROM ccms_schedule WHERE channel_name = $1 and schedule= $2`, channelname, scheduledate)
	if err != nil {
		return err
	}

	return err
}
