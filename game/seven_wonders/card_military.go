package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

type CardMilitary struct {
	Card
	Strength int
}

func NewCardMilitary(
	name string,
	cost Cost,
	strength int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	for _, c := range NewCard(
		name,
		CardKindMilitary,
		cost,
		freeWith,
		makesFree,
		players...,
	) {
		d = d.Push(CardMilitary{
			c.(Card),
			strength,
		})
	}
	return nil
}
