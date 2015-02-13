package seven_wonders

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type BuildAction struct {
	Card   int
	Deal   int
	Chosen bool
}

func (a *BuildAction) IsComplete() bool {
	return a.Chosen
}

func (a *BuildAction) DealOptions(player int, g *Game) []map[int]int {
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

func (a *BuildAction) Execute(player int, g *Game) {
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

func (a *BuildAction) Output(player int, g *Game) string {
	c := g.Hands[player][a.Card].(Carder)
	buf := bytes.NewBufferString("{{b}}Building:{{_b}} ")
	buf.WriteString(RenderCard(c))
	buf.WriteString("\n")
	_, coins := g.CanBuildCard(player, c)
	buf.WriteString(g.RenderDeals(player, coins))
	return buf.String()
}
