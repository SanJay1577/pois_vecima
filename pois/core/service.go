package core

import (
	"fmt"
	serviceLog "pois/config/logging"
	sh "pois/shared"
	"strconv"
	"strings"
)

/*
*
Params : aliasName
Description : Alias have channel return the channelname
Returns : channelname, true -> if exists else false-> no channels for the corresponding channels
*
*/
func GetChannelNameByAliasName(alias string) (string, bool) {
	channelName := sh.AliasChannelMap.GetChannel(alias)
	if channelName == "" {
		return "", false
	}
	return channelName, true
}

/*
Input Argument : ChannelName , Date and Time
Return Type : Duration, SchedulerPath
*/
func GetScheduleByChannelAndTime(channelName string, date string, time string) string {

	ch := sh.ChannelSchedules{Channel: channelName, Date: date}
	serviceLog.EsamLog.Infof("[POST] channel schedule map size [%v] ", len(sh.ChannelScheduleData.GetAllSchedules()))
	serviceLog.EsamLog.Infof("[POST] channel schedule map value [%v] ", sh.ChannelScheduleData.GetAllSchedules())

	if timeDurationMap, ok := sh.ChannelScheduleData.GetSchedule(ch); ok {
		if schedules, ok := timeDurationMap[time]; ok {

			if len(schedules.Schedules) > 1 {
				var duration string
				for index, schedule := range schedules.Schedules {
					if index != 0 {
						duration = ProcessDuration(duration, schedule.ScheduledLength)
					} else {
						duration = schedule.ScheduledLength
					}
				}
				return duration

			} else {
				return schedules.Schedules[0].ScheduledLength
			}

		} else {
			serviceLog.EsamLog.Infof("[POST] no schedule found for the time %v", time)
			return ""
		}
	} else {
		serviceLog.EsamLog.Infof("[POST] no schedule found for the day %v", date)
		return ""
	}

}

// ProcessDuration will extract and calculate schedule duration and return it as a string
func ProcessDuration(duration1 string, duration2 string) string {

	hhs, mms, sss, err1 := ExtractDuration(duration1)
	if err1 != "" {
		return "noTimeForADay"
	}

	hh2, mm2, ss2, err2 := ExtractDuration(duration2)
	if err2 != "" {
		return "noTimeForADay"
	}

	durationInSec := ((hhs * 3600) + (mms * 60) + (sss * 1)) + ((hh2 * 3600) + (mm2 * 60) + (ss2 * 1))

	hhs = durationInSec / 3600
	if durationInSec%3600 > 0 {
		mms = (durationInSec - (hhs * 3600)) / 60
		if (durationInSec-(hhs*3600))%60 > 0 {
			sss = (durationInSec - (hhs * 3600) - (mms * 60))
		} else {
			sss = 0
		}
	} else {
		mms = 0
		sss = 0
	}

	serviceLog.EsamLog.Debugf("[POST] Final Duration calculated %v HH %v MM %v SS", hhs, mms, sss)
	return fmt.Sprintf("%02d", hhs) + fmt.Sprintf("%02d", mms) + fmt.Sprintf("%02d", sss)
}

// ExtractDuration will time from the given string
func ExtractDuration(duration string) (int, int, int, string) {

	durationComponents := strings.Split(duration, "")

	hhs, err := strconv.Atoi(strings.Join(durationComponents[:2], ""))

	if err != nil {
		//executes if there is any error
		serviceLog.EsamLog.Errorf("[POST] string to number parsing error for value %v", strings.Join(durationComponents[:2], ""))
		return 0, 0, 0, "noTimeForADay"
	}

	mms, err := strconv.Atoi(strings.Join(durationComponents[2:4], ""))
	if err != nil {
		//executes if there is any error
		serviceLog.EsamLog.Errorf("[POST] string to number parsing error for value %v", strings.Join(durationComponents[2:4], ""))
		return 0, 0, 0, "noTimeForADay"
	}

	sss, err := strconv.Atoi(strings.Join(durationComponents[4:], ""))
	if err != nil {
		//executes if there is any error
		serviceLog.EsamLog.Errorf("[POST] string to number parsing error for value %v", strings.Join(durationComponents[4:], ""))
		return 0, 0, 0, "noTimeForADay"
	}

	return hhs, mms, sss, ""

}
