package helper

import "github.com/Miniand/brdg.me/command"

const (
	Mick = iota
	Steve
	BJ
	Greg
	Wallace
	Pete
)

var Players = []string{"Mick", "Steve", "BJ", "Greg", "Wallace", "Pete"}

type Commander interface {
	Commands(player string) []command.Command
}

func Cmd(g Commander, p int, input string) error {
	_, err := command.CallInCommands(
		Players[p],
		g,
		input,
		g.Commands(Players[p]),
	)
	return err
}
