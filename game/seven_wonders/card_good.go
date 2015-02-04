package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

const (
	GoodCoin = iota
	GoodWood
	GoodStone
	GoodOre
	GoodClay
	GoodPapyrus
	GoodTextile
	GoodGlass
)

var RawGoods = []int{
	GoodWood,
	GoodStone,
	GoodOre,
	GoodClay,
}

var ManufacturedGoods = []int{
	GoodPapyrus,
	GoodTextile,
	GoodGlass,
}

type CardGood struct {
	Card
	Goods  []int
	Amount int
}

func NewCardGood(
	name string,
	kind int,
	cost Cost,
	goods []int,
	amount int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	if goods == nil || len(goods) == 0 {
		panic("no goods")
	}
	for _, c := range NewCard(
		name,
		kind,
		cost,
		freeWith,
		makesFree,
		players...,
	) {
		d = d.Push(CardGood{
			c.(Card),
			goods,
			amount,
		})
	}
	return d
}

func NewCardGoodRaw(
	name string,
	cost Cost,
	goods []int,
	amount int,
	players ...int,
) card.Deck {
	return NewCardGood(
		name,
		CardKindRaw,
		cost,
		goods,
		amount,
		nil,
		nil,
		players...,
	)
}

func NewCardGoodManufactured(
	name string,
	cost Cost,
	goods []int,
	amount int,
	players ...int,
) card.Deck {
	return NewCardGood(
		name,
		CardKindManufactured,
		cost,
		goods,
		amount,
		nil,
		nil,
		players...,
	)
}

func NewCardGoodCommercial(
	name string,
	cost Cost,
	goods []int,
	amount int,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	return NewCardGood(
		name,
		CardKindCommercial,
		cost,
		goods,
		amount,
		freeWith,
		makesFree,
		players...,
	)
}
