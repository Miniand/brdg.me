package love_letter

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type CharBaron struct{}

func (p CharBaron) Name() string { return "Baron" }
func (p CharBaron) Number() int  { return Baron }
func (p CharBaron) Text() string {
	return "Compare hands with another player, lowest card is eliminated"
}

func (p CharBaron) Play(g *Game, player int, args ...string) error {
	target, err := g.ParseTarget(player, args...)
	if err != nil {
		return err
	}

	g.DiscardCard(player, Baron)

	if target == player {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s played %s, but had nobody to target so just discarded the card",
			g.RenderName(player),
			RenderCard(Baron),
		)))
		return nil
	}

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s and is comparing hands with %s to see who has a lower card",
		g.RenderName(player),
		RenderCard(Baron),
		g.RenderName(target),
	)))
	playerCard := helper.IntRemove(Baron, g.Hands[player], 1)[0]
	targetCard := g.Hands[target][0]
	logFmt := "You have %s, %s has %s"
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		logFmt,
		RenderCard(playerCard),
		g.RenderName(target),
		RenderCard(targetCard),
	), []string{g.Players[player]}))
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		logFmt,
		RenderCard(targetCard),
		g.RenderName(player),
		RenderCard(playerCard),
	), []string{g.Players[target]}))

	eliminate := -1
	diff := Cards[playerCard].Number() - Cards[targetCard].Number()
	if diff < 0 {
		eliminate = player
		g.Hands[player] = []int{playerCard}
	} else if diff > 0 {
		eliminate = target
		g.Hands[target] = []int{targetCard}
	}

	if eliminate == -1 {
		g.Log.Add(log.NewPublicMessage("The cards were equal, nobody is eliminated"))
	} else {
		g.Eliminate(eliminate)
	}

	return nil
}
