package email

import "os"

func modelTestShouldRun() bool {
	return os.Getenv("TEST_DATABASE") != ""
}
