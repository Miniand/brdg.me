package seven_wonders

import "github.com/Miniand/brdg.me/command"

const (
	Mick = iota
	Steve
	Greg
)

var players = []string{
	"Mick",
	"Steve",
	"Greg",
}

func cmd(g *Game, player int, input string) error {
	_, err := command.CallInCommands(
		players[player],
		g,
		input,
		g.Commands(players[player]),
	)
	return err
}
