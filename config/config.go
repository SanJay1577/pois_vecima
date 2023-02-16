package config

import (
	"fmt"
	"pois/version"
	"sync"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
	"github.com/spf13/viper"
)

var config *viper.Viper
var cfgLock sync.Mutex

// Default values for the config data is defined
func loadConfig(cfg string) *viper.Viper {
	vp := viper.New()
	vp.SetDefault("app.db.dialect", "postgres")
	vp.SetDefault("app.db.datasource", "dbname=pois user=postgres password=password sslmode=disable")
	vp.SetDefault("app.db.dir", "migrations")
	vp.SetDefault("app.db.table", "migrations")

	vp.SetDefault("app.log.accessLog", true)
	vp.SetDefault("app.log.ccmsLog", true)
	vp.SetDefault("app.log.aliasLog", true)
	vp.SetDefault("app.log.esamLog", true)
	vp.SetDefault("app.log.config.dirName", "log/")
	vp.SetDefault("app.log.config.logToStdout", false)
	vp.SetDefault("app.log.config.logJSON", false)
	vp.SetDefault("app.log.config.level", "info")
	vp.SetDefault("app.log.config.maxSize", 1)
	vp.SetDefault("app.log.config.maxBackups", 3)
	vp.SetDefault("app.log.config.maxAge", 5)

	vp.SetDefault("api.enable", true)
	vp.SetDefault("api.port", 8130)
	vp.SetDefault("api.tlsPort", 8131)
	vp.SetDefault("api.missport", 2244)
	vp.SetDefault("api.tlsCertPath", "./config/localhost.crt")
	vp.SetDefault("api.tlsPrivateKeyPath", "./config/localhost.key")
	vp.SetDefault("api.services", []string{"ccms , alias"})
	vp.SetDefault("api.rootPath", "/pois/")
	vp.SetDefault("api.channelPath", "channels")
	vp.SetDefault("api.aliasPath", "channels/alias")
	vp.SetDefault("api.prometheusRootPath", "/metrics")
	vp.SetDefault("api.swaggerRootPath", "/docs")
	vp.SetDefault("api.responseTimeout", 10)

	vp.SetDefault("esam.enable", true)
	vp.SetDefault("esam.port", 8150)
	vp.SetDefault("esam.tlsport", 8151)
	vp.SetDefault("esam.rootPath", "/esam/v1/*/request")
	vp.SetDefault("esam.providers", []string{"millicom", "discovery", "comcast", "verizon", "altice"})
	vp.SetDefault("esam.responseTimeout", 10)

	vp.SetDefault("esam.response.acquisitionPointIdentity", "delete")
	vp.SetDefault("esam.response.utcPoint", "delete")
	vp.SetDefault("esam.response.signalType", "delete")
	vp.SetDefault("esam.response.streamTimeType", "delete")
	vp.SetDefault("esam.response.scete35Payload", "delete")
	vp.SetDefault("esam.response.scte35DecodeError", "delete")
	vp.SetDefault("esam.response.spliceCommandType", "delete")
	vp.SetDefault("esam.response.spliceEventId", "delete")
	vp.SetDefault("esam.response.spliceDuration", "delete")
	vp.SetDefault("esam.response.noScheduleForTheDay", "delete")
	vp.SetDefault("esam.response.noScheduleForTheTime", "noop")
	vp.SetDefault("esam.response.validRecord", "replace")

	//read all the configurations from config file.
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath(cfg)
	err := vp.ReadInConfig()
	if err != nil {
		fmt.Errorf("%s - Error in reading config file: %s", version.Version.Application(), err.Error())
		return nil
	}

	return vp
}

var ZapConfigData *zaplogger.Config

// Loading the configuration files from the respective config json.

// Initializing all the config file parsing methods.
// chage intilaze logging
func Initialize(cfg string) error {
	if config == nil {
		config = loadConfig(cfg)
	}
	ZapConfigData = &zaplogger.Config{
		DirName:     GetConfig().GetString("app.log.config.dirName"),
		LogToStdout: GetConfig().GetBool("app.log.config.logToStdout"),
		LogJSON:     GetConfig().GetBool("app.log.config.logJSON"),
		Level:       GetConfig().GetString("app.log.config.level"),
		MaxSize:     GetConfig().GetInt("app.log.config.maxSize"),
		MaxBackups:  GetConfig().GetInt("app.log.config.maxBackups"),
		MaxAge:      GetConfig().GetInt("app.log.config.maxAge"),
	}
	return nil
}

// Methods are defined to get all the config values with mutex lock feature.
func GetConfig() *viper.Viper {
	cfgLock.Lock()
	defer cfgLock.Unlock()
	return config
}
