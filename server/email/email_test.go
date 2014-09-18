package email

import (
	"os"
	"testing"
)

func modelTestShouldRun() bool {
	return os.Getenv("TEST_DATABASE") != ""
}
