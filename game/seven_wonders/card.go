package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

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
	oc, ok := other.(Card)
	if !ok {
		return 0, false
	}
	kindDiff := c.Kind - oc.Kind
	if kindDiff != 0 {
		return kindDiff, true
	}
	switch {
	case c.Name < oc.Name:
		return -1, true
	case c.Name == oc.Name:
		return 0, true
	default:
		return 1, true
	}
}

func (c Card) GetCard() Card {
	return c
}
