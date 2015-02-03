package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

type CardCivilian struct {
	Card
	VP int
}

func NewCardCivilian(
	name string,
	cost Cost,
	vp int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	for _, c := range NewCard(
		name,
		CardKindCivilian,
		cost,
		freeWith,
		makesFree,
		players...,
	) {
		d = d.Push(CardCivilian{
			c.(Card),
			vp,
		})
	}
	return nil
}
