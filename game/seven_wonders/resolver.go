package seven_wonders

import "github.com/Miniand/brdg.me/command"

type Resolver interface {
	String(player int, g *Game) string
	WhoseTurn(g *Game) []string
	Commands() []command.Command
}
