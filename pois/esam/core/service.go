package core

import (
	"fmt"
	serviceLog "pois/config/logging"
	"strconv"
	"strings"
	"time"
)

func UtcFormater(utcPoint string) (string, string) {
	timeFormat, err := time.Parse(time.RFC3339, utcPoint)
	if err != nil {
		serviceLog.EsamLog.Infof("[POSt] Error parsing utcpoint: ", err)
	}
	month := fmt.Sprintf("%02d", int(timeFormat.Month()))
	date := fmt.Sprintf("%02d", timeFormat.Day())
	hours := fmt.Sprintf("%02d", timeFormat.Hour())
	minutes := fmt.Sprintf("%02d", timeFormat.Minute())
	seconds := fmt.Sprintf("%02d", timeFormat.Second())
	dateMonth := fmt.Sprintf("%v%v", month, date)
	timeSeconds := fmt.Sprintf("%v%v%v", hours, minutes, seconds)
	return dateMonth, timeSeconds
}

// The below set of function will convert the type of duration that we received from ccms
// The duration value will be convert into PT00H00M30S format
// and will be passed to the replace response field.
func DurationTypeConverstion(duration string) string {
	serviceLog.EsamLog.Infof("[POST] Duration %v ", duration)
	if duration == "" {
		return ""
	}
	hh, mm, ss := ExtractDuration(duration)
	modifiedDuration := "PT" + strconv.Itoa(hh) + "H" + strconv.Itoa(mm) + "M" + strconv.Itoa(ss) + "S"
	return modifiedDuration
}

// This function will extract the duration from CCMS and will return three diffrent interger values
// such as hours minutes and seconds in hh mm ss format
// The return values will be used in the durationTypeConverstion function
func ExtractDuration(duration string) (int, int, int) {
	// hours minutes and seconds variables are declared and initiated to return from the function
	durationComponents := strings.Split(duration, "")
	hoursValue, err := strconv.Atoi(strings.Join(durationComponents[:2], ""))
	if err != nil {
		serviceLog.EsamLog.Errorf("[POST] string to number parsing error for value %v Error: %v", strings.Join(durationComponents[:2], ""), err)
		return 0, 0, 0
	}
	minutesValue, err := strconv.Atoi(strings.Join(durationComponents[2:4], ""))
	if err != nil {
		serviceLog.EsamLog.Errorf("[POST] string to number parsing error for value %v Error: %v", strings.Join(durationComponents[2:4], ""), err)
		return 0, 0, 0
	}
	secondsValue, err := strconv.Atoi(strings.Join(durationComponents[4:], ""))
	if err != nil {
		serviceLog.EsamLog.Errorf("[POST] string to number parsing error for value %v Error: %v", strings.Join(durationComponents[4:], ""), err)
		return 0, 0, 0
	}
	return hoursValue, minutesValue, secondsValue
}
