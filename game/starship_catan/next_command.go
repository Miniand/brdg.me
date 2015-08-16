package starship_catan

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type NextCommand struct{}

func (c NextCommand) Name() string { return "next" }

func (c NextCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Next(p)
}

func (c NextCommand) Usage(player string, context interface{}) string {
	return "{{b}}next{{_b}} to advance to the next card"
}

func (g *Game) CanNext(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight ||
		g.FlightCards.Len() == 0 || g.GainResources != nil {
		return false
	}
	if g.CardFinished {
		return true
	}
	c, _ := g.FlightCards.Pop()
	actioner, ok := c.(Actioner)
	return !ok || !actioner.RequiresAction()
}

func (g *Game) Next(player int) error {
	if !g.CanNext(player) {
		return errors.New("you can't advance to the next card")
	}
	return g.NextSectorCard()
}
