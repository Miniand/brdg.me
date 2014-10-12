package parser

import (
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"regexp"
	"strings"
)

func GetPlainEmailBody(r io.Reader) (*mail.Message, string, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, "", err
	}
	body, _, err := GetPlainEmailBodyReader(m.Body,
		m.Header.Get("Content-Type"),
		m.Header.Get("Content-Transfer-Encoding"))
	return m, body, err
}

func GetPlainEmailBodyReader(r io.Reader, contentType string,
	contentTransferEncoding string) (string, string, error) {
	var body, foundContentType string
	// Extract body
	if contentType == "" {
		// No content type, assume plain
		rawBody, err := ioutil.ReadAll(r)
		if err != nil {
			return "", "", err
		}
		body = string(rawBody)
		foundContentType = "text/plain"
	} else {
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
					part.Header.Get("Content-Type"),
					part.Header.Get("Content-Transfer-Encoding"))
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
	}
	// Convert based on content transfer encoding
	switch contentTransferEncoding {
	case "quoted-printable":
		body = DecodeQuotedPrintable(body)
	case "base64":
		body = DecodeBase64(body)
	}
	return body, foundContentType, nil
}

func DecodeQuotedPrintable(body string) string {
	return regexp.MustCompile(`=[0-9A-F]{2}`).ReplaceAllStringFunc(
		regexp.MustCompile(`=\r\n`).ReplaceAllString(body, ""),
		func(repl string) string {
			b, err := hex.DecodeString(repl[1:])
			if err != nil {
				panic(err.Error())
			}
			return string(b)
		})
}

func DecodeBase64(s string) string {
	output, _ := base64.StdEncoding.DecodeString(s)
	return string(output)
}

func StripHtmlTags(in string) string {
	return regexp.MustCompile(TagMatchRegexString()).ReplaceAllString(in, "")
}

func TagMatchRegexString() string {
	return "<[^>]*>"
}
