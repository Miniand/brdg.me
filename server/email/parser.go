package main

import (
	"bytes"
	"errors"
	"github.com/beefsack/brdg.me/command"
	"github.com/beefsack/brdg.me/game"
	"github.com/beefsack/brdg.me/render"
	"labix.org/v2/mgo/bson"
	"mime/multipart"
	"net/textproto"
	"regexp"
	"strings"
)

// Search for an email address
func ParseFrom(from string) string {
	reg := regexp.MustCompile(EmailSearchRegexString())
	return reg.FindString(from)
}

// Search for a BSON objectid to match to a game (length 24 hex string)
func ParseSubject(subject string) string {
	reg := regexp.MustCompile("\\b[a-f0-9]{24}\\b")
	return reg.FindString(subject)
}

// Find contiguous lines as commands until the first blank line
func ParseBody(body string) string {
	return strings.Replace(strings.Replace(body, "\r\n", "\n", -1),
		"\r", "\n", -1)
}

// Run commands on the game, email relevant people and handle action issues
func HandleCommandText(player, gameId string, commandText string) error {
	if gameId == "" {
		commands := Commands()
		err := command.CallInCommands(player, nil, commandText, commands)
		if err != nil {
			// Print help
			body := bytes.NewBufferString("")
			if err != command.NO_COMMAND_FOUND {
				body.WriteString(err.Error())
				body.WriteString("\n\n")
			}
			body.WriteString("Welcome to brdg.me!\n\n")
			htmlOutput, err := render.RenderHtml(render.OutputCommands(player,
				nil, command.AvailableCommands(player, nil, commands)))
			if err != nil {
				return err
			}
			body.WriteString(htmlOutput)
			err = SendMail([]string{player}, body.String())
			if err != nil {
				return err
			}
		}
	} else {
		gm, err := LoadGame(bson.ObjectIdHex(gameId))
		if err != nil {
			return err
		}
		g, err := gm.ToGame()
		if err != nil {
			return err
		}
		commands := append(g.Commands(), Commands()...)
		initialWhoseTurn := g.WhoseTurn()
		err = command.CallInCommands(player, g, commandText, commands)
		commErrs := []string{}
		header := ""
		if err != nil {
			header = err.Error()
		}
		commErr := CommunicateGameTo(gm.Id, g, []string{player}, header, false)
		if commErr != nil {
			commErrs = append(commErrs, commErr.Error())
		}
		if err != command.NO_COMMAND_FOUND {
			_, err := UpdateGame(bson.ObjectIdHex(gameId), g)
			if err != nil {
				return err
			}
			// Email any players who now have a turn
			commErr = CommunicateGameTo(gm.Id, g,
				WhoseTurnNow(g, initialWhoseTurn), "", false)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			// Update again to handle saves during render, ie for logger
			_, err = UpdateGame(bson.ObjectIdHex(gameId), g)
			if err != nil {
				return err
			}
		}
		if len(commErrs) > 0 {
			return errors.New(strings.Join(commErrs, "\n"))
		}
	}
	return nil
}

func WhoseTurnNow(g game.Playable, initialWhoseTurn []string) []string {
	initialWhoseTurnMap := map[string]bool{}
	for _, p := range initialWhoseTurn {
		initialWhoseTurnMap[p] = true
	}
	whoseTurnNow := []string{}
	for _, p := range g.WhoseTurn() {
		if !initialWhoseTurnMap[p] {
			whoseTurnNow = append(whoseTurnNow, p)
		}
	}
	return whoseTurnNow
}

func CommunicateGameTo(id interface{}, g game.Playable, to []string,
	header string, initial bool) error {
	if header != "" {
		header += "\n\n"
	}
	if g.IsFinished() {
		winners := g.Winners()
		header += "The game is over"
		if len(winners) == 0 {
			header += ", it was a draw!"
		} else {
			header += ", the winners were: " + strings.Join(winners, ", ")
		}
	} else {
		header = "Current turn: " + strings.Join(g.WhoseTurn(), ", ")
	}
	for _, p := range to {
		pHeader := header
		rawOutput, err := g.RenderForPlayer(p)
		commands := append(g.Commands(), Commands()...)
		available := command.AvailableCommands(p, g, commands)
		if len(available) > 0 {
			pHeader += "\n\n" + render.OutputCommands(p, g, available)
		}
		if err != nil {
			return err
		}
		raw := pHeader + "\n\n" + rawOutput
		terminalOutput, err := render.RenderTerminal(raw)
		if err != nil {
			return err
		}
		htmlOutput, err := render.RenderHtml(raw)
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
		// Handle Message-ID and In-Reply-To for email threading
		threadingHeaders := ""
		messageId := id.(bson.ObjectId).Hex() + "@brdg.me"
		if initial {
			// We create the base Message-ID
			threadingHeaders = "Message-Id: <" + messageId + ">\r\n"
		} else {
			// We create a unique Message-ID and set the In-Reply-To to original
			threadingHeaders = "In-Reply-To: <" + messageId + ">\r\n"
		}
		err = SendMail([]string{p},
			"Subject: "+g.Name()+" ("+id.(bson.ObjectId).Hex()+")\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: multipart/alternative; boundary="+data.Boundary()+"\r\n"+
				threadingHeaders+
				buf.String())
		if err != nil {
			return err
		}
	}
	return nil
}
