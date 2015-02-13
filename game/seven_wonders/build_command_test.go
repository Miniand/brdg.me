package seven_wonders

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

func TestGame_CanBuildCard_free(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{},
		{},
		{},
	}
	can, coins := g.CanBuildCard(Mick, Cards[CardLumberYard])
	assert.True(t, can)
	assert.Len(t, coins, 0)
}

func TestGame_CanBuildCard_prereq(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardTrainingGround]},
		{},
		{},
	}
	can, coins := g.CanBuildCard(Mick, Cards[CardCircus])
	assert.True(t, can)
	assert.Len(t, coins, 0)
}

func TestGame_CanBuildCard_owned(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardLoom]},
		{},
		{},
	}
	can, _ := g.CanBuildCard(Mick, Cards[CardLoom])
	assert.False(t, can)
}

func TestGame_CanBuildCard_self(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardTreeFarm], Cards[CardClayPit], Cards[CardLoom]},
		{},
		{},
	}
	can, coins := g.CanBuildCard(Mick, Cards[CardHaven])
	assert.True(t, can)
	assert.Len(t, coins, 0)
}

func TestGame_CanBuildCard_poor(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardTreeFarm], Cards[CardClayPit], Cards[CardLoom]},
		{},
		{},
	}
	can, _ := g.CanBuildCard(Mick, Cards[CardArsenal])
	assert.False(t, can)
}

func TestGame_CanBuildCard_trade(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardClayPit], Cards[CardLoom]},
		{Cards[CardTreeFarm]},
		{},
	}
	can, coins := g.CanBuildCard(Mick, Cards[CardHaven])
	assert.True(t, can)
	assert.Len(t, coins, 1)
	assert.Equal(t, map[int]int{
		DirRight: 2,
	}, coins[0])
}

func TestGame_CanBuildCard_tradePoor(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardClayPit]},
		{Cards[CardTreeFarm]},
		{Cards[CardLoom]},
	}
	can, _ := g.CanBuildCard(Mick, Cards[CardHaven])
	assert.False(t, can)
}

func TestGame_CanBuildCard_tradeDiscount(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))
	g.Cards = []card.Deck{
		{Cards[CardClayPit], Cards[CardEastTradingPost]},
		{Cards[CardTreeFarm]},
		{Cards[CardLoom]},
	}
	can, coins := g.CanBuildCard(Mick, Cards[CardHaven])
	assert.True(t, can)
	assert.Len(t, coins, 1)
	assert.Equal(t, map[int]int{
		DirLeft:  2,
		DirRight: 1,
	}, coins[0])
}
