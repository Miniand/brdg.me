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
	Commands() []command.Command
}

func Cmd(g Commander, p int, input string) error {
	_, err := command.CallInCommands(
		Players[p],
		g,
		input,
		command.AvailableCommands(Players[p], g, g.Commands()),
	)
	return err
}
