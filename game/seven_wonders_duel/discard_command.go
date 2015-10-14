package seven_wonders_duel

import (
	"errors"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type DiscardCommand struct{}

func (c DiscardCommand) Name() string { return "discard" }

func (c DiscardCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	search, err := input.ReadToEndOfLine()
	if err != nil {
		return "", err
	}

	buildableLocs := g.Layout.Buildable()
	names := make([]string, len(buildableLocs))
	for k, l := range buildableLocs {
		names[k] = Cards[g.Layout.At(l)].Name
	}

	card, err := helper.MatchStringInStrings(
		strings.TrimSpace(search),
		names,
	)
	if err != nil {
		return "", err
	}
	return "", g.Discard(pNum, buildableLocs[card])
}

func (c DiscardCommand) Usage(player string, context interface{}) string {
	return "{{b}}discard [card]{{_b}} to discard a card for 2 coins + the number of yellow cards you have, eg. {{b}}discard quarry{{_b}}"
}

func (g *Game) CanDiscard(player int) bool {
	return g.CanPlay(player)
}

func (g *Game) Discard(player int, loc Loc) error {
	if !g.CanDiscard(player) {
		return errors.New("can't discard at the moment")
	}
	if !g.Layout.CanBuild(loc) {
		return errors.New("that card isn't available")
	}
	g.DiscardedCards = append(g.DiscardedCards, g.Layout.At(loc))
	g.ModifyCoins(player, 2+g.PlayerCardTypeCount(player, CardTypeCommercial))
	g.CurrentPlayer = Opponent(g.CurrentPlayer)
	g.Layout = g.Layout.Remove(loc)
	if len(g.Layout) == 0 {
		g.NextAge()
	}
	return nil
}
