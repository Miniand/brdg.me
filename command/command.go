package command

import (
	"errors"
	"strings"
)

// Command is a regexp based command parser, providing an interface for
// authorisation and instructions.
type Command interface {
	// Parses the input string for the command, the return is nil if it could not parse the command, or if it could
	Parse(input string) []string
	CanCall(player string, context interface{}) bool
	Call(player string, context interface{}, args []string) error
	Usage(player string, context interface{}) string
}

// Tried to call a command given a range of command parsers.  Errors if it was
// unable to match to any commands.
func CallInCommands(player string, context interface{}, input string,
	commands []Command) (err error) {
	numRun := 0
	for {
		input = strings.TrimSpace(input)
		initialInput := input
		for _, c := range commands {
			if c.CanCall(player, context) {
				args := c.Parse(input)
				if args != nil {
					// The command matches
					numRun++
					// Trim the matched text out of the input string
					input = input[len(args[0]):]
					err = c.Call(player, context, args)
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
