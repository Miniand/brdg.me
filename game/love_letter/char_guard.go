package love_letter

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type CharGuard struct{}

func (p CharGuard) Name() string { return "Guard" }
func (p CharGuard) Number() int  { return Guard }
func (p CharGuard) Text() string {
	return "Guess another player's card to eliminate them, except for Guard"
}

func (p CharGuard) Play(g *Game, player int, args ...string) error {
	if len(args) != 2 {
		return errors.New("please specify which player to target and what card you think they are, eg. play guard steve handmaid")
	}

	target, err := g.ParseTarget(player, args...)
	if err != nil {
		return err
	}
	if target == player {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s played %s, but had nobody to target so just discarded the card",
			g.RenderName(player),
			RenderCard(Priest),
		)))
		return nil
	}

	names := map[int]string{}
	for i, c := range Cards {
		names[i] = c.Name()
	}
	card, err := helper.MatchStringInStringMap(args[1], names)
	if err != nil {
		return err
	}
	if card == Guard {
		return errors.New("you can't use Guard against other Guards")
	}

	prefix := fmt.Sprintf(
		"%s played %s and guessed that %s is a %s, ",
		g.RenderName(player),
		RenderCard(Guard),
		g.RenderName(target),
		RenderCard(card),
	)

	if _, ok := helper.IntFind(card, g.Hands[target]); ok {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%sand was correct!",
			prefix,
		)))
		g.Eliminate(target)
	} else {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%sbut was incorrect",
			prefix,
		)))
	}
	return nil
}
