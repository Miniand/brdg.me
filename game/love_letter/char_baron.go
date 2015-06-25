package love_letter

import (
	"errors"
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
	targets := g.AvailableTargets(player)
	if len(targets) > 0 && len(args) != 1 {
		return errors.New("please specify the other player to compare hands with")
	} else if len(targets) == 0 {
		if len(args) == 0 {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"%s played %s, but had nobody to target so just discarded the card",
				g.RenderName(player),
				RenderCard(Baron),
			)))
			return nil
		}
		return errors.New("because there are no possible targets for your Baron, please play it without arguments to discard it without activating it")
	}

	target, err := helper.MatchStringInStrings(args[0], g.Players)
	if err != nil {
		return err
	}
	if _, ok := helper.IntFind(target, targets); !ok {
		return ErrCannotTarget
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
	} else if diff > 0 {
		eliminate = target
	}

	if eliminate == -1 {
		g.Log.Add(log.NewPublicMessage("The cards were equal, nobody is eliminated"))
	} else {
		g.Eliminate(eliminate)
	}
	return nil
}
