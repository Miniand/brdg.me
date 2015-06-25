package love_letter

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type CharKing struct{}

func (p CharKing) Name() string { return "King" }
func (p CharKing) Number() int  { return King }
func (p CharKing) Text() string {
	return "Trade your hand with another player"
}
func (p CharKing) Colour() string { return render.Blue }

func (p CharKing) Play(g *Game, player int, args ...string) error {
	if _, ok := helper.IntFind(Countess, g.Hands[player]); ok {
		return errors.New("you must play the Countess")
	}

	target, err := g.ParseTarget(player, false, args...)
	if err != nil {
		return err
	}

	g.DiscardCard(player, King)

	if target == player {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s played %s, but had nobody to target so just discarded the card",
			g.RenderName(player),
			RenderCard(King),
		)))
		return nil
	}

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s and swapped hands with %s",
		g.RenderName(player),
		RenderCard(King),
		g.RenderName(target),
	)))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"You traded your %s for %s",
		RenderCard(g.Hands[player][0]),
		RenderCard(g.Hands[target][0]),
	), []string{g.Players[player]}))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"You traded your %s for %s",
		RenderCard(g.Hands[target][0]),
		RenderCard(g.Hands[player][0]),
	), []string{g.Players[target]}))

	return nil
}
