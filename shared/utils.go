package shared

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"pois/config"
	serviceLog "pois/config/logging"
	"pois/models"
	"pois/store"

	"git.eng.vecima.com/cloud/golib/v4/httpservice"
)

var Stores *store.DataStore

// Channel schedule
type Schedules struct {
	Schedules []ScheduleInfo
}

var threadProfile = pprof.Lookup("threadcreate")

type ScheduleInfo struct {
	EventType           string
	ScheduledDate       string
	ScheduledTime       string
	WindowStartTime     string
	WindowDurationTime  string
	BreakWithinWindow   string
	PositionWithinBreak string
	ScheduledLength     string
	ActualAiredTime     string
	ActualAiredLength   string
	ActualAiredPosition string
	SpotIdentification  string
	StatusCode          string
	UserDefined         string
}

type ChannelSchedules struct {
	Channel string
	Date    string
}

// Structure defined to get the channelSchedulemap
type ChannelScheduleMaps struct {
	channelScheduleMap map[ChannelSchedules]map[string]Schedules
	channelLock        sync.Mutex
	expireAtTimestamp  int64
	stop               chan struct{}
	wg                 sync.WaitGroup
}

// Structure to form failure json api response
type ResponseMessage struct {
	Message string `json:"message"`
}

func (channel *ChannelScheduleMaps) StopCleanup() {
	close(channel.stop)
	channel.wg.Wait()
}

// The function initialize the cleauploop to check the schedule expiration time
// and delete the schedule.
func InitailizeCleanUp(cleanupInterval time.Duration) {
	ChannelScheduleData.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		ChannelScheduleData.CleanScheduleData(cleanupInterval)
		ChannelScheduleData.wg.Done()
	}(cleanupInterval)
}

// The Below function is responsible of the cleaning the schedule information based on the expiration time
func (channel *ChannelScheduleMaps) CleanScheduleData(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()
	for {
		select {
		case <-channel.stop:
			return
		case <-timer.C:
			for key, _ := range channel.channelScheduleMap {
				channel.channelLock.Lock()
				if channel.expireAtTimestamp <= time.Now().Unix() {
					delete(channel.channelScheduleMap, key)
				}
				channel.channelLock.Unlock()
			}

		}
	}
}

func InternalServerErrorResponse() *httpservice.HttpServiceResponse {

	response := &httpservice.HttpServiceResponse{
		Status:      http.StatusInternalServerError,
		ContentType: "text/plain",
		Header:      nil,
		Body:        []byte("Internal server error"),
	}
	return response
}

// initializing channel schedules map data structure
var ChannelScheduleData = ChannelScheduleMaps{channelScheduleMap: make(map[ChannelSchedules]map[string]Schedules), stop: make(chan struct{})}

//methods to get locking feature for schedule data

// Delete a particular channel
func (channel *ChannelScheduleMaps) DelSchedule(key ChannelSchedules) {
	// Exculisve lock for schedule storage
	channel.channelLock.Lock()
	delete(channel.channelScheduleMap, key)
	// Exculisve lock for schedule storage is un locked
	defer channel.channelLock.Unlock()
}

// Set schedules for a particular channel
func (channel *ChannelScheduleMaps) SetSchedule(key ChannelSchedules, value map[string]Schedules) {
	// Exculisve lock for schedule storage

	channel.channelLock.Lock()
	channel.channelScheduleMap[key] = value
	channel.expireAtTimestamp = time.Now().Unix() + int64(config.GetConfig().GetInt("api.responseTimeout"))
	// Exculisve lock for schedule storage is un locked
	defer channel.channelLock.Unlock()
}

// Fetch schedules for a particular channel
func (channel *ChannelScheduleMaps) GetSchedule(key ChannelSchedules) (map[string]Schedules, bool) {
	// Exculisve lock for schedule storage
	channel.channelLock.Lock()
	// Exculisve lock for schedule storage is un locked
	defer channel.channelLock.Unlock()
	//returning channelScheduleMap
	schedule, ok := channel.channelScheduleMap[key]
	return schedule, ok
}

