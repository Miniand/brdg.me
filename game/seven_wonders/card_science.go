package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

const (
	FieldMathematics = iota
	FieldEngineering
	FieldTheology
)

var AllFields = []int{
	FieldMathematics,
	FieldEngineering,
	FieldTheology,
}

type CardScience struct {
	Card
	Fields []int
}

func NewCardScience(
	name string,
	cost Cost,
	fields []int,
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
		d = d.Push(CardScience{
			c.(Card),
			fields,
		})
	}
	return nil
}
