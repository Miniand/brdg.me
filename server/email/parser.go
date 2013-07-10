package main

import (
	"bytes"
	"errors"
	"github.com/beefsack/boredga.me/game"
	"github.com/beefsack/boredga.me/render"
	"labix.org/v2/mgo/bson"
	"log"
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
func ParseBody(body string) [][]string {
	commandSplitReg := regexp.MustCompile("\\s+")
	// Convert all CRLF and CR to LF
	cleanedBody := strings.Replace(strings.Replace(body, "\r\n", "\n", -1),
		"\r", "\n", -1)
	// Break block down to lines of commands
	commands := strings.Split(cleanedBody, "\n")
	parsedCommands := [][]string{}
	// Break each command down to parts with spaces
	contentStarted := false
	for _, command := range commands {
		trimmedCommand := strings.TrimSpace(command)
		if trimmedCommand == "" {
			if contentStarted {
				// Don't allow any blank lines after initial content
				break
			}
		} else {
			contentStarted = true
			parsedCommands = append(parsedCommands, commandSplitReg.Split(
				trimmedCommand, -1))
		}
	}
	return parsedCommands
}

// Run commands on the game, email relevant people and handle action issues
func HandleCommands(player, gameId string, commands [][]string) error {
	if gameId == "" {
		// Either starting a new game or just print help
		if len(commands) > 0 && len(commands[0]) > 2 &&
			strings.ToLower(commands[0][0]) == "new" {
			gType := game.Collection()[commands[0][1]]
			if gType != nil {
				// We found the game, lets try to start it
				players := commands[0][2:]
				players = append(players, player)
				g, err := gType(players)
				if err != nil {
					// The game couldn't start
					err := SendMail([]string{player}, "Couldn't start game: "+
						err.Error())
					if err != nil {
						return err
					}
				} else {
					// We started the game, lets kick it off
					gm, err := SaveGame(g)
					if err != nil {
						return err
					}
					err = CommunicateGameTo(gm.Id, g, g.PlayerList(),
						"You have been invited by "+player+
							" to play "+g.Name()+" by email!")
					if err != nil {
						return err
					}
				}
			} else {
				err := SendMail([]string{player}, "Invalid command.\n\n"+
					GeneralHelpText())
				if err != nil {
					return err
				}
			}
		} else {
			// Print help
			err := SendMail([]string{player}, GeneralHelpText())
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
		initialWhoseTurn := g.WhoseTurn()
		if err != nil {
			return err
		}
		commandRun := false
		for _, command := range commands {
			err = g.PlayerAction(player, command[0], command[1:])
			if err != nil {
				// Don't try to do any more commands
				break
			}
			commandRun = true
		}
		commErrs := []string{}
		header := ""
		if err != nil {
			header = err.Error()
		}
		commErr := CommunicateGameTo(gm.Id, g, []string{player}, header)
		if commErr != nil {
			commErrs = append(commErrs, commErr.Error())
		}
		if commandRun {
			_, err := UpdateGame(bson.ObjectIdHex(gameId), g)
			if err != nil {
				return err
			}
			// Email any players who now have a turn
			commErr = CommunicateGameTo(gm.Id, g,
				WhoseTurnNow(g, initialWhoseTurn), "")
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
		}
		if len(commErrs) > 0 {
			return errors.New(strings.Join(commErrs, "\n"))
		}
	}
	return nil
}

func GeneralHelpText() string {
	gameNames := []string{}
	for _, g := range game.RawCollection() {
		gameNames = append(gameNames, g.Identifier()+"   ("+g.Name()+
			")")
	}
	return `Welcome to boredga.me!

Please start a new game by using the "new" command like below, but using the game name and player emails you want.

new tic_tac_toe player1@example.com player2@example.com

The available games are:

` + strings.Join(gameNames, "\n")
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
	header string) error {
	var footer string
	if g.IsFinished() {
		winners := g.Winners()
		footer = "The game is over"
		if len(winners) == 0 {
			footer += ", it was a draw!"
		} else {
			footer += ", the winners were: " + strings.Join(winners, ", ")
		}
	} else {
		footer = "Current turn: " + strings.Join(g.WhoseTurn(), ", ")
	}
	for _, p := range to {
		if header != "" {
			header += "\n\n"
		}
		rawOutput, err := g.RenderForPlayer(p)
		if err != nil {
			return err
		}
		raw := header + rawOutput + "\n\n" + footer
		terminalOutput, err := render.RenderTerminal(raw, g)
		if err != nil {
			return err
		}
		htmlOutput, err := render.RenderHtml(raw, g)
		if err != nil {
			return err
		}
		// Make a multipart message
		buf := &bytes.Buffer{}
		data := multipart.NewWriter(buf)
		htmlW, err := data.CreatePart(textproto.MIMEHeader{
			"Content-Type": []string{"text/html"},
		})
		if err != nil {
			return err
		}
		_, err = htmlW.Write([]byte("<pre>" + htmlOutput))
		if err != nil {
			return err
		}
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
		err = data.Close()
		if err != nil {
			return err
		}
		err = SendMail([]string{p},
			"Subject: "+g.Name()+" ("+id.(bson.ObjectId).Hex()+")\n"+
				"MIME-Version: 1.0\n"+
				"Content-Type: multipart/alternative; boundary="+data.Boundary()+"\n"+
				buf.String())
		log.Println("Subject: " + g.Name() + " (" + id.(bson.ObjectId).Hex() + ")\n" +
			"MIME-Version: 1.0\n" +
			"Content-Type: multipart/alternative; boundary=" + data.Boundary() + "\n" +
			buf.String())
		if err != nil {
			return err
		}
	}
	return nil
}