// Fetch all schedules
func (channel *ChannelScheduleMaps) GetAllSchedules() map[ChannelSchedules]map[string]Schedules {
	// Exculisve lock for schedule storage
	channel.channelLock.Lock()
	// Exculisve lock for schedule storage is un locked
	defer channel.channelLock.Unlock()
	// returning all the map for loops
	return channel.channelScheduleMap
}

// structure defined with exclusive lock for ccms channel names
type AliasChannels struct {
	aliasForChannel map[string]string
	dbLock          sync.Mutex
}

var AliasChannelMap = AliasChannels{aliasForChannel: make(map[string]string)}

// SetChannel adds new channel to the alias name with key and values
func (alias *AliasChannels) SetChannel(key string, value string) {
	//obtain an exclusive lock
	alias.dbLock.Lock()
	//set the key/value in the map
	alias.aliasForChannel[key] = value
	//release the exclusive lock
	defer alias.dbLock.Unlock()
}

// Get the channel name using alias name with the required key
func (alias *AliasChannels) GetChannel(key string) string {
	// obtaining the exculsive lock feature
	alias.dbLock.Lock()
	// after returing the channel unlcok the exclusive lock
	defer alias.dbLock.Unlock()
	//  get the channel from the given key
	return alias.aliasForChannel[key]
}

// Getting all the channel informations for range and loop through
func (alias *AliasChannels) GetAllChannels() map[string]string {
	// obtaining the exculsive lock feature for channel data structure
	alias.dbLock.Lock()
	// after returing all the channels unlcok the exclusive lock
	defer alias.dbLock.Unlock()
	return alias.aliasForChannel
}

// Delete a channel name using aliasn name method
func (alias *AliasChannels) DelChannel(key string) {
	//obtain an exclusive lock
	alias.dbLock.Lock()
	delete(alias.aliasForChannel, key)
	//release the exclusive lock
	defer alias.dbLock.Unlock()
}

/*
extract the path variables from a given request url, returns the path variables in a slices of string

path : is the requested URL Path
base path : API ROOT Path example /pois/
rootPath : is the alias and ccms root path

//Extract the  path params returns the params in slice
//Example : /pois/v1/channels/cnn/06122022 -> returns the pathVaribles [cnn 06122022]
*/
func ExtractPathVariable(path string, basePath string, vmin int, vmax int, rootPath string) ([]string, int, error) {
	var err error
	version := vmin
	pathVariables := make([]string, 0, 0)

	if path != "" {
		path = strings.TrimPrefix(path, basePath)
		pathComponents := strings.Split(path, "/")

		if len(pathComponents) < 2 {
			return pathVariables, version, fmt.Errorf("invalid URI")
		}

		// The first component is a version number (e.g. "v1").
		if strings.HasPrefix(pathComponents[0], "v") {
			versionString := strings.TrimPrefix(pathComponents[0], "v")
			version, err = strconv.Atoi(versionString)
		}

		if err != nil || version < vmin || version > vmax {
			if err == nil {
				err = fmt.Errorf("invalid or unsupported version path: '%s'", pathComponents[0])
			}
			return pathVariables, version, err
		}

		//append remaining components
		pathComponents = pathComponents[1:]
		path = strings.Join(pathComponents, "/")
		path = strings.TrimPrefix(strings.Trim(path, ""), rootPath+"/")

		if path != "" {
			pathVariables = strings.Split(path, "/")
			return pathVariables, version, nil
		}

	}
	return pathVariables, version, nil
}

// this function checks the provided file exists
// return true if exists, otherwise return false
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/*
Validating the channel name, channel name accepts the alphanumeric and special character(a-zA-Z0-9\.\-\|!@#$%&_><+=')
*/
func ValidateChannelName(channelName string) bool {

	//channel name should accepts alphanumeric and special character
	IsValidChannelName := regexp.MustCompile(`^[a-zA-Z0-9\.\-\|!@#$%&_><+=']*$`).MatchString(channelName)
	if IsValidChannelName {
		return true
	} else {
		return false
	}
}

