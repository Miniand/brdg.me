package starship_catan

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type CompleteCommand struct{}

func (c CompleteCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("complete", 1, input)
}

func (c CompleteCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanComplete(p)
}

func (c CompleteCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}

	a := command.ExtractNamedCommandArgs(args)
	adventure, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("you must specify an adventure number")
	}

	return "", g.Complete(p, adventure)
}

func (c CompleteCommand) Usage(player string, context interface{}) string {
	return "{{b}}complete #{{_b}} to complete an adventure, eg. {{b}}complete 2{{_b}}"
}
