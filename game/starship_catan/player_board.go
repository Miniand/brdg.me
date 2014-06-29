package starship_catan

import "github.com/Miniand/brdg.me/game/card"

type PlayerBoard struct {
	Player              int
	Resources           map[int]int
	Modules             map[int]int
	CompletedAdventures card.Deck
	Colonies            card.Deck
	TradingPosts        card.Deck
	TradeShips          int
	ColonyShips         int
}

func NewPlayerBoard(player int) *PlayerBoard {
	pb := &PlayerBoard{
		Player: player,
		Resources: map[int]int{
			ResourceTrade:   2,
			ResourceScience: 1,
		},
		Modules:             map[int]int{},
		CompletedAdventures: card.Deck{},
		Colonies:            card.Deck{},
		TradingPosts:        card.Deck{},
	}
	if player == 0 {
	} else {
	}
	return pb
}
