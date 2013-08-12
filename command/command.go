package command

import (
	"errors"
	"github.com/beefsack/brdg.me/game"
	"strings"
)

type Command interface {
	// Parses the input string for the command, the return is nil if it could not parse the command, or if it could
	Parse(input string) []string
	CanCall(player string, g *game.Playable) bool
	Call(player string, g *game.Playable, args []string) error
	Usage(player string, g *game.Playable) string
}

func CallInCommands(player string, g *game.Playable, input string,
	commands []Command) (err error) {
	numRun := 0
	for {
		input = strings.TrimSpace(input)
		initialInput := input
		for _, c := range commands {
			if c.CanCall(player, g) {
				args := c.Parse(input)
				if args != nil {
					// The command matches
					numRun++
					// Trim the matched text out of the input string
					input = input[len(args[0]):]
					err = c.Call(player, g, args)
					break
				}
			}
		}
		if err != nil || input == initialInput {
			// No commands ran or there was an error, stop running
			break
		}
	}
	if numRun == 0 {
		err = errors.New("We couldn't find any commands in the text you sent, please make sure your commands are at the top")
	}
	return
}
