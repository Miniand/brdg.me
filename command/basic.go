package command

import (
	"fmt"
	"github.com/beefsack/brdg.me/game"
	"regexp"
	"strings"
)

const ARGUMENT_REGEXP = `\s*(\b[^\s]+\b)`

type BasicCommand struct {
	Name        string
	NumArgs     int
	CanCallFunc func(player string, g *game.Playable) bool
	CallFunc    func(player string, g *game.Playable, args []string) error
	UsageFunc   func(player string, g *game.Playable) string
}

func (b BasicCommand) Parse(input string) []string {
	return regexp.MustCompile(fmt.Sprintf(`(?im)\s*%s`+strings.Repeat(
		ARGUMENT_REGEXP, b.NumArgs)+`\s*$`, b.Name)).FindStringSubmatch(input)
}

func (b BasicCommand) CanCall(player string, g *game.Playable) bool {
	return b.CanCallFunc(player, g)
}

func (b BasicCommand) Call(player string, g *game.Playable, args []string) error {
	return b.CallFunc(player, g, args[1:])
}

func (b BasicCommand) Usage(player string, g *game.Playable) string {
	if b.UsageFunc == nil {
		return b.Name
	}
	return b.UsageFunc(player, g)
}
