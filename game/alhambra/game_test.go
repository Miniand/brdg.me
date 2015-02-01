package alhambra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_ScoreType(t *testing.T) {
	g := &Game{}
	g.Start(players)

	g.Boards[Mick].Grid = Grid{
		Vect{0, 1}: NewTile(TileTypePavillion, 0),
		Vect{0, 2}: NewTile(TileTypeSeraglio, 0),
		Vect{0, 3}: NewTile(TileTypeSeraglio, 0),
		Vect{0, 4}: NewTile(TileTypeArcades, 0),
		Vect{0, 5}: NewTile(TileTypeChambers, 0),
		Vect{0, 6}: NewTile(TileTypeChambers, 0),
	}

	g.Boards[Steve].Grid = Grid{
		Vect{0, 1}: NewTile(TileTypeArcades, 0),
		Vect{0, 2}: NewTile(TileTypeChambers, 0),
		Vect{0, 3}: NewTile(TileTypeSeraglio, 0),
		Vect{0, 4}: NewTile(TileTypeTower, 0),
		Vect{0, 5}: NewTile(TileTypeArcades, 0),
		Vect{0, 6}: NewTile(TileTypeArcades, 0),
		Vect{0, 7}: NewTile(TileTypeChambers, 0),
	}

	g.Boards[BJ].Grid = Grid{
		Vect{0, 1}: NewTile(TileTypeGarden, 0),
		Vect{0, 2}: NewTile(TileTypeTower, 0),
		Vect{0, 3}: NewTile(TileTypeArcades, 0),
		Vect{0, 4}: NewTile(TileTypeArcades, 0),
		Vect{0, 5}: NewTile(TileTypeChambers, 0),
	}

	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick}, 1},
	}, g.ScoreType(TileTypePavillion, 1))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick}, 8},
	}, g.ScoreType(TileTypePavillion, 2))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick}, 16},
	}, g.ScoreType(TileTypePavillion, 3))

	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick}, 2},
	}, g.ScoreType(TileTypeSeraglio, 1))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick}, 9},
		{[]int{Steve}, 2},
	}, g.ScoreType(TileTypeSeraglio, 2))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick}, 17},
		{[]int{Steve}, 9},
	}, g.ScoreType(TileTypeSeraglio, 3))

	assert.Equal(t, []RoundTypeScore{
		{[]int{Steve}, 3},
	}, g.ScoreType(TileTypeArcades, 1))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Steve}, 10},
		{[]int{BJ}, 3},
	}, g.ScoreType(TileTypeArcades, 2))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Steve}, 18},
		{[]int{BJ}, 10},
		{[]int{Mick}, 3},
	}, g.ScoreType(TileTypeArcades, 3))

	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick, Steve}, 2},
	}, g.ScoreType(TileTypeChambers, 1))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick, Steve}, 7},
	}, g.ScoreType(TileTypeChambers, 2))
	assert.Equal(t, []RoundTypeScore{
		{[]int{Mick, Steve}, 15},
		{[]int{BJ}, 4},
	}, g.ScoreType(TileTypeChambers, 3))
}
