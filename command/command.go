package command

import (
	"errors"
	"strings"
)

var NO_COMMAND_FOUND = errors.New(
	`Invalid command, please make sure your commands are at the top, and please check the "You can:" section to see a list of available commands`)

// Command is a regexp based command parser, providing an interface for
// authorisation and instructions.
type Command interface {
	// Parses the input string for the command, the return is nil if it could not parse the command, or if it could
	Parse(input string) []string
	CanCall(player string, context interface{}) bool
	Call(player string, context interface{}, args []string) (output string,
		err error)
	Usage(player string, context interface{}) string
}

// Tried to call a command given a range of command parsers.  Errors if it was
// unable to match to any commands.
func CallInCommands(player string, context interface{}, input string,
	commands []Command) (output string, err error) {
	var commandOutput string
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
					commandOutput, err = c.Call(player, context, args)
					// If we got some output, add it
					if commandOutput != "" {
						if output != "" {
							output += "\n"
						}
						output += commandOutput
					}
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
		err = NO_COMMAND_FOUND
	}
	return
}

func AvailableCommands(player string, context interface{},
	commands []Command) (available []Command) {
	for _, c := range commands {
		if c.CanCall(player, context) {
			available = append(available, c)
		}
	}
	return
}

func CommandUsages(player string, context interface{},
	commands []Command) (usages []string) {
	for _, c := range commands {
		usage := c.Usage(player, context)
		if usage != "" {
			usages = append(usages, usage)
		}
	}
	return
}
