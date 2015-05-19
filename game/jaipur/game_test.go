package jaipur

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	assert.Len(t, g.Deck, 40)
	assert.Equal(t, 5, len(g.Hands[helper.Mick])+g.Camels[helper.Mick])
	assert.Equal(t, 5, len(g.Hands[helper.Steve])+g.Camels[helper.Steve])
	assert.Len(t, g.Goods, len(TradeGoods))
	assert.Len(t, g.Bonuses, 3)
	assert.Len(t, g.Bonuses[3], 7)
	assert.Len(t, g.Bonuses[4], 6)
	assert.Len(t, g.Bonuses[5], 5)
}

func TestGame_Decode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	data, err := g.Encode()
	assert.NoError(t, err)
	assert.NoError(t, g.Decode(data))
}

func TestParseGoodStr(t *testing.T) {
	for _, test := range []struct {
		Input    string
		Expected []int
	}{
		{"Camel", []int{GoodCamel}},
		{"2 Camels", []int{GoodCamel, GoodCamel}},
		{"2 Camels -1 dia", nil},
		{"2 gold gold di", []int{GoodDiamond, GoodGold, GoodGold, GoodGold}},
	} {
		actual, err := ParseGoodStr(test.Input)
		if test.Expected == nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, helper.IntSort(test.Expected), helper.IntSort(actual))
		}
	}
}
