package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

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
	bodyWithFooter := fmt.Sprintf(
		"%s\n\n\n{{c \"gray\"}}To no longer receive emails or game invites, please reply with {{b}}unsubscribe{{_b}}.{{_c}}",
		body)
	terminalOutput, err := render.RenderTerminal(bodyWithFooter)
	if err != nil {
		return err
	}
	htmlOutput, err := render.RenderHtml(bodyWithFooter)
	if err != nil {
		return err
	}
	// Make a multipart message
	buf := &bytes.Buffer{}
	data := multipart.NewWriter(buf)
	// Write plain version
	plainW, err := data.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"text/plain"},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return err
	}
	src := []byte(terminalOutput)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	_, err = plainW.Write(dst)
	if err != nil {
		return err
	}
	// Write HTML version
	htmlW, err := data.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{`text/html; charset="UTF-8"`},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return err
	}
	src = []byte(fmt.Sprintf(
		`<pre style="color:white;background-color:#1d1f21;padding:13px;font-size:13px;line-height:15px;font-family:DejaVu Sans Mono,monospace,Segoe UI Symbol;white-space:pre-wrap;">%s`,
		htmlOutput,
	))
	dst = make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	_, err = htmlW.Write(dst)
	if err != nil {
		return err
	}
	// Wrap up and send with headers
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
