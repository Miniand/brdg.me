package scommand

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	comm "github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/server/email"
	"github.com/Miniand/brdg.me/server/model"
)

type DumpCommand struct {
	gameModel *model.GameModel
}

func (c DumpCommand) Parse(input string) []string {
	return comm.ParseNamedCommandNArgs("dump", 0, input)
}

func (c DumpCommand) CanCall(player string, context interface{}) bool {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || !ok {
		return false
	}
	return u.Admin
}

func (c DumpCommand) Call(player string, context interface{},
	args []string) (string, error) {
	if !c.CanCall(player, context) {
		return "", errors.New("you aren't allowed to do that")
	}

	buf := &bytes.Buffer{}
	data := multipart.NewWriter(buf)
	f, err := data.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"application/octet-stream"},
		"Content-Disposition":       {"attachment;filename=.game"},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return "", fmt.Errorf("Unable to create file part: %v", err)
	}
	wr := bufio.NewWriter(f)
	raw := bytes.NewBufferString(fmt.Sprintf("%s\n", c.gameModel.Type))
	raw.Write(c.gameModel.State)
	encoded := base64.StdEncoding.EncodeToString(raw.Bytes())
	wr.WriteString(encoded)
	wr.Flush()

	if err := data.Close(); err != nil {
		return "", fmt.Errorf("Unable to close multipart writer: %v", err)
	}
	headers := []string{
		fmt.Sprintf("Subject: brdg.me game dump of %s", c.gameModel.Id),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s",
			data.Boundary()),
	}
	if err := email.SendMail([]string{player}, fmt.Sprintf(
		"%s\r\n%s",
		strings.Join(headers, "\r\n"),
		buf.String(),
	)); err != nil {
		return "", fmt.Errorf("Unable to send dump email: %v", err)
	}
	return "You have been emailed a game dump", nil
}

func (c DumpCommand) Usage(player string, context interface{}) string {
	return "{{b}}dump{{_b}} to dump the game data which can be used for troubleshooting"
}
