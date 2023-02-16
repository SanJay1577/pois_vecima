package ccms

import (
	"fmt"
	"pois/version"

	serviceLog "pois/config/logging"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
)

var ccmsLogger *zaplogger.Logger

// CCMS Log Initialization
func Initialize() error {
	ccmsLogger = zaplogger.NewDomainLogger("ccms")

	ccmsLogger.Infof("ccms initializing.")

	//define service level log initialization implementation

	//Initialize ccms transactional Log level
	if err := serviceLog.CcmsLogInitialization("ccms"); err != nil {
		return fmt.Errorf("%s - Alias Log Initialization problem: %s", version.Version.ApplicationName(), err.Error())
	}

	ccmsLogger.Infof("ccms initialization complete.")
	return nil
}
