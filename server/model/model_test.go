package model

import (
	"os"
)

func modelTestShouldRun() bool {
	return os.Getenv("TEST_DATABASE") != ""
}

func cleanTestingDatabase() {
	session, err := Connect()
	if err != nil {
		panic(err.Error())
	}
	defer session.Close()
	session.DB(DatabaseName()).DropDatabase()
}
