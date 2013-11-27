package model

import (
	"os"
)

func modelTestShouldRun() bool {
	return os.Getenv("TEST_DB") != ""
}

func cleanTestingDatabase() {
	if err := dropDb(); err != nil {
		panic(err.Error())
	}
	if err := initDb(); err != nil {
		panic(err.Error())
	}
}

func init() {
	os.Setenv("BRDGME_DB_DATABASE", DatabaseName()+"_test")
}
