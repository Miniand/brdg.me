package love_letter

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type CharPrince struct{}

func (p CharPrince) Name() string { return "Prince" }
func (p CharPrince) Number() int  { return Prince }
func (p CharPrince) Text() string {
	return "Choose a player to discard and draw a new card"
}

func (p CharPrince) Play(g *Game, player int, args ...string) error {
	if _, ok := helper.IntFind(Countess, g.Hands[player]); ok {
		return errors.New("you must play the Countess")
	}

	target, err := g.ParseTarget(player, args...)
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

	g.DiscardCardLog(target, g.Hands[target][0])
	g.DrawCard(target)

	return nil
}
