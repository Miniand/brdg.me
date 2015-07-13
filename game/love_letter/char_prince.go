package love_letter

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type CharPrince struct{}

func (p CharPrince) Name() string { return "Prince" }
func (p CharPrince) Number() int  { return Prince }
func (p CharPrince) Text() string {
	return "Choose a player (or yourself) to discard and draw a new card"
}
func (p CharPrince) Colour() string { return render.Magenta }

func (p CharPrince) Play(g *Game, player int, args ...string) error {
	if _, ok := helper.IntFind(Countess, g.Hands[player]); ok {
		return errors.New("you must play the Countess")
	}

	target, err := g.ParseTarget(player, true, args...)
	if err != nil {
		return err
	}

	g.DiscardCard(player, Prince)

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s and made %s discard their hand and draw a new card",
		g.RenderName(player),
		RenderCard(Prince),
		g.RenderName(target),
	)))

	curRound := g.Round
	g.DiscardCardLog(target, g.Hands[target][0])
	if g.Round == curRound && !g.Eliminated[target] {
		g.DrawCard(target)
	}

	return nil
}
