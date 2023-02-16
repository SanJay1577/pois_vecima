package api

import (
	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
)

var apiLogger *zaplogger.Logger

func Initialize() error {
	apiLogger = zaplogger.NewDomainLogger("api")

	apiLogger.Infof("api initializing.")

	ready := make(chan error)

	// Start up the Wew Service.
	go listenAndServeRequests(ready)
	err := <-ready
	if err == nil {
		apiLogger.Infof("api initialization complete. %v", err)
	}

	return err
}
