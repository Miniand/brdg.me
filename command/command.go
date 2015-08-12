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
	// Name is the name of the command and is used to decide which to call.
	Name() string
	Call(player string, context interface{}, input *Parser) (output string,
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
	p := NewParserString(input)
	for {
		commandOutput, err = CallOneInCommands(
			player,
			context,
			p,
			commands,
		)
		p.ReadString('\n') // Flush line to prepare for more commands
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
	p *Parser,
	commands []Command,
) (output string, err error) {
	p.ReadSpace() // Clear leading space
	cmd, err := p.ReadWord()
	if err != nil {
		// Unable to read command.
		err = ErrNoCommandFound
		return
	}
	lowerCmd := strings.ToLower(cmd)
	foundCommand := false
	for _, c := range commands {
		if strings.ToLower(c.Name()) != lowerCmd {
			continue
		}
		foundCommand = true
		output, err = c.Call(player, context, p)
		break
	}
	if !foundCommand {
		// No commands ran or there was an error, stop running
		err = ErrNoCommandFound
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
