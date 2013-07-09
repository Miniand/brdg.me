package main

import (
	"errors"
	"github.com/beefsack/boredga.me/game"
	"labix.org/v2/mgo/bson"
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
		commErr := CommunicateGameTo(g, []string{player}, header)
		if commErr != nil {
			commErrs = append(commErrs, commErr.Error())
		}
		if commandRun {
			// Email any players who now have a turn
			commErr = CommunicateGameTo(g, WhoseTurnNow(g, initialWhoseTurn),
				"")
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

func CommunicateGameTo(g game.Playable, to []string, header string) error {
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
		output, err := g.RenderForPlayer(p)
		if err != nil {
			return err
		}
		err = SendMail([]string{p}, header+output+"\n\n"+footer)
		if err != nil {
			return err
		}
	}
	return nil
}
