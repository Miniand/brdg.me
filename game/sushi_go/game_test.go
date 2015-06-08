package sushi_go

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
}

func TestGame_Encode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	_, err := g.Encode()
	assert.NoError(t, err)
}

func TestGame_Decode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	data, err := g.Encode()
	assert.NoError(t, err)
	assert.NoError(t, g.Decode(data))
}

func TestGame_Score_maki(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	score, _ := g.Score()
	assert.Equal(t, []int{0, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardMakiRoll1}
	score, _ = g.Score()
	assert.Equal(t, []int{6, 0, 0}, score)

	g.Played[helper.Steve] = []int{CardMakiRoll1}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 3, 0}, score)

	g.Played[helper.BJ] = []int{CardMakiRoll1}
	score, _ = g.Score()
	assert.Equal(t, []int{2, 2, 2}, score)

	g.Played[helper.Steve] = []int{CardMakiRoll2}
	score, _ = g.Score()
	assert.Equal(t, []int{1, 6, 1}, score)
}

func TestGame_Score_pudding(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	g.Played[helper.Mick] = []int{CardPudding}
	score, _ := g.Score()
	assert.Equal(t, []int{0, 0, 0}, score)

	g.Round = 3
	score, _ = g.Score()
	assert.Equal(t, []int{6, -3, -3}, score)

	g.Played[helper.BJ] = []int{CardPudding, CardPudding}
	score, _ = g.Score()
	assert.Equal(t, []int{0, -6, 6}, score)

	g.Played[helper.Mick] = []int{CardPudding, CardPudding}
	g.Played[helper.Steve] = []int{CardPudding, CardPudding}
	score, _ = g.Score()
	assert.Equal(t, []int{0, 0, 0}, score)
}

func TestGame_Score_nigiri(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	g.Played[helper.Mick] = []int{CardEggNigiri}
	score, _ := g.Score()
	assert.Equal(t, []int{1, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardEggNigiri, CardWasabi}
	score, _ = g.Score()
	assert.Equal(t, []int{1, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardWasabi, CardEggNigiri}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 0, 0}, score)

	g.Played[helper.Steve] = []int{CardSalmonNigiri}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 2, 0}, score)

	g.Played[helper.Steve] = []int{CardWasabi, CardSalmonNigiri}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 6, 0}, score)

	g.Played[helper.BJ] = []int{CardSquidNigiri}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 6, 3}, score)

	g.Played[helper.BJ] = []int{CardWasabi, CardSquidNigiri}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 6, 9}, score)
}

func TestGame_Score_tempura(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	g.Played[helper.Mick] = []int{CardTempura}
	score, _ := g.Score()
	assert.Equal(t, []int{0, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardTempura, CardTempura}
	score, _ = g.Score()
	assert.Equal(t, []int{5, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardTempura, CardTempura, CardTempura}
	score, _ = g.Score()
	assert.Equal(t, []int{5, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardTempura, CardTempura, CardTempura, CardTempura}
	score, _ = g.Score()
	assert.Equal(t, []int{10, 0, 0}, score)
}

func TestGame_Score_sashimi(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	g.Played[helper.Mick] = []int{CardSashimi}
	score, _ := g.Score()
	assert.Equal(t, []int{0, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardSashimi, CardSashimi}
	score, _ = g.Score()
	assert.Equal(t, []int{0, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardSashimi, CardSashimi, CardSashimi}
	score, _ = g.Score()
	assert.Equal(t, []int{10, 0, 0}, score)

	g.Played[helper.Mick] = []int{CardSashimi, CardSashimi, CardSashimi, CardSashimi}
	score, _ = g.Score()
	assert.Equal(t, []int{10, 0, 0}, score)
}

func TestGame_Score_dumpling(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))

	g.Played[helper.Mick] = []int{
		CardDumpling,
	}
	score, _ := g.Score()
	assert.Equal(t, []int{1, 0, 0}, score)

	g.Played[helper.Mick] = []int{
		CardDumpling,
		CardDumpling,
	}
	score, _ = g.Score()
	assert.Equal(t, []int{3, 0, 0}, score)

	g.Played[helper.Mick] = []int{
		CardDumpling,
		CardDumpling,
		CardDumpling,
	}
	score, _ = g.Score()
	assert.Equal(t, []int{6, 0, 0}, score)

	g.Played[helper.Mick] = []int{
		CardDumpling,
		CardDumpling,
		CardDumpling,
		CardDumpling,
	}
	score, _ = g.Score()
	assert.Equal(t, []int{10, 0, 0}, score)

	g.Played[helper.Mick] = []int{
		CardDumpling,
		CardDumpling,
		CardDumpling,
		CardDumpling,
		CardDumpling,
	}
	score, _ = g.Score()
	assert.Equal(t, []int{15, 0, 0}, score)

	g.Played[helper.Mick] = []int{
		CardDumpling,
		CardDumpling,
		CardDumpling,
		CardDumpling,
		CardDumpling,
		CardDumpling,
	}
	score, _ = g.Score()
	assert.Equal(t, []int{15, 0, 0}, score)
}
