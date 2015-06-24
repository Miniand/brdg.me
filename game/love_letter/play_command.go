package love_letter

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("play", 1, -1, input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanPlay(pNum)
}

func (c PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)

	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}

	a := command.ExtractNamedCommandArgs(args)
	if len(a) == 0 {
		return "", errors.New("you must specify which card to play")
	}

	card, err := strconv.Atoi(a[0])
	if err != nil {
		names := map[int]string{}
		for _, c := range g.Hands[pNum] {
			names[c] = Cards[c].Name()
		}
		if card, err = helper.MatchStringInStringMap(a[0], names); err != nil {
			return "", err
		}
	}
	return "", g.Play(pNum, card, a[1:]...)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play # (#...){{_b}} to play a card, eg. {{b}}play handmaid{{_b}} or {{b}}play guard steve princess{{b}}"
}

func (g *Game) CanPlay(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Play(player, card int, args ...string) error {
	if !g.CanPlay(player) {
		return errors.New("unable to play right now")
	}
	if _, ok := helper.IntFind(card, g.Hands[player]); !ok {
		return errors.New("you don't have that card")
	}
	if err := Cards[card].Play(g, player, args...); err != nil {
		return err
	}
	curRound := g.Round
	g.DiscardCard(player, card)
	if g.Round == curRound {
		// Only go to the next player if the round didn't just end.
		g.NextPlayer()
	}
	return nil
}
