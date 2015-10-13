package seven_wonders_duel

import (
	"errors"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type PlayCommand struct{}

func (c PlayCommand) Name() string { return "play" }

func (c PlayCommand) Call(
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
	return "", g.Play(pNum, buildableLocs[card])
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play [card]{{_b}} to build a card, eg. {{b}}play {{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	if g.Phase != PhasePlay {
		return false
	}
	return g.CurrentPlayer == player
}

func (g *Game) Play(player int, loc Loc) error {
	if !g.CanPlay(player) {
		return errors.New("can't play at the moment")
	}
	if !g.Layout.CanBuild(loc) {
		return errors.New("that card isn't playable")
	}
	g.PlayerCards[player] = append(g.PlayerCards[player], g.Layout.At(loc))
	g.CurrentPlayer = Opponent(g.CurrentPlayer)
	g.Layout = g.Layout.Remove(loc)
	if len(g.Layout) == 0 {
		g.NextAge()
	}
	return nil
}
