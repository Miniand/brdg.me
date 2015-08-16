package splendor

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type TakeCommand struct{}

func (c TakeCommand) Name() string { return "take" }

func (c TakeCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	pNum, found := g.PlayerNum(player)
	if !found {
		return "", errors.New("could not find player")
	}
	args, err := input.ReadLineArgs()
	if err != nil || len(args) == 0 {
		return "", errors.New("please specify two or three tokens")
	}
	tokens := []int{}
	gemStrings := GemStrings()
	for _, a := range args {
		t, err := helper.MatchStringInStringMap(a, gemStrings)
		if err != nil {
			return "", err
		}
		tokens = append(tokens, t)
	}
	return "", g.Take(pNum, tokens)
}

func (c TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take ## ## (##){{_b}} to take two or three tokens, eg. {{b}}take di di{{_b}}.  If you take two you must take two of the same type of tokens, and there must be at least four in the supply.  If you take three, they must be three different tokens."
}

func (g *Game) CanTake(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseMain
}

func (g *Game) Take(player int, tokens []int) error {
	if !g.CanTake(player) {
		return errors.New("unable to take right now")
	}
	switch l := len(tokens); l {
	case 2:
		if tokens[0] != tokens[1] {
			return errors.New("must take the same type of tokens when taking two")
		}
		if g.Tokens[tokens[0]] < 4 {
			return errors.New("can only take two when there are four or more remaining")
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s took {{b}}2 %s{{_b}}",
			g.RenderName(player),
			RenderResourceColour(ResourceStrings[tokens[0]], tokens[0]),
		)))
	case 3:
		tokenStrs := []string{}
		for i, t := range tokens {
			if t == tokens[(i+1)%l] {
				return errors.New("must take different tokens when taking three")
			}
			if g.Tokens[t] == 0 {
				return errors.New("there aren't enough tokens remaning to take that")
			}
			tokenStrs = append(tokenStrs, render.Bold(
				RenderResourceColour(ResourceStrings[t], t)))
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s took %s",
			g.RenderName(player),
			render.CommaList(tokenStrs),
		)))
	default:
		return errors.New("can only take two or three tokens")
	}
	amount := cost.Cost{}
	for _, t := range tokens {
		amount[t] += 1
	}
	g.PlayerBoards[player].Tokens = g.PlayerBoards[player].Tokens.Add(amount)
	g.Tokens = g.Tokens.Sub(amount)
	g.NextPhase()
	return nil
}
