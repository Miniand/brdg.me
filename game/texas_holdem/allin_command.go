package texas_holdem

import "github.com/Miniand/brdg.me/command"

type AllinCommand struct{}

func (ac AllinCommand) Name() string { return "allin" }

func (ac AllinCommand) Call(
	player string,
	context interface{},
	input *command.Parser,
) (string, error) {
	g := context.(*Game)
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	return "", g.AllIn(playerNum)
}

func (ac AllinCommand) Usage(player string, context interface{}) string {
	return "{{b}}allin{{_b}} to bet all your money and go all in"
}

func (g *Game) CanAllin(player int) bool {
	return g.CurrentPlayer == player && g.PlayerMoney[player] > 0 &&
		!g.IsFinished()
}
