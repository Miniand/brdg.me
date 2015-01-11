package king_of_tokyo

import "github.com/Miniand/brdg.me/command"

const (
	Mick = iota
	Steve
	BJ
	Walas
)

var names = []string{"Mick", "Steve", "BJ", "Walas"}

func cmd(g *Game, player int, input string) error {
	_, err := command.CallInCommands(g.Players[player], g, input, g.Commands())
	return err
}
