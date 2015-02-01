package alhambra

import "github.com/Miniand/brdg.me/game/card"

type PlayerBoard struct {
	Grid    Grid
	Reserve []Tile
	Cards   card.Deck
	Place   []Tile
	Points  int
}

func NewPlayerBoard() PlayerBoard {
	return PlayerBoard{
		Grid:    NewGrid(),
		Reserve: []Tile{},
		Cards:   card.Deck{},
		Place:   []Tile{},
	}
}

func (b PlayerBoard) TileCounts() map[int]int {
	counts := map[int]int{}
	for _, t := range b.Grid {
		if t.Type == TileTypeEmpty {
			continue
		}
		counts[t.Type]++
	}
	return counts
}
