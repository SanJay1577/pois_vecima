package alias

import (
	"fmt"

	serviceLog "pois/config/logging"
	version "pois/version"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"
)

var aliasLogger *zaplogger.Logger

// Alias Log Initialization
func Initialize() error {

	aliasLogger = zaplogger.NewDomainLogger("alias")

	aliasLogger.Infof("alias initializing.")

	//define service level log initialization implementation

	//Initialize Alias Transcational Log Level
	if err := serviceLog.AliasLogInitialization("alias"); err != nil {
		return fmt.Errorf("%s - Alias Log Initialization problem: %s", version.Version.ApplicationName(), err.Error())
	}

	aliasLogger.Infof("alias initialization complete.")

	return nil
}
