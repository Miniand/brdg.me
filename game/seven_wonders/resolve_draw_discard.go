package seven_wonders

import (
	"fmt"

	"github.com/Miniand/brdg.me/command"
)

type ResolveDrawDiscard struct {
	Player int
}

func (rdd ResolveDrawDiscard) String(player int, g *Game) string {
	if player != rdd.Player {
		return ""
	}
	return fmt.Sprintf(
		"{{b}}Choose a discarded card to take:{{_b}}\n\n%s",
		g.RenderCardList(player, g.Discard, true, false),
	)
}

func (rdd ResolveDrawDiscard) WhoseTurn(g *Game) []string {
	return []string{g.Players[rdd.Player]}
}

func (rdd ResolveDrawDiscard) Commands() []command.Command {
	return []command.Command{
		TakeCommand{},
	}
}
