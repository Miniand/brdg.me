package texas_holdem

import (
	"github.com/Miniand/brdg.me/command"
)

type FoldCommand struct{}

func (fc FoldCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("fold", 0, input)
}

func (fc FoldCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return false
	}
	currentBet := g.CurrentBet()
	return g.CurrentPlayer == playerNum && g.Bets[playerNum] < currentBet &&
		!g.IsFinished()
}

func (fc FoldCommand) Call(player string, context interface{}, args []string) error {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return err
	}
	return g.Fold(playerNum)
}

func (fc FoldCommand) Usage(player string, context interface{}) string {
	return "{{b}}fold{{_b}} to forfeit this hand"
}
