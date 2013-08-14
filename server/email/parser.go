package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/beefsack/brdg.me/command"
	"github.com/beefsack/brdg.me/game"
	"github.com/beefsack/brdg.me/render"
	"github.com/beefsack/brdg.me/server/model"
	"labix.org/v2/mgo/bson"
	"regexp"
	"strings"
)

var currentGameId interface{}

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
	unsubscribed, err := UserIsUnsubscribed(player)
	if (err == nil && unsubscribed) || gameId == "" {
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
			body.WriteString(render.OutputCommands(player, nil,
				command.AvailableCommands(player, nil, commands)))
			err = SendRichMail([]string{player}, "Welcome to brdg.me!",
				body.String(), []string{})
			if err != nil {
				return err
			}
		}
	} else {
		var initialEliminated []string
		gm, err := model.LoadGame(bson.ObjectIdHex(gameId))
		if err != nil {
			return err
		}
		g, err := gm.ToGame()
		if err != nil {
			return err
		}
		currentGameId = gm.Id
		commands := append(g.Commands(), Commands()...)
		initialWhoseTurn := g.WhoseTurn()
		eliminator, isEliminator := g.(game.Eliminator)
		if isEliminator {
			initialEliminated = eliminator.EliminatedPlayerList()
		}
		err = command.CallInCommands(player, g, commandText, commands)
		header := ""
		if err != nil {
			header = err.Error()
		}
		commErrs := []string{}
		commErr := CommunicateGameTo(gm.Id, g, []string{player}, header, false)
		if commErr != nil {
			commErrs = append(commErrs, commErr.Error())
		}
		if err != command.NO_COMMAND_FOUND {
			_, err := model.UpdateGame(bson.ObjectIdHex(gameId), g)
			if err != nil {
				return err
			}
			// Email any players who now have a turn
			commErr = CommunicateGameTo(gm.Id, g,
				WhoseTurnNow(g, initialWhoseTurn), "", false)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			// Email any players who were eliminated this turn
			if isEliminator {
				commErr = CommunicateGameTo(gm.Id, g, FindNewStringsInSlice(
					initialEliminated, eliminator.EliminatedPlayerList()),
					"You have been eliminated from the game.", false)
				if commErr != nil {
					commErrs = append(commErrs, commErr.Error())
				}
			}
			// Update again to handle saves during render, ie for logger
			_, err = model.UpdateGame(bson.ObjectIdHex(gameId), g)
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
	return FindNewStringsInSlice(initialWhoseTurn, g.WhoseTurn())
}

func FindNewStringsInSlice(oldSlice, newSlice []string) (newStrings []string) {
	oldSliceMap := map[string]bool{}
	for _, s := range oldSlice {
		oldSliceMap[s] = true
	}
	for _, s := range newSlice {
		if !oldSliceMap[s] {
			newStrings = append(newStrings, s)
		}
	}
	return
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
		header += "Current turn: " + strings.Join(g.WhoseTurn(), ", ")
	}
	for _, p := range to {
		unsubscribed, err := UserIsUnsubscribed(p)
		if err == nil && unsubscribed {
			continue
		}
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
		body := pHeader + "\n\n" + rawOutput
		subject := fmt.Sprintf("%s (%s)", g.Name(), id.(bson.ObjectId).Hex())
		extraHeaders := []string{}
		messageId := id.(bson.ObjectId).Hex() + "@brdg.me"
		if initial {
			// We create the base Message-ID
			extraHeaders = append(extraHeaders,
				fmt.Sprintf("Message-Id: <%s>", messageId))
		} else {
			// We create a unique Message-ID and set the In-Reply-To to original
			extraHeaders = append(extraHeaders,
				fmt.Sprintf("In-Reply-To: <%s>", messageId))
		}
		err = SendRichMail([]string{p}, subject, body, extraHeaders)
		if err != nil {
			return err
		}
	}
	return nil
}
