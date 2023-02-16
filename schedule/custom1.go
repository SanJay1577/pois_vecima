package schedule

import (
	"fmt"
	serviceLog "pois/config/logging"
	cr "pois/pois/core"
	sh "pois/shared"
)

// custom is a custom implementation logic to extract and calculate schedule duration
// and return it as a string
// Will extarct and return the first break position for the given time
func custom(channelName string, date string, time string) (string, string) {

	ch := sh.ChannelSchedules{Channel: channelName, Date: date}

	serviceLog.EsamLog.Infof("[POST] channel schedule map size [%v] ", len(sh.ChannelScheduleData.GetAllSchedules()))
	serviceLog.EsamLog.Infof("[POST] channel schedule map value [%v] ", sh.ChannelScheduleData.GetAllSchedules())

	if timeDurationMap, ok := sh.ChannelScheduleData.GetSchedule(ch); ok {
		if schedules, ok := timeDurationMap[time]; ok {

			if len(schedules.Schedules) > 1 {
				var duration string
				for _, schedule := range schedules.Schedules {
					if schedule.PositionWithinBreak == "001" {
						duration = customProcessDuration(schedule.ScheduledLength, schedule.PositionWithinBreak)
					}
				}
				if duration != "" {
					return duration, ""
				} else {
					return "", "noTimeForADay"
				}
			} else {
				return schedules.Schedules[0].ScheduledLength, ""
			}

		} else {
			serviceLog.EsamLog.Infof("[POST] no schedule found for the time %v", time)
			return "", "noTimeForADay"
		}
	} else {
		serviceLog.EsamLog.Infof("[POST] no schedule found for the day %v", date)
		return "", "noScheduleForADay"
	}

}

// processDuration will extract and calculate schedule duration and return it as a string
func customProcessDuration(duration string, position string) string {

	hhs, mms, sss, err1 := cr.ExtractDuration(duration)
	if err1 != "" {
		return "noTimeForADay"
	}

	durationInSec := ((hhs * 3600) + (mms * 60) + (sss * 1))

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

// init to register provider with its implememtaion logic signatures
func init() {
	registerProvider("millicom", custom)
	registerProvider("comcast", custom)
}
