package love_letter

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type CharPriest struct{}

func (p CharPriest) Name() string { return "Priest" }
func (p CharPriest) Number() int  { return Priest }
func (p CharPriest) Text() string {
	return "Look at another player's hand"
}

func (p CharPriest) Play(g *Game, player int, args ...string) error {
	target, err := g.ParseTarget(player, args...)
	if err != nil {
		return err
	}

	g.DiscardCard(player, Priest)

	if target == player {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s played %s, but had nobody to target so just discarded the card",
			g.RenderName(player),
			RenderCard(Priest),
		)))
		return nil
	}

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s and looked at %s's hand",
		g.RenderName(player),
		RenderCard(Priest),
		g.RenderName(target),
	)))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"%s has %s",
		g.RenderName(target),
		render.CommaList(RenderCards(g.Hands[target])),
	), []string{g.Players[player]}))

	return nil
}
