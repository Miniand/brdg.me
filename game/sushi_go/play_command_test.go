package sushi_go

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestPlayCommand_Call(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
	assert.Equal(t, helper.Players[:3], g.WhoseTurn())

	// Mick plays a card
	mickCard := g.Hands[helper.Mick][0]
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 1"))
	assert.Equal(t, CardPlayed, g.Hands[helper.Mick][0])
	assert.Equal(t, []int{mickCard}, g.Playing[helper.Mick])
	assert.Equal(t, helper.Players[1:3], g.WhoseTurn())

	// BJ plays a card
	bjCard := g.Hands[helper.BJ][1]
	assert.NoError(t, helper.Cmd(g, helper.BJ, "play 2"))
	assert.Equal(t, CardPlayed, g.Hands[helper.BJ][1])
	assert.Equal(t, []int{bjCard}, g.Playing[helper.BJ])
	assert.Equal(t, helper.Players[1:2], g.WhoseTurn())

	// Steve plays a card
	steveHandLen := len(g.Hands[helper.Steve])
	steveCard := g.Hands[helper.Steve][8]
	assert.NoError(t, helper.Cmd(g, helper.Steve, "play 9"))
	assert.Len(t, g.Hands[helper.Steve], steveHandLen-1)
	// Plays should have happened now.
	assert.Equal(t, []int{mickCard}, g.Played[helper.Mick])
	assert.Equal(t, []int{bjCard}, g.Played[helper.BJ])
	assert.Equal(t, []int{steveCard}, g.Played[helper.Steve])
	assert.Equal(t, helper.Players[:3], g.WhoseTurn())
}

func TestPlayCommand_Call_chopsticks(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	// Prepare hands
	g.Hands[helper.Mick] = []int{
		CardDumpling,
		CardMakiRoll3,
		CardMakiRoll2,
		CardMakiRoll1,
	}
	g.Hands[helper.Steve] = []int{
		CardDumpling,
		CardSalmonNigiri,
		CardSquidNigiri,
		CardEggNigiri,
	}
	g.Played[helper.Mick] = []int{
		CardSquidNigiri,
		CardEggNigiri,
		CardDumpling,
	}
	g.Played[helper.Steve] = []int{
		CardPudding,
		CardChopsticks,
		CardSashimi,
	}

	// Mick tries to play two cards but can't without chopsticks
	assert.Error(t, helper.Cmd(g, helper.Mick, "play 1 2"))
	// It should work after giving Mick chopsticks
	g.Played[helper.Mick][1] = CardChopsticks
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 1 2"))

	// Play with the rest
	assert.NoError(t, helper.Cmd(g, helper.Steve, "play 3"))
	assert.NoError(t, helper.Cmd(g, helper.BJ, "play 2"))

	// Make sure Mick got his hand from Steve
	assert.Equal(t, []int{
		CardDumpling,
		CardSalmonNigiri,
		CardEggNigiri,
	}, g.Hands[helper.Mick])
	// Make sure BJ got his hand from Mick
	assert.Equal(t, []int{
		CardMakiRoll2,
		CardMakiRoll1,
		CardChopsticks,
	}, g.Hands[helper.BJ])
}

func TestPlayCommand_Call_dummyPlayTwo(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))

	g.Played[helper.Mick] = []int{CardChopsticks}
	g.Hands[helper.Mick] = []int{CardMakiRoll1, CardMakiRoll2}
	assert.Error(t, helper.Cmd(g, helper.Mick, "play 1 2"))

	// Should be fine if dummy has already had a card played.
	g.Playing[Dummy] = []int{CardMakiRoll1}
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 1 2"))
}