// This function compares the current system date with provided data.
// returns true only current system date must today's date or future date
// if it's previous date returns false won't allow to create a schedules for the given date with the provided channel.
func CompareDate(day string, month string, year string) bool {

	modifiedDate := month + "/" + day + "/" + year

	//get the current system date for compare
	currentSystemDate := time.Now()
	formatedCurrentDate := currentSystemDate.Format("01022006")

	// channels allow to create schedule on the same date as current or future
	return modifiedDate <= formatedCurrentDate

}

// This function validate the date returns day, month, year and true if the date is valid
func ValidateDate(date string) (string, string, string, bool) {

	splice_string := strings.Split(date, "")

	if len(splice_string) == 8 {

		month := splice_string[0:2]
		day := splice_string[2:4]
		year := splice_string[4:]

		dayStr := strings.Join(day, "")
		monthStr := strings.Join(month, "")
		yearStr := strings.Join(year, "")

		stringDate := monthStr + "/" + dayStr + "/" + yearStr

		_, err := time.Parse("01/02/2006", stringDate)
		if err != nil {
			fmt.Println("Error parsing time", err)
			return "", "", "", false
		} else {
			return dayStr, monthStr, yearStr, true
		}

	}
	return "", "", "", false
}

// Validate the schedular components that satisfy the given condition. if any of the components fails it returns false
func ValidateScheduleFileds(schedulerComponents []string) bool {
	//scheduled Time field have 6bytes length
	if len([]byte(schedulerComponents[2])) != 6 {
		return false
	}

	// window start time have 4byte length
	if len([]byte(schedulerComponents[3])) != 4 {
		return false
	}

	//Window Duration Time - 4byte length
	if len([]byte(schedulerComponents[4])) != 4 {
		return false
	}
	//Break Number with in Window - 4 byte length
	if len([]byte(schedulerComponents[5])) != 3 {
		return false
	}

	// Position Number With In Break - 4 byte length
	if len([]byte(schedulerComponents[6])) != 3 {
		return false
	}
	// Scheduled Length - 6 byte length
	if len([]byte(schedulerComponents[7])) != 6 {
		return false
	}

	//Actual Aired Time - 6 byte length
	if len([]byte(schedulerComponents[8])) != 6 {
		return false
	}

	//Actual Aired Length - 8 byte length
	if len([]byte(schedulerComponents[9])) != 8 {
		return false
	}

	//Actual Aire Position With in Break - 3byte length
	if len([]byte(schedulerComponents[10])) != 3 {
		return false
	}
	// Spot Identification - Bytes 62 - 72bytes

	if len([]byte(schedulerComponents[11])) != 11 {
		return false
	}

	//Status Code  - 4 byte length
	if len([]byte(schedulerComponents[12])) != 4 {
		return false
	}

	return true
}

