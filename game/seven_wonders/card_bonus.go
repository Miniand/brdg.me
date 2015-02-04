package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

type CardBonus struct {
	Card
	TargetKinds []int
	Directions  []int
	VP          int
	Coins       int
}

func NewCardBonus(
	name string,
	kind int,
	cost Cost,
	targetKinds, directions []int,
	vp, coins int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	if targetKinds == nil || len(targetKinds) == 0 {
		panic("no targetKinds")
	}
	if directions == nil || len(directions) == 0 {
		panic("no directions")
	}
	for _, c := range NewCard(
		name,
		kind,
		cost,
		freeWith,
		makesFree,
		players...,
	) {
		d = d.Push(CardBonus{
			c.(Card),
			targetKinds,
			directions,
			vp,
			coins,
		})
	}
	return d
}
