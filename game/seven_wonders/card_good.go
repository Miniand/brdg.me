package seven_wonders

import (
	"encoding/gob"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
)

func init() {
	gob.Register(CardGood{})
}

type GoodsProducer interface {
	GoodsProduced() []cost.Cost
}

type GoodsTrader interface {
	GoodsTraded() []cost.Cost
}

type CardGood struct {
	Card
	Goods []cost.Cost
}

func NewCardGood(
	name string,
	kind int,
	cost cost.Cost,
	goods []cost.Cost,
	freeWith, makesFree []string,
) CardGood {
	if goods == nil || len(goods) == 0 {
		panic("no goods")
	}
	return CardGood{
		NewCard(name, kind, cost, freeWith, makesFree),
		goods,
	}
}

func NewCardGoodRaw(
	name string,
	cost cost.Cost,
	goods []cost.Cost,
) CardGood {
	return NewCardGood(
		name,
		CardKindRaw,
		cost,
		goods,
		nil,
		nil,
	)
}

func NewCardGoodManufactured(
	name string,
	cost cost.Cost,
	goods []cost.Cost,
) CardGood {
	return NewCardGood(
		name,
		CardKindManufactured,
		cost,
		goods,
		nil,
		nil,
	)
}

func NewCardGoodCommercial(
	name string,
	cost cost.Cost,
	goods []cost.Cost,
	freeWith, makesFree []string,
) CardGood {
	return NewCardGood(
		name,
		CardKindCommercial,
		cost,
		goods,
		freeWith,
		makesFree,
	)
}

func (c CardGood) SuppString() string {
	parts := []string{}
	for _, g := range c.Goods {
		parts = append(parts, RenderResourceList(g.Ints(), " "))
	}
	return strings.Join(parts, "/")
}

func (c CardGood) GoodsProduced() []cost.Cost {
	return c.Goods
}

func (c CardGood) GoodsTraded() []cost.Cost {
	if c.Card.Kind == CardKindCommercial {
		// Commercial cards don't trade to neighbours.
		return []cost.Cost{}
	}
	return c.Goods
}
