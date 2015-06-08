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

func (b PlayerBoard) CurrencyValue(currency int) int {
	value := 0
	for _, c := range b.Cards {
		crd := c.(Card)
		if crd.Currency == currency {
			value += crd.Value
		}
	}
	return value
}
