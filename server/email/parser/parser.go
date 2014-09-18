package parser

import (
	"regexp"
	"strings"
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

// Search for an email address
func ParseFrom(from string) string {
	reg := regexp.MustCompile(EmailSearchRegexString())
	return strings.ToLower(reg.FindString(from))
}

// Search for a BSON objectid to match to a game (length 24 hex string)
func ParseSubject(subject string) string {
	reg := regexp.MustCompile("\\b[a-f0-9]{8}-([a-f0-9]{4}-){3}.{12}\\b")
	return reg.FindString(subject)
}

// Find contiguous lines as commands until the first blank line
func ParseBody(body string) string {
	return strings.Replace(strings.Replace(body, "\r\n", "\n", -1),
		"\r", "\n", -1)
}