/*
This function preprocess the schedular file store the schedules in the memory along with channel name
*/
func PreprocessSchedulerFile(scheduleContent string, channelName string, dateMonth string) (bool, int, int) {

	validScheduleCount := 0
	invalidSchdeuleCount := 0

	serviceLog.CcmsLog.Infof("[%v] TID:[%v] ChannelName:[%v] Date:[%v]", "PUT", "", channelName, dateMonth)

	var scheduleMap = make(map[string]Schedules)

	AliasChannelMap.SetChannel(channelName, channelName)

	ch := ChannelSchedules{Channel: channelName, Date: dateMonth}
	ChannelScheduleData.DelSchedule(ch)
	DeleteSchedule(channelName, dateMonth)

	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key [%v] ", "PUT", "", ch)

	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] ChannelName:[%v] Date:[%v] Delete old records", "PUT", "", channelName, dateMonth)

	for _, line := range strings.Split(strings.TrimSuffix(scheduleContent, "\n"), "\n") {
		//length of the each schedule at least 77 byte

		if len([]byte(line)) >= 77 {
			//Read the content line by line
			schedulerComponents := strings.Split(line, " ")
			if len(schedulerComponents) >= 13 {
				//first field is Match With LOI

				if schedulerComponents[0] != "LOI" || len([]byte(schedulerComponents[1])) != 4 {
					// Debug Statement for line miss
					serviceLog.CcmsLog.Errorf("[%v] TID:[%v] key:[%v] Not a 'LOI' event type - scheduleInfos:[%v] ", "PUT", "", ch, line)
					invalidSchdeuleCount += 1
					continue
				}

				//check the schedule date is match with the current provided data
				if schedulerComponents[1] != dateMonth {
					serviceLog.CcmsLog.Errorf("[%v] TID:[%v] key:[%v] Not a valid schedule date - scheduleInfos:[%v] ", "PUT", "", ch, line)
					invalidSchdeuleCount += 1
					continue
				}

				ch.Date = schedulerComponents[1]

				//2nd param - Scheduled Date Bytes 5-8
				if ValidateScheduleFileds(schedulerComponents) {

					var scheduleInfos ScheduleInfo
					scheduleInfos.EventType = schedulerComponents[0]
					scheduleInfos.ScheduledDate = schedulerComponents[1]
					scheduleInfos.ScheduledTime = schedulerComponents[2]
					scheduleInfos.WindowStartTime = schedulerComponents[3]
					scheduleInfos.WindowDurationTime = schedulerComponents[4]
					scheduleInfos.BreakWithinWindow = schedulerComponents[5]
					scheduleInfos.PositionWithinBreak = schedulerComponents[6]
					scheduleInfos.ScheduledLength = schedulerComponents[7]
					scheduleInfos.ActualAiredTime = schedulerComponents[8]
					scheduleInfos.ActualAiredLength = schedulerComponents[9]
					scheduleInfos.ActualAiredPosition = schedulerComponents[10]
					scheduleInfos.SpotIdentification = schedulerComponents[11]
					scheduleInfos.StatusCode = schedulerComponents[12]
					scheduleInfos.UserDefined = strings.Trim(schedulerComponents[13], "\r")

					var schdeulemodel models.Schedule
					schdeulemodel.Channelname = channelName
					schdeulemodel.Schedule = dateMonth
					schdeulemodel.Eventtype = schedulerComponents[0]
					schdeulemodel.Scheduleddate = schedulerComponents[1]
					schdeulemodel.Scheduledtime = schedulerComponents[2]
					schdeulemodel.Windowstarttime = schedulerComponents[3]
					schdeulemodel.Windowdurationtime = schedulerComponents[4]
					schdeulemodel.Breakwithinwindow = schedulerComponents[5]
					schdeulemodel.Positionwithinbreak = schedulerComponents[6]
					schdeulemodel.Scheduledlength = schedulerComponents[7]
					schdeulemodel.Actualairedtime = schedulerComponents[8]
					schdeulemodel.Actualairedlength = schedulerComponents[9]
					schdeulemodel.Actualairedposition = schedulerComponents[10]
					schdeulemodel.Spotidentification = schedulerComponents[11]
					schdeulemodel.Statuscode = schedulerComponents[12]
					schdeulemodel.Userdefined = strings.Trim(schedulerComponents[13], "\r")

					err := CreateSchedule(schdeulemodel)

					if err != nil {
						serviceLog.CcmsLog.Errorf("[%v] TID:[%v] key:[%v] Error inserting schedule into DB - scheduleInfos:[%v] Error: [%v]", "PUT", "", ch, line, err)
						invalidSchdeuleCount += 1
						continue
					} else {
						//scheduleStored Count
						validScheduleCount += 1

						serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] scheduleInfos:[%v] ", "PUT", "", ch, scheduleInfos)

						if schedules, ok := scheduleMap[schedulerComponents[2]]; ok {
							//schedules.Schedules = append(schedules.Schedules, scheduleInfos)
							shInfo := schedules.AddNewSchedule(scheduleInfos)
							serviceLog.CcmsLog.Debugf("[%v] TID:[%v] appending new schedule", "PUT", "")
							serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] shInfo:[%v] ", "PUT", "", ch, shInfo)
							serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] schedules:[%v] ", "PUT", "", ch, schedules)
							scheduleMap[schedulerComponents[2]] = schedules
						} else {
							shInfo := []ScheduleInfo{}
							schedules := Schedules{shInfo}
							schedules.AddNewSchedule(scheduleInfos)
							serviceLog.CcmsLog.Debugf("[%v] TID:[%v] new schedule", "PUT", "")
							serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] shInfo:[%v] ", "PUT", "", ch, shInfo)
							serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] schedules:[%v] ", "PUT", "", ch, schedules)
							scheduleMap[schedulerComponents[2]] = schedules
						}
					}

				} else {
					serviceLog.CcmsLog.Errorf("[%v] TID:[%v] key:[%v] Schedule have less than 77 bytes - scheduleInfos:[%v] ", "PUT", "", ch, line)
					invalidSchdeuleCount += 1
					continue
				}
			} else {
				serviceLog.CcmsLog.Errorf("[%v] TID:[%v] key:[%v] Schedule have less than 13 fields - scheduleInfos:[%v] ", "PUT", "", ch, line)
				invalidSchdeuleCount += 1
				continue
			}
		} else {
			if strings.HasPrefix(line, "REM") || strings.HasPrefix(line, "END") {
				continue
			} else {
				serviceLog.CcmsLog.Errorf("[%v] TID:[%v] key:[%v] Schedule have less than 77 bytes - scheduleInfos:[%v] ", "PUT", "", ch, line)
				invalidSchdeuleCount += 1
				continue
			}
		}
	}

	ChannelScheduleData.SetSchedule(ch, scheduleMap)

	serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] Map:[%v] ", "PUT", "", ch, ChannelScheduleData.GetAllSchedules())

	serviceLog.CcmsLog.Infof("[%v] TID:[%v] ChannelName:[%v] Date:[%v] Exit from preprocess", "PUT", "", channelName, dateMonth)

	return len(ChannelScheduleData.GetAllSchedules()) > 0, validScheduleCount, invalidSchdeuleCount
}

