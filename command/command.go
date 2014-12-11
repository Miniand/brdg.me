package command

import (
	"errors"
	"strings"
)

var ErrNoCommandFound = errors.New(
	`Invalid command, please check that your command matches one of the available commands`)

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
func CallInCommands(
	player string,
	context interface{},
	input string,
	commands []Command,
) (output string, err error) {
	return CallInCommandsPostHook(
		player,
		context,
		input,
		commands,
		nil,
	)
}

func CallInCommandsPostHook(
	player string,
	context interface{},
	input string,
	commands []Command,
	postCommandHook func() error,
) (output string, err error) {
	var commandOutput string
	numRun := 0
	outputs := []string{}
	for {
		input, commandOutput, err = CallOneInCommands(
			player,
			context,
			input,
			commands,
		)
		if commandOutput != "" {
			outputs = append(outputs, commandOutput)
		}
		if postCommandHook != nil && err == nil {
			err = postCommandHook()
		}
		if err == ErrNoCommandFound {
			if numRun > 0 {
				// One of the commands ran successfully
				err = nil
			}
			break
		} else if err != nil {
			break
		}
		numRun += 1
	}
	output = strings.Join(outputs, "\n")
	return
}

func CallOneInCommands(
	player string,
	context interface{},
	input string,
	commands []Command,
) (remaining, output string, err error) {
	remaining = strings.TrimSpace(input)
	initialInput := remaining
	for _, c := range commands {
		if c.CanCall(player, context) {
			args := c.Parse(remaining)
			if args != nil {
				// Trim the matched text out of the remaining string
				remaining = remaining[len(args[0]):]
				output, err = c.Call(player, context, args)
				break
			}
		}
	}
	if err == nil && remaining == initialInput {
		// No commands ran or there was an error, stop running
		err = ErrNoCommandFound
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
