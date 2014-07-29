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
	FriendOfThePeople   bool
	HeroOfThePeople     bool
}

func NewPlayerBoard(player int) *PlayerBoard {
	pb := &PlayerBoard{
		Player: player,
		Resources: map[int]int{
			ResourceTrade:      2,
			ResourceScience:    1,
			ResourceAstro:      25,
			ResourceColonyShip: 1,
			ResourceTradeShip:  1,
			ResourceBooster:    2,
			ResourceCannon:     1,
		},
		Modules:             map[int]int{},
		CompletedAdventures: card.Deck{},
		Colonies:            card.Deck{},
		TradingPosts:        card.Deck{},
	}
	pb.Colonies = pb.Colonies.Push(StartingCards()[player])
	return pb
}
