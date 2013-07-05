package main

import (
	"testing"
)

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
