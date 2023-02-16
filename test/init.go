package testing

import (
	"fmt"
	"pois/app"
)

//initialize ccms and alias API and external server for testing

func initialize() error {

	//initilalize the ccms and alias api
	if err := app.Initialize("cfg/"); err != nil {
		return fmt.Errorf("failed to initialize the CMC API: %s", err.Error())
	}
	return nil
}
