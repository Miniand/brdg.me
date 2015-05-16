package jaipur

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestSellCommand_Call(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.NoError(t, helper.Cmd(g, helper.Mick, "sell 2 gold"))
}
