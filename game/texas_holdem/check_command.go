package texas_holdem

import "github.com/Miniand/brdg.me/command"

type CheckCommand struct{}

func (cc CheckCommand) Name() string { return "check" }

func (cc CheckCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.Check(playerNum)
}

func (cc CheckCommand) Usage(player string, context interface{}) string {
	return "{{b}}check{{_b}} to continue without betting more money"
}

func (g *Game) CanCheck(player int) bool {
	currentBet := g.CurrentBet()
	return g.CurrentPlayer == player && g.Bets[player] == currentBet &&
		!g.IsFinished()
}
