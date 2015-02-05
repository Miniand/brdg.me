package seven_wonders

import (
	"encoding/gob"

	"github.com/Miniand/brdg.me/game/card"
)

func init() {
	gob.Register(Card{})
}

type Carder interface {
	card.Comparer
	GetCard() Card
}

type CardForPlayers struct {
	Card    string
	Players []int
}

type Card struct {
	Name      string
	Kind      int
	Cost      Cost
	FreeWith  []string
	MakesFree []string
}

func NewCard(
	name string,
	kind int,
	cost Cost,
	freeWith, makesFree []string,
) Card {
	if cost == nil {
		cost = Cost{}
	}
	if freeWith == nil {
		freeWith = []string{}
	}
	if makesFree == nil {
		makesFree = []string{}
	}
	return Card{
		name,
		kind,
		cost,
		freeWith,
		makesFree,
	}
}

func (c Card) Compare(other card.Comparer) (int, bool) {
	oc, ok := other.(Carder)
	if !ok {
		return 0, false
	}
	otherCard := oc.GetCard()
	kindDiff := c.Kind - otherCard.Kind
	if kindDiff != 0 {
		return kindDiff, true
	}
	switch {
	case c.Name < otherCard.Name:
		return -1, true
	case c.Name == otherCard.Name:
		return 0, true
	default:
		return 1, true
	}
}

func (c Card) GetCard() Card {
	return c
}
