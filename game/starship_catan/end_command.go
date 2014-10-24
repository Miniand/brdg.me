package starship_catan

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type EndCommand struct{}

func (c EndCommand) Parse(input string) []string {
	return command.ParseNamedCommand("end", input)
}

func (c EndCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanEnd(p)
}

func (c EndCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.End(p)
}

func (c EndCommand) Usage(player string, context interface{}) string {
	return "{{b}}end{{_b}} to end the flight early"
}

func (g *Game) CanEnd(player int) bool {
	return g.CanNext(player)
}

func (g *Game) End(player int) error {
	if !g.CanEnd(player) {
		return errors.New("cannot end the flight at the moment")
	}
	return g.EndFlight()
}
