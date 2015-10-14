package seven_wonders_duel

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

var WonderInputRegexp = regexp.MustCompile(`(?i)\s+with\s+`)

type WonderCommand struct{}

func (c WonderCommand) Name() string { return "wonder" }

func (c WonderCommand) Call(
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

	parts := WonderInputRegexp.Split(strings.TrimSpace(search), 2)
	if len(parts) != 2 {
		return "", errors.New("you must specify a wonder name and a card to use, eg. wonder pyramids with quarry")
	}

	wonders := []string{}
	for _, w := range g.PlayerWonders[pNum] {
		wonders = append(wonders, Cards[w].Name)
	}
	wonder, err := helper.MatchStringInStrings(
		parts[0],
		wonders,
	)
	if err != nil {
		return "", err
	}

	buildableLocs := g.Layout.Buildable()
	names := make([]string, len(buildableLocs))
	for k, l := range buildableLocs {
		names[k] = Cards[g.Layout.At(l)].Name
	}
	card, err := helper.MatchStringInStrings(
		parts[1],
		names,
	)
	if err != nil {
		return "", err
	}
	return "", g.Wonder(pNum, wonder, buildableLocs[card])
}

func (c WonderCommand) Usage(player string, context interface{}) string {
	return "{{b}}wonder [card]{{_b}} to wonder a card for 2 coins + the number of yellow cards you have, eg. {{b}}wonder quarry{{_b}}"
}

func (g *Game) CanWonder(player int) bool {
	return len(g.PlayerWonders[player]) > 0 && g.CanPlay(player)
}

func (g *Game) Wonder(player, wonder int, loc Loc) error {
	if !g.CanWonder(player) {
		return errors.New("can't build a wonder at the moment")
	}
	if !g.Layout.CanBuild(loc) {
		return errors.New("that card isn't available")
	}
	if wonder < 0 || wonder >= len(g.PlayerWonders[player]) {
		return errors.New("invalid wonder number")
	}
	g.PlayerCards[player] = append(
		g.PlayerCards[player],
		g.PlayerWonders[player][wonder],
	)
	g.PlayerWonders[player] = append(
		g.PlayerWonders[player][:wonder],
		g.PlayerWonders[player][wonder+1:]...,
	)
	g.CurrentPlayer = Opponent(g.CurrentPlayer)
	g.Layout = g.Layout.Remove(loc)
	if len(g.Layout) == 0 {
		g.NextAge()
	}
	return nil
}