// Add new schdule to Schedules struct
func (schedule *Schedules) AddNewSchedule(scheduleInfo ScheduleInfo) []ScheduleInfo {
	schedule.Schedules = append(schedule.Schedules, scheduleInfo)
	return schedule.Schedules
}

// CreateSchedule make schedule create a call to the DB and return error if it fails
func CreateSchedule(schedule models.Schedule) error {

	serviceLog.CcmsLog.Infof("Create schedule for channel: %v", schedule.Channelname)

	schedules, err := Stores.CreateSchedule(schedule)
	if err != nil {
		serviceLog.CcmsLog.Errorf("error creating schedule for channel: %v, %v", schedule.Channelname, err)
		return err
	}

	if len(schedules) > 0 {
		serviceLog.CcmsLog.Infof("Query Exec for channel: %v", schedule.Channelname)
	}

	return nil

}

// LoadChannelScheduleAndAlias will load channel alias and current and future date schedule information from DB to in-memory
func LoadChannelScheduleAndAlias() {
	currentSystemDate := time.Now()

	monthdate := fmt.Sprintf("%02d%02d", int(currentSystemDate.Month()), currentSystemDate.Day())

	AddChannelAliasToInMemory()

	AddScheduleToInMemory(monthdate)
}

// AddChannelAliasToInMemory will load channel alias from DB to in-memory
func AddChannelAliasToInMemory() error {

	serviceLog.AliasLog.Infof("get all channel alias from DB")

	alias, err := Stores.FindAllAlias()
	if err != nil {
		serviceLog.AliasLog.Errorf("no channel alias found in DB")
		return err
	} else {
		serviceLog.AliasLog.Infof("Query Exec for channel alias")

		if len(alias) > 0 {

			serviceLog.AliasLog.Infof("channel alias found in DB")

			for _, aliasStr := range alias {
				AliasChannelMap.SetChannel(aliasStr.AliasName, aliasStr.Channelname)
				AliasChannelMap.SetChannel(aliasStr.Channelname, aliasStr.Channelname)
			}

			return nil

		} else {
			serviceLog.AliasLog.Errorf("no channel alias found in DB")
			return nil
		}
	}

}

