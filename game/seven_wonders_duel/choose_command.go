package seven_wonders_duel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

type ChooseCommand struct{}

func (c ChooseCommand) Name() string { return "choose" }

func (c ChooseCommand) Call(
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
	availableWonders := g.AvailableWonders()
	wonders := make([]string, len(availableWonders))
	for i, w := range availableWonders {
		wonders[i] = Cards[w].Name
	}
	wonder, err := helper.MatchStringInStrings(
		strings.TrimSpace(search),
		wonders,
	)
	if err != nil {
		return "", err
	}
	return "", g.ChooseWonder(pNum, wonder)
}

func (c ChooseCommand) Usage(player string, context interface{}) string {
	return "{{b}}choose [wonder]{{_b}} to choose a wonder, eg. {{b}}choose colossus{{_b}}"
}

func (g *Game) AvailableWonders() []int {
	n := len(g.RemainingWonders) % 4
	if n == 0 {
		n = 4
	}
	return g.RemainingWonders[:n]
}

func (g *Game) CanChoose(player int) bool {
	if g.Phase != PhaseChooseWonder {
		return false
	}
	// Player order is 0 1 1 0 1 0 0 1
	return (byte(105)&(byte(1)<<(byte(len(g.RemainingWonders))-1)) == 0) !=
		(player == 1)
}

func (g *Game) ChooseWonder(player, wonder int) error {
	if !g.CanChoose(player) {
		return errors.New("can't choose at the moment")
	}
	if wonder < 0 || wonder > len(g.RemainingWonders) {
		return errors.New("not a valid wonder")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s chose
%s`,
		g.PlayerName(player),
		Cards[g.RemainingWonders[wonder]].RenderMultiline(0),
	)))
	g.PlayerWonders[player] = append(
		g.PlayerWonders[player],
		g.RemainingWonders[wonder],
	)
	g.RemainingWonders = append(
		g.RemainingWonders[:wonder],
		g.RemainingWonders[wonder+1:]...,
	)
	if len(g.RemainingWonders) == 0 {
		g.StartAge()
	} else if len(g.AvailableWonders()) == 1 {
		g.ChooseWonder(Opponent(player), 0)
	}
	return nil
}
