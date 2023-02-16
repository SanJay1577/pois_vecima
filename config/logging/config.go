package logging

import (
	"pois/config"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
)

const (
	DefaultLogDirName = "/var/log/pois/"
	DefaultLogLevel   = zaplogger.LoggingLevelInfo
	LogMaxSize        = 100
	LogMaxBackups     = 0
	LogMaxAge         = 30
)

// Defined globally so they can be used by all packages in the project
var CcmsLog *zaplogger.Logger
var AliasLog *zaplogger.Logger
var EsamLog *zaplogger.Logger

// CCMS log Transcational Log Initialization
func CcmsLogInitialization(logId string) error {
	var err error

	level, _ := zaplogger.ToLoggingLevel(config.ZapConfigData.Level)
	//if the log initialization is true
	if config.GetConfig().GetBool("app.log.aliasLog") {
		CcmsLog, err = zaplogger.NewLogger(&zaplogger.LoggerConfig{

			ID:         "ccms",
			FileName:   config.ZapConfigData.DirName + logId + ".log",
			MaxSize:    config.ZapConfigData.MaxSize,
			MaxBackups: config.ZapConfigData.MaxBackups,
			MaxAge:     config.ZapConfigData.MaxAge,
			IsSugar:    true,
			LogType:    zaplogger.LoggerTypeConsole,
			LogLevel:   level,
		})
		if err != nil {
			return err
		}
	} else {
		CcmsLog = zaplogger.NullLogger()
	}

	return err
}

// CCMS log Transcational Log Initialization
func AliasLogInitialization(logId string) error {
	var err error
	level, _ := zaplogger.ToLoggingLevel(config.ZapConfigData.Level)
	if config.GetConfig().GetBool("app.log.aliasLog") {
		AliasLog, err = zaplogger.NewLogger(&zaplogger.LoggerConfig{

			ID:         "alias",
			FileName:   config.ZapConfigData.DirName + logId + ".log",
			MaxSize:    config.ZapConfigData.MaxSize,
			MaxBackups: config.ZapConfigData.MaxBackups,
			MaxAge:     config.ZapConfigData.MaxAge,
			IsSugar:    true,
			LogType:    zaplogger.LoggerTypeConsole,
			LogLevel:   level,
		})

		if err != nil {
			return err
		}
	} else {

		//default null logger
		AliasLog = zaplogger.NullLogger()
	}

	return err
}

// Esam Transcational Log Initialization
func EsamLogInitialization(logId string) error {
	var err error
	level, _ := zaplogger.ToLoggingLevel(config.ZapConfigData.Level)
	if config.GetConfig().GetBool("app.log.esamLog") {
		EsamLog, err = zaplogger.NewLogger(&zaplogger.LoggerConfig{
			ID:         "esam",
			FileName:   config.ZapConfigData.DirName + logId + ".log",
			MaxSize:    config.ZapConfigData.MaxSize,
			MaxBackups: config.ZapConfigData.MaxBackups,
			MaxAge:     config.ZapConfigData.MaxAge,
			IsSugar:    true,
			LogType:    zaplogger.LoggerTypeConsole,
			LogLevel:   level,
		})

		if err != nil {
			return err
		}
	} else {
		EsamLog = zaplogger.NullLogger()
	}

	return err
}

func logSync() {
	CcmsLog.Sync()
	AliasLog.Sync()
	EsamLog.Sync()
}