// AddScheduleToInMemory will load current and future date channel schedule information from DB to in-memory
func AddScheduleToInMemory(scheduleDate string) error {

	serviceLog.CcmsLog.Infof("get schedule from DB for date: %v", scheduleDate)

	schedule, err := Stores.FindScheduleByDate(scheduleDate)
	if err != nil {
		serviceLog.CcmsLog.Errorf("no schedule from DB for date: %v, %v", scheduleDate, err)
		return err
	}
	serviceLog.CcmsLog.Infof("Query Exec for date: %v", scheduleDate)

	if len(schedule) > 0 {

		serviceLog.CcmsLog.Infof("schedule found for date: %v", scheduleDate)

		for _, sch := range schedule {

			ch := ChannelSchedules{Channel: sch.Channelname, Date: sch.Scheduleddate}

			var scheduleInfos ScheduleInfo
			scheduleInfos.EventType = sch.Eventtype
			scheduleInfos.ScheduledDate = sch.Scheduleddate
			scheduleInfos.ScheduledTime = sch.Scheduledtime
			scheduleInfos.WindowStartTime = sch.Windowstarttime
			scheduleInfos.WindowDurationTime = sch.Windowdurationtime
			scheduleInfos.BreakWithinWindow = sch.Breakwithinwindow
			scheduleInfos.PositionWithinBreak = sch.Positionwithinbreak
			scheduleInfos.ScheduledLength = sch.Scheduledlength
			scheduleInfos.ActualAiredTime = sch.Actualairedtime
			scheduleInfos.ActualAiredLength = sch.Actualairedlength
			scheduleInfos.ActualAiredPosition = sch.Actualairedposition
			scheduleInfos.SpotIdentification = sch.Spotidentification
			scheduleInfos.StatusCode = sch.Statuscode
			scheduleInfos.UserDefined = sch.Userdefined

			if scheduleMap, ok := ChannelScheduleData.GetSchedule(ch); ok {

				if schedules, ok := scheduleMap[sch.Scheduledtime]; ok {
					//schedules.Schedules = append(schedules.Schedules, scheduleInfos)
					shInfo := schedules.AddNewSchedule(scheduleInfos)
					serviceLog.CcmsLog.Debugf("[%v] TID:[%v] appending new schedule", "PUT", "")
					serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] shInfo:[%v] ", "PUT", "", ch, shInfo)
					serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] schedules:[%v] ", "PUT", "", ch, schedules)
					scheduleMap[sch.Scheduledtime] = schedules
				} else {
					shInfo := []ScheduleInfo{}
					schedules := Schedules{shInfo}
					schedules.AddNewSchedule(scheduleInfos)
					serviceLog.CcmsLog.Debugf("[%v] TID:[%v] new schedule", "PUT", "")
					serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] shInfo:[%v] ", "PUT", "", ch, shInfo)
					serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] schedules:[%v] ", "PUT", "", ch, schedules)
					scheduleMap[sch.Scheduledtime] = schedules
				}
				ChannelScheduleData.SetSchedule(ch, scheduleMap)
			} else {
				var scheduleMap = make(map[string]Schedules)

				shInfo := []ScheduleInfo{}
				schedules := Schedules{shInfo}
				schedules.AddNewSchedule(scheduleInfos)
				serviceLog.CcmsLog.Debugf("[%v] TID:[%v] new schedule", "PUT", "")
				serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] shInfo:[%v] ", "PUT", "", ch, shInfo)
				serviceLog.CcmsLog.Debugf("[%v] TID:[%v] key:[%v] schedules:[%v] ", "PUT", "", ch, schedules)
				scheduleMap[sch.Scheduledtime] = schedules
				ChannelScheduleData.SetSchedule(ch, scheduleMap)
			}
		}

	} else {
		serviceLog.CcmsLog.Errorf("no schedule from DB for date: %v", scheduleDate)
		return nil
	}
	return nil
}

