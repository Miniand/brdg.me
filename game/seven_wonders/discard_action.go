package seven_wonders

import (
	"bytes"
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type DiscardAction struct {
	Card int
}

func (a *DiscardAction) IsComplete() bool {
	return true
}

func (a *DiscardAction) Execute(player int, g *Game) {
	g.Discard = g.Discard.Push(g.Hands[player][a.Card])
	g.Hands[player] = append(
		g.Hands[player][:a.Card],
		g.Hands[player][a.Card+1:]...,
	)
	g.Coins[player] += 3
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded a card for %s",
		g.PlayerName(player),
		RenderMoney(3),
	)))
}

func (a *DiscardAction) Output(player int, g *Game) string {
	c := g.Hands[player][a.Card].(Carder)
	buf := bytes.NewBufferString("discarding ")
	buf.WriteString(RenderCard(c))
	return buf.String()
}
