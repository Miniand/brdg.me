package main

import (
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
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

func TagMatchRegexString() string {
	return "<[^>]*>"
}

func StripHtmlTags(in string) string {
	return regexp.MustCompile(TagMatchRegexString()).ReplaceAllString(in, "")
}

func GetPlainEmailBody(r io.Reader) (string, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return "", err
	}
	body, _, err := GetPlainEmailBodyReader(m.Body,
		m.Header.Get("Content-Type"))
	return body, err
}

func GetPlainEmailBodyReader(r io.Reader, contentType string) (string, string,
	error) {
	var body, foundContentType string
	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", "", err
	}
	foundContentType = mediatype
	if mediatype == "text/plain" {
		rawBody, err := ioutil.ReadAll(r)
		if err != nil {
			return "", "", err
		}
		body = string(rawBody)
	} else if mediatype == "text/html" {
		rawBody, err := ioutil.ReadAll(r)
		if err != nil {
			return "", "", err
		}
		body = StripHtmlTags(string(rawBody))
	} else if strings.Contains(mediatype, "multipart") &&
		params["boundary"] != "" {
		// Recurse parts
		mpr := multipart.NewReader(r, params["boundary"])
		for {
			part, err := mpr.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return "", "", err
				}
			}
			pBody, pContentType, err := GetPlainEmailBodyReader(part,
				part.Header.Get("Content-Type"))
			if err != nil {
				return "", "", err
			}
			if pContentType == "text/plain" {
				body = pBody
				break
			} else if pBody != "" {
				body = pBody
			}
		}
	}
	return body, foundContentType, nil
}
