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
	Free   bool
	Wonder bool
	Deal   int
	Chosen bool
}

func (a *BuildAction) IsComplete() bool {
	return a.Chosen
}

func (a *BuildAction) DealOptions(player int, g *Game) []map[int]int {
	if a.Free {
		return []map[int]int{}
	}
	_, coins := g.CanBuildCard(player, a.GetCard(player, g))
	return coins
}

func (a *BuildAction) ChooseDeal(player int, g *Game, n int) error {
	_, coins := g.CanBuildCard(player, a.GetCard(player, g))
	if n < 0 || n >= len(coins) {
		return errors.New("that is not a valid deal number")
	}
	a.Deal = n
	a.Chosen = true
	g.CheckHandComplete()
	return nil
}

func (a *BuildAction) HandlePostActionExecute(player int, g *Game) {
	// The new card will be the last in the player's hand
	c := g.Cards[player][len(g.Cards[player])-1]
	if post, ok := c.(PostActionExecuteHandler); ok {
		post.HandlePostActionExecute(player, g)
	}
}

func (a *BuildAction) GetCard(player int, g *Game) Carder {
	if a.Wonder {
		return g.RemainingWonderStages(player)[0].(Carder)
	}
	return g.Hands[player][a.Card].(Carder)
}

func (a *BuildAction) Execute(player int, g *Game) {
	c := a.GetCard(player, g)
	crd := c.GetCard()
	paymentString := ""
	if a.Free {
		paymentString = " for free"
	} else {
		g.Coins[player] -= crd.Cost[GoodCoin]

		_, coins := g.CanBuildCard(player, c)
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

	if a.Free {
		for _, c := range g.PlayerThings(player) {
			if free, ok := c.(FreeBuilder); ok {
				free.HandleFreeBuild()
			}
		}
	}
}

func (a *BuildAction) Output(player int, g *Game) string {
	c := a.GetCard(player, g)
	buf := bytes.NewBufferString("building ")
	buf.WriteString(RenderCard(c))
	if !a.Free {
		_, coins := g.CanBuildCard(player, c)
		if len(coins) > 0 {
			buf.WriteString("\n")
			buf.WriteString(g.RenderDeals(player, coins))
		}
	}
	return buf.String()
}
