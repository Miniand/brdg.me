package texas_holdem

import (
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type CallCommand struct{}

func (cc CallCommand) Name() string { return "call" }

func (cc CallCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Call(playerNum)
}

func (cc CallCommand) Usage(player string, context interface{}) string {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(
		"{{b}}call{{_b}} to increase your bet by %d to match the current bet",
		g.CurrentBet()-g.Bets[playerNum])
}

func (g *Game) CanCall(player int) bool {
	currentBet := g.CurrentBet()
	return g.CurrentPlayer == player && g.Bets[player] < currentBet &&
		g.PlayerMoney[player] > currentBet-g.Bets[player] &&
		!g.IsFinished()
}
