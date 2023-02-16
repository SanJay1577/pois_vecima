package api

import (
	"fmt"
	"net/http"

	config "pois/config"
	"pois/pois"

	"git.eng.vecima.com/cloud/golib/v4/httpservice"
	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
)

// This is handler type for the status API.
type apiStatusHandler struct{}

// HandleHttpService is tha application status handler.
func (h *apiStatusHandler) HandleHttpService(request *http.Request, transactionId uint64) (*httpservice.HttpServiceResponse, error) {
	response := &httpservice.HttpServiceResponse{
		Status: http.StatusOK,
		Header: nil,
	}
	return response, nil
}

var statusHandler httpservice.HttpServiceHandler = &apiStatusHandler{}

// ListenAndServerRequests sets up the HTTP request listener
// for the RESTful Web Service.
func listenAndServeRequests(ready chan error) {
	// Fetch the current configuration

	//getting the logging configuation
	logconfig := config.GetConfig().GetBool("app.log.accessLog")

	var accessLog *zaplogger.Logger

	if logconfig {
		var err error
		level, _ := zaplogger.ToLoggingLevel(config.ZapConfigData.Level)
		accessLog, err = zaplogger.NewLogger(&zaplogger.LoggerConfig{

			ID:         "access",
			FileName:   config.ZapConfigData.DirName + "access.log",
			MaxSize:    config.ZapConfigData.MaxSize,
			MaxBackups: config.ZapConfigData.MaxBackups,
			MaxAge:     config.ZapConfigData.MaxAge,
			IsSugar:    false,
			LogType:    zaplogger.LoggerTypeJSON,
			LogLevel:   level,
		})
		if err != nil {
			ready <- fmt.Errorf("unable to create an access log. (%s)", err.Error())
			return
		}
	} else {
		accessLog = zaplogger.NullLogger()
	}

	serviceAPIs := make([]*httpservice.HttpServiceAPI, 0, 0)

	serviceAPIs = append(serviceAPIs,
		&httpservice.HttpServiceAPI{
			ApiRoot:      config.GetConfig().GetString("api.rootpath") + "status/",
			ReadACL:      nil,
			WriteACL:     nil,
			CheckHeaders: false,
			Handler:      statusHandler,
			WebService:   nil,
		})
	poisServiceAPIs := pois.GetServiceAPIs(nil, "/")
	serviceAPIs = append(serviceAPIs, poisServiceAPIs...)

	// Debug
	for i, api := range serviceAPIs {
		apiLogger.Debugf("serviceAPI(%d): %v", i, api.ApiRoot)

	}

	// Creating the HTTP service.
	service := &httpservice.HttpService{
		Interfaces:         []string{},
		TlsInterfaces:      []string{},
		Port:               uint16(config.GetConfig().GetInt("api.port")),
		TlsPort:            uint16(config.GetConfig().GetInt("api.tlsport")),
		TlsCertificateFile: config.GetConfig().GetString("api.tlsCertPath"),
		TlsKeyFile:         config.GetConfig().GetString("api.tlsPrivateKeyPath"),
		API:                serviceAPIs,
		AccessLogger:       accessLog,
	}

	if err := httpservice.New(service); err != nil {
		ready <- fmt.Errorf("unable to create HTTP service: %s", err.Error())
		return
	}

	ready <- nil
}
