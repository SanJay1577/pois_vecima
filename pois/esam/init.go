package esam

import (
	"fmt"
	serviceLog "pois/config/logging"
	"pois/version"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
)

var esamLogger *zaplogger.Logger

// All global esam values are registerd in the below code.
var ResponseType string
var ResponseValidation bool                  // Response type from teh config file will be stored here
var DurationValue, DurationValidation string // defining a value for duration value and validations.

func Initialize() error {

	esamLogger = zaplogger.NewDomainLogger("esam")
	esamLogger.Infof("esam initializing.")

	if err := serviceLog.EsamLogInitialization("esam"); err != nil {
		return fmt.Errorf("%s - Esam Log Initialization problem: %s", version.Version.ApplicationName(), err.Error())
	}
	esamLogger.Infof("esam initializing complete.")
	return nil
}
