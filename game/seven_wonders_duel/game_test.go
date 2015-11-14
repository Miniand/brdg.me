package seven_wonders_duel

import (
	"fmt"
	"testing"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func atFirstRound(t *testing.T) *Game {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	for g.Phase == PhaseChooseWonder {
		var p int
		for pNum, pName := range g.Players {
			if len(g.Commands(pName)) > 0 {
				p = pNum
				break
			}
		}
		assert.NoError(t, helper.Cmd(g, p, fmt.Sprintf(
			"choose %s",
			Cards[g.AvailableWonders()[0]].Name,
		)))
	}
	return g
}

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
}

func TestOpponent(t *testing.T) {
	assert.Equal(t, 1, Opponent(0))
	assert.Equal(t, 0, Opponent(1))
}

func TestGame_Decode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	data, err := g.Encode()
	assert.NoError(t, err)
	newG := &Game{}
	assert.NoError(t, newG.Decode(data))
}

func TestGame_TradeResourceCount(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.True(t, g.TradeGoodCount(helper.Mick).IsZero())
	assert.True(t, g.TradeGoodCount(helper.Steve).IsZero())

	g.PlayerCards[helper.Mick] = []int{
		CardSawmill,
		CardLoggingCamp,
		CardGlassblower,
		CardForum,
	}
	assert.Equal(t, cost.Cost{
		GoodWood:  3,
		GoodGlass: 1,
	}, g.TradeGoodCount(helper.Mick))
}

func TestGame_TradePrices(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.Equal(t, cost.Cost{
		GoodWood:    2,
		GoodStone:   2,
		GoodClay:    2,
		GoodGlass:   2,
		GoodPapyrus: 2,
	}, BaseTradePrices())
	assert.Equal(t, BaseTradePrices(), g.TradePrices(helper.Mick))
	assert.Equal(t, BaseTradePrices(), g.TradePrices(helper.Steve))

	g.PlayerCards[helper.Mick] = []int{
		CardSawmill,
		CardLoggingCamp,
		CardGlassblower,
		CardForum,
	}
	assert.Equal(t, cost.Cost{
		GoodWood:    5,
		GoodStone:   2,
		GoodClay:    2,
		GoodGlass:   3,
		GoodPapyrus: 2,
	}, g.TradePrices(helper.Steve))

	g.PlayerCards[helper.Steve] = []int{
		CardWoodReserve,
		CardCustomsHouse,
	}
	assert.Equal(t, cost.Cost{
		GoodWood:    1,
		GoodStone:   2,
		GoodClay:    2,
		GoodGlass:   1,
		GoodPapyrus: 1,
	}, g.TradePrices(helper.Steve))
}
