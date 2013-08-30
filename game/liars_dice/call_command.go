package liars_dice

import (
	"github.com/Miniand/brdg.me/command"
)

type CallCommand struct{}

func (c CallCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("call", 0, input)
}

func (c CallCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return !g.IsFinished() && g.WhoseTurn()[0] == player && g.BidQuantity != 0
}

func (c CallCommand) Call(player string, context interface{}, args []string) (
	output string, err error) {
	g := context.(*Game)
	quantity := 0
	for _, pd := range g.PlayerDice {
		for _, d := range pd {
			if d == g.BidValue || d == 1 {
				quantity++
			}
		}
	}
	if quantity < g.BidQuantity {
		// Caller was correct
		g.PlayerDice[g.BidPlayer] = g.PlayerDice[g.BidPlayer][1:]
	} else {
		// Bidder was correct
		g.PlayerDice[g.CurrentPlayer] = g.PlayerDice[g.CurrentPlayer][1:]
	}
	if !g.IsFinished() {
		g.StartRound()
		g.CurrentPlayer = g.NextActivePlayer(g.CurrentPlayer)
	}
	return
}

func (c CallCommand) Usage(player string, context interface{}) string {
	return "{{b}}call{{_b}} to call the last bidder if you think their bid is too high."
}
