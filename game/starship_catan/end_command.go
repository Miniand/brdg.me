package starship_catan

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type EndCommand struct{}

func (c EndCommand) Name() string { return "end" }

func (c EndCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
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
