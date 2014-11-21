package king_of_tokyo

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/stretchr/testify/assert"
)

const (
	Mick = iota
	Steve
	BJ
)

var names = []string{"Mick", "Steve", "BJ"}

func cmd(t *testing.T, g *Game, player int, input string) {
	_, err := command.CallInCommands(g.Players[player], g, input, g.Commands())
	assert.NoError(t, err)
}
