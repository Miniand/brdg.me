package seven_wonders

import (
	"encoding/gob"
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

func init() {
	gob.Register(&BuildAction{})
	gob.Register(&WonderAction{})
	gob.Register(&DiscardAction{})
}

type Actioner interface {
	IsComplete() bool
	Execute(player int, g *Game)
	Output(player int, g *Game) string
}

type DealOptioner interface {
	DealOptions(player int, g *Game) []map[int]int
	ChooseDeal(player int, g *Game, n int) error
}

type WonderAction struct {
	Card   int
	Deal   int
	Chosen bool
}

func (a WonderAction) IsComplete() bool {
	return a.Chosen
}

func (a WonderAction) DealOptions(player int, g *Game) []map[int]int {
	return nil
}

type DiscardAction struct {
	Card int
}

func (a DiscardAction) IsComplete() bool {
	return true
}

func (a DiscardAction) Execute(player int, g *Game) {
	g.Cards[player] = append(
		g.Cards[player][:a.Card],
		g.Cards[player][a.Card+1:]...,
	)
	g.Coins[player] += 3
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded a card for %s",
		g.PlayerName(player),
		RenderMoney(3),
	)))
}
