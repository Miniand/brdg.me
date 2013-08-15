package main

import (
	"bytes"
	"fmt"
	"github.com/Miniand/brdg.me/render"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
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

func GetPlainEmailBody(r io.Reader) (*mail.Message, string, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, "", err
	}
	body, _, err := GetPlainEmailBodyReader(m.Body,
		m.Header.Get("Content-Type"))
	return m, body, err
}

func GetPlainEmailBodyReader(r io.Reader, contentType string) (string, string,
	error) {
	var body, foundContentType string
	if contentType == "" {
		// No content type, assume plain
		rawBody, err := ioutil.ReadAll(r)
		if err != nil {
			return "", "", err
		}
		return string(rawBody), "text/plain", nil
	}
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

func FromAddr() string {
	from := os.Getenv("BRDGME_EMAIL_SERVER_SMTP_FROM")
	if from == "" {
		from = "play@brdg.me"
	}
	return from
}

func SendMailAuth() smtp.Auth {
	return smtp.PlainAuth("", FromAddr(), "password", "mail.brdg.me")
}

func SendRichMail(to []string, subject string, body string,
	extraHeaders []string) error {
	terminalOutput, err := render.RenderTerminal(body)
	if err != nil {
		return err
	}
	htmlOutput, err := render.RenderHtml(body)
	if err != nil {
		return err
	}
	// Make a multipart message
	buf := &bytes.Buffer{}
	data := multipart.NewWriter(buf)
	plainW, err := data.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain"},
	})
	if err != nil {
		return err
	}
	_, err = plainW.Write([]byte(terminalOutput))
	if err != nil {
		return err
	}
	htmlW, err := data.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{`text/html; charset="UTF-8"`},
	})
	if err != nil {
		return err
	}
	_, err = htmlW.Write([]byte(`<pre style="color:#000000;">` + htmlOutput))
	if err != nil {
		return err
	}
	err = data.Close()
	if err != nil {
		return err
	}
	headers := []string{
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s",
			data.Boundary()),
	}
	headers = append(headers, extraHeaders...)
	return SendMail(to,
		fmt.Sprintf("%s\r\n%s", strings.Join(headers, "\r\n"), buf.String()))
}

func SendMail(to []string, data string) error {
	smtpAddr := os.Getenv("BRDGME_EMAIL_SERVER_SMTP_ADDR")
	if smtpAddr == "" {
		smtpAddr = "localhost:25"
	}
	return smtp.SendMail(smtpAddr, SendMailAuth(), FromAddr(), to, []byte(data))
}
