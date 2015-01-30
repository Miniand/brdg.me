package alhambra

import "github.com/Miniand/brdg.me/game/card"

type PlayerBoard struct {
	Grid    Grid
	Reserve []Tile
	Cards   card.Deck
}

func NewPlayerBoard() PlayerBoard {
	return PlayerBoard{
		Grid:    NewGrid(),
		Reserve: []Tile{},
		Cards:   card.Deck{},
	}
}
