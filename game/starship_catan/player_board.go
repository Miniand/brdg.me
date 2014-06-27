package starship_catan

import "github.com/Miniand/brdg.me/game/card"

type PlayerBoard struct {
	Player              int
	Astro               int
	CompletedAdventures card.Deck
	TradeShips          int
	ColonyShips         int
	Modules             map[int]int
}