// FindSchedule will fetch the schedule information based on channel and date from DB
func FindSchedule(channelName string, scheduleDate string) error {
	serviceLog.CcmsLog.Infof("get schedule from DB for channel: %v", channelName)
	schedule, err := Stores.FindSchedule(channelName, scheduleDate)
	if err != nil {
		serviceLog.CcmsLog.Errorf("no schedule from DB for channel: %v, %v", channelName, err)
		return nil
	}
	serviceLog.CcmsLog.Infof("Query Exec for channel: %v", channelName)
	if len(schedule) > 0 {
		serviceLog.CcmsLog.Infof("schedule found for channel: %v", channelName)
	} else {
		serviceLog.CcmsLog.Errorf("no schedule from DB for channel: %v", channelName)
		return nil
	}
	return nil

}

// DeleteSchedule will deletes the channel schedule for the given date
func DeleteSchedule(channelName string, schdeuleDate string) error {

	serviceLog.CcmsLog.Infof("delete schdeule from channel: %v", channelName)

	err := Stores.DeleteSchedule(channelName, schdeuleDate)
	if err != nil {
		serviceLog.CcmsLog.Errorf("error deleting schdeule from channel: %v, %v", channelName, err)
		return err
	}
	serviceLog.CcmsLog.Infof("Query Exec for channel: %v", channelName)
	serviceLog.CcmsLog.Infof("deleted schedule: %v from channel: %v", channelName, schdeuleDate)

	return nil

}

// ConcatinateScheduleInformation to make the cahnnel schedule information into a single line
func ConcatinateScheduleInformation(schedules map[string]Schedules) string {
	var scheduleLine string
	for _, element := range schedules {
		for _, schedule := range element.Schedules {
			scheduleLine += schedule.EventType + " " + schedule.ScheduledDate + " " + schedule.ScheduledTime + " " +
				schedule.WindowStartTime + " " + schedule.WindowDurationTime + " " + schedule.BreakWithinWindow + " " +
				schedule.PositionWithinBreak + " " + schedule.ScheduledLength + " " + schedule.ActualAiredTime + " " +
				schedule.ActualAiredLength + " " + schedule.ActualAiredPosition + " " + schedule.SpotIdentification + " " +
				schedule.StatusCode + " " + schedule.UserDefined + "\n"
		}
	}
	return scheduleLine
}

// GetAlias will fetch the alias name for a channel from DB
func GetAlias(channelName string) error {

	serviceLog.AliasLog.Infof("get Alias from DB for channel: %v", channelName)

	alias, err := Stores.FindAlias(channelName)
	if err != nil {
		serviceLog.AliasLog.Errorf("no Alias from DB for channel: %v, %v", channelName, err)
		return nil
	}
	serviceLog.AliasLog.Infof("Query Exec for channel: %v", channelName)

	if len(alias) > 0 {
		serviceLog.AliasLog.Infof("Alias found for channel: %v", channelName)
	} else {
		serviceLog.AliasLog.Errorf("no Alias from DB for channel: %v", channelName)
		return nil
	}

	return nil

}

// CreateAlias will create a alias name for a channel in DB
func CreateAlias(aliasmodel models.Alias) error {

	serviceLog.AliasLog.Infof("craete Alias for channel: %v", aliasmodel.Channelname)

	alias, err := Stores.CreateAlias(aliasmodel)
	if err != nil {
		serviceLog.AliasLog.Errorf("error creating alias for channel: %v with alias: %v, %v", aliasmodel.Channelname, aliasmodel.AliasName, err)
		return nil
	}
	if len(alias) > 0 {
		serviceLog.AliasLog.Infof("Query Exec for channel: %v", aliasmodel.Channelname)
	}

	return nil

}

// DeleteAlias will delete alias name for a channel in DB
func DeleteAlias(channelName string, aliasName string) error {

	serviceLog.AliasLog.Infof("delete alias from channel: %v", channelName)

	err := Stores.DeleteAlias(channelName, aliasName)
	if err != nil {
		serviceLog.AliasLog.Errorf("error deleting alias from channel: %v alias: %v, %v", channelName, aliasName, err)
		return nil
	}
	serviceLog.AliasLog.Infof("Query Exec for channel: %v", channelName)
	serviceLog.AliasLog.Infof("deleted alias: %v from channel: %v", channelName, aliasName)

	return nil

}
