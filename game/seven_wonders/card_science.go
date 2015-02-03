package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

const (
	FieldMathematics = iota
	FieldEngineering
	FieldTheology
)

type CardScience struct {
	Card
	Field int
}

func NewCardScience(
	name string,
	cost Cost,
	field int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	for _, c := range NewCard(
		name,
		CardKindScientific,
		cost,
		freeWith,
		makesFree,
		players...,
	) {
		d = d.Push(CardMilitary{
			c.(Card),
			field,
		})
	}
	return nil
}
