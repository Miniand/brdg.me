package main

import (
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
		if err != nil {
			return err
		}
		for _, command := range commands {
			err = g.PlayerAction(player, command[0], command[1:])
			if err != nil {
				// Stop here and email
			}
		}
	}
	return nil
}
