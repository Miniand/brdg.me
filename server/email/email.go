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
	imageOutput, width, height, err := render.RenderImageMeta(bodyWithFooter)
	if err != nil {
		return err
	}
	htmlOutput := fmt.Sprintf(
		`<img width="%d" height="%d" style="min-width:%dpx;min-height:%dpx;" src="cid:game.png@brdg.me" />`,
		width,
		height,
		width,
		height,
	)
	// Make a multipart message
	buf := &bytes.Buffer{}
	data := multipart.NewWriter(buf)
	// Write HTML version
	htmlW, err := data.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{`text/html; charset="UTF-8"`},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return err
	}
	src := []byte(htmlOutput)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	_, err = htmlW.Write(dst)
	if err != nil {
		return err
	}
	// Write image
	imageW, err := data.CreatePart(textproto.MIMEHeader{
		"Content-ID":                []string{"<game.png@brdg.me>"},
		"Content-Type":              []string{"image/png"},
		"Content-Disposition":       []string{"inline"},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return err
	}
	src = []byte(imageOutput)
	dst = make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	_, err = imageW.Write(dst)
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
		fmt.Sprintf("Content-Type: multipart/related; boundary=%s",
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
