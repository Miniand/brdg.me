package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type FoundColonyCommand struct{}

func (c FoundColonyCommand) Name() string { return "found" }

func (c FoundColonyCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.FoundColony(p)
}

func (c FoundColonyCommand) Usage(player string, context interface{}) string {
	return "{{b}}found{{_b}} to found a colony here"
}

func (g *Game) CanFoundColony(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight ||
		len(g.FlightCards) == 0 ||
		g.PlayerBoards[player].Resources[ResourceColonyShip] == 0 {
		return false
	}
	c, _ := g.FlightCards.Pop()
	_, ok := c.(ColonyCard)
	return ok
}

func (g *Game) FoundColony(player int) error {
	var c card.Card

	if !g.CanFoundColony(player) {
		return errors.New("you are not able to found a colony")
	}
	c, g.FlightCards = g.FlightCards.Pop()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s founded a colony on %s`, g.RenderName(player), c)))
	g.PlayerBoards[player].Colonies = g.PlayerBoards[player].Colonies.Push(c)
	g.PlayerBoards[player].Resources[ResourceColonyShip] -= 1
	g.ReplaceCard()
	g.MarkCardActioned()
	g.RecalculatePeopleCards()
	return nil
}
