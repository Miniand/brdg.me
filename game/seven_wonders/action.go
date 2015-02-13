package seven_wonders

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

func init() {
	gob.Register(BuildAction{})
	gob.Register(WonderAction{})
	gob.Register(DiscardAction{})
}

type Actioner interface {
	IsComplete() bool
	Execute(player int, g *Game)
}

type DealOptioner interface {
	DealOptions() []map[int]int
	ChooseDeal(n int) error
}

type BuildAction struct {
	Card   int
	Deal   int
	Chosen bool
}

func (a BuildAction) IsComplete() bool {
	return a.Chosen
}

func (a BuildAction) DealOptions(player int, g *Game) []map[int]int {
	_, coins := g.CanBuildCard(player, g.Hands[player][a.Card].(Carder))
	return coins
}

func (a *BuildAction) ChooseDeal(player int, g *Game, n int) error {
	_, coins := g.CanBuildCard(player, g.Hands[player][a.Card].(Carder))
	if n < 0 || n >= len(coins) {
		return errors.New("that is not a valid deal number")
	}
	a.Deal = n
	a.Chosen = true
	g.CheckHandComplete()
	return nil
}

func (a BuildAction) Execute(player int, g *Game) {
	c := g.Hands[player][a.Card].(Carder)
	crd := c.GetCard()
	g.Coins[player] -= crd.Cost[GoodCoin]

	_, coins := g.CanBuildCard(player, c)
	paymentString := ""
	if len(coins) > 0 {
		parts := []string{}
		for dir, amt := range coins[a.Deal] {
			targetPlayer := g.NumFromPlayer(player, dir)
			g.Coins[player] -= amt
			g.Coins[targetPlayer] += amt
			parts = append(parts, fmt.Sprintf(
				"%s %s",
				g.PlayerName(targetPlayer),
				RenderMoney(amt),
			))
		}
		if len(parts) > 0 {
			paymentString = fmt.Sprintf(
				", paying %s",
				render.CommaList(parts),
			)
		}
	}

	g.Cards[player] = g.Cards[player].Push(c)
	g.Hands[player] = append(
		g.Hands[player][:a.Card],
		g.Hands[player][a.Card+1:]...,
	)

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s built %s%s",
		g.PlayerName(player),
		RenderCard(c),
		paymentString,
	)))
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
