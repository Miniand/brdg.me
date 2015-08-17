package texas_holdem

import (
	"github.com/Miniand/brdg.me/command"
)

type FoldCommand struct{}

func (fc FoldCommand) Name() string { return "fold" }

func (fc FoldCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Fold(playerNum)
}

func (fc FoldCommand) Usage(player string, context interface{}) string {
	return "{{b}}fold{{_b}} to forfeit this hand"
}

func (g *Game) CanFold(player int) bool {
	currentBet := g.CurrentBet()
	return g.CurrentPlayer == player && g.Bets[player] < currentBet &&
		!g.IsFinished()
}
