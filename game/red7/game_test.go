package red7

import (
	"log"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func crd(input string) int {
	card, ok := ParseCard(input)
	if !ok {
		log.Fatalf("could not parse card %s", input)
	}
	return card
}

func crds(inputs ...string) []int {
	out := make([]int, len(inputs))
	for i, in := range inputs {
		out[i] = crd(in)
	}
	return out
}

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.Len(t, g.Hands[helper.Mick], 7)
	assert.Len(t, g.Hands[helper.Steve], 7)
	assert.Len(t, g.Palettes[helper.Mick], 1)
	assert.Len(t, g.Palettes[helper.Steve], 1)
}

func TestGame_CurrentRule(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.Equal(t, SuitRed, g.CurrentRule())
	g.DiscardPile = append(g.DiscardPile, crd("b5"))
	assert.Equal(t, SuitBlue, g.CurrentRule())
	g.DiscardPile = append(g.DiscardPile, crd("y5"))
	assert.Equal(t, SuitYellow, g.CurrentRule())
}

func TestGame_Decode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	data, err := g.Encode()
	assert.NoError(t, err)
	g2 := &Game{}
	assert.NoError(t, g2.Decode(data))
}

func TestGame_Leader(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	g.Palettes = [][]int{
		// Mick
		crds("y3"),
		// Steve
		crds("b4"),
	}
	leader, palette := g.Leader()
	assert.Equal(t, helper.Steve, leader)
	assert.Equal(t, crds("b4"), palette)
}

func TestGame_NextPlayer(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:4]))
	assert.Equal(t, helper.Steve, g.NextPlayer(helper.Mick))
	g.Eliminated[helper.Mick] = true
	assert.Equal(t, helper.Steve, g.NextPlayer(helper.Greg))
	g.Eliminated[helper.Steve] = true
	assert.Equal(t, helper.BJ, g.NextPlayer(helper.Greg))
}

func TestEndPoints(t *testing.T) {
	assert.Equal(t, 40, EndPoints(2))
	assert.Equal(t, 35, EndPoints(3))
	assert.Equal(t, 30, EndPoints(4))
}

func TestGame_EndRound(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	// We expect that the deck will be 1 less after the card is scored.
	initialLen := len(g.Deck)
	assert.NoError(t, helper.Cmd(g, g.CurrentPlayer, "done"))
	assert.Len(t, g.Deck, initialLen-1)
}
