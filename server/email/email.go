package main

import (
	"regexp"
)

func EmailRegexString() string {
	return "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,4}"
}

func EmailValidationRegexString() string {
	return "^" + EmailRegexString() + "$"
}

func EmailSearchRegexString() string {
	return "\\b" + EmailRegexString() + "\\b"
}

func ValidateEmail(email string) bool {
	reg := regexp.MustCompile(EmailValidationRegexString())
	return reg.MatchString(email)
}
