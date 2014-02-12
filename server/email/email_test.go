package email

import (
	"os"
	"testing"
)

func modelTestShouldRun() bool {
	return os.Getenv("TEST_DATABASE") != ""
}

func TestValidateEmail(t *testing.T) {
	if !ValidateEmail("beefsack@gmail.com") {
		t.Error("beefsack@gmail.com should be valid")
	}
	if ValidateEmail("   fdsafad@egg.com") {
		t.Error("Shouldn't validate with whitespace")
	}
	if ValidateEmail("shonkydonk") {
		t.Error("Shouldn't validate ridiculous emails")
	}
}

func TestDecodeQuotedPrintable(t *testing.T) {
	expected := `If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy.`
	actual := DecodeQuotedPrintable("If you believe that truth=3Dbeauty, then surely =\r\nmathematics is the most beautiful branch of philosophy.")
	if actual != expected {
		t.Fatalf(`Expected: "%s" but got "%s"`, expected, actual)
	}
}
