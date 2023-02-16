package testing

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	//initialize and setup the test api

	if err := initialize(); err != nil {

		fmt.Printf("Unable to initialize testing API : %s", err.Error())
		os.Exit(1)
	}

	//Now run all the tests
	os.Exit(m.Run())
}
