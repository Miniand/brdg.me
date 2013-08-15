package main

import (
	"testing"
)

func TestMissingUserIsSubscribed(t *testing.T) {
	if modelTestShouldRun() {
		unsubscribed, err := UserIsUnsubscribed("farfdsiahfdufdhas@fdfhdsak.dfisa")
		if err != nil {
			t.Fatal(err)
		}
		if unsubscribed {
			t.Fatal("User was unsubscribed when they shouldn't have been!")
		}
	}
}
