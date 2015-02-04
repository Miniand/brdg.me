package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

const (
	DirLeft  = -1
	DirDown  = 0
	DirRight = 1
)

var (
	DirAll        = []int{DirLeft, DirDown, DirRight}
	DirNeighbours = []int{DirLeft, DirRight}
	DirSelf       = []int{DirDown}
)

type CardCommercialTrade struct {
	Card
	Directions []int
	Goods      []int
}

func NewCardCommercialTrade(
	name string,
	cost Cost,
	directions, goods []int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	for _, c := range NewCard(
		name,
		CardKindCommercial,
		cost,
		freeWith,
		makesFree,
		players...,
	) {
		d = d.Push(CardCommercialTrade{
			c.(Card),
			directions,
			goods,
		})
	}
	return d
}

type CardCommercialTavern struct {
	Card
}

func NewCardCommercialTavern(players ...int) card.Deck {
	d := card.Deck{}
	for _, c := range NewCard(
		CardTavern,
		CardKindCommercial,
		nil,
		nil,
		nil,
		players...,
	) {
		d = d.Push(CardCommercialTavern{c.(Card)})
	}
	return d
}
