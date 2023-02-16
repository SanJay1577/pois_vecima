package models

// DB model for ccms_schedule
type Schedule struct {
	Channelname         string `db:"channel_name"`
	Schedule            string `db:"schedule"`
	Eventtype           string `db:"event_type"`
	Scheduleddate       string `db:"scheduled_date"`
	Scheduledtime       string `db:"scheduled_time"`
	Windowstarttime     string `db:"window_start_time"`
	Windowdurationtime  string `db:"window_duration_time"`
	Breakwithinwindow   string `db:"break_within_window"`
	Positionwithinbreak string `db:"position_within_break"`
	Scheduledlength     string `db:"scheduled_length"`
	Actualairedtime     string `db:"actual_aired_time"`
	Actualairedlength   string `db:"actual_aired_length"`
	Actualairedposition string `db:"actual_aired_position"`
	Spotidentification  string `db:"spot_identification"`
	Statuscode          string `db:"status_code"`
	Userdefined         string `db:"user_defined"`
}
