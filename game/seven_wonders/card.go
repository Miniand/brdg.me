package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

const (
	CardKindRaw = iota
	CardKindManufactured
	CardKindCivilian
	CardKindScientific
	CardKindCommercial
	CardKindMilitary
	CardKindGuild
	TokenVictory
	TokenDefeat
	WonderStage
)

type Carder interface {
	GetCard() Card
}

type Cost map[int]int

type Card struct {
	Name       string
	Kind       int
	Cost       Cost
	FreeWith   []string
	MakesFree  []string
	MinPlayers int
}

func NewCard(
	name string,
	kind int,
	cost Cost,
	freeWith, makesFree []string,
	players ...int,
) card.Deck {
	d := card.Deck{}
	if cost == nil {
		cost = Cost{}
	}
	if freeWith == nil {
		freeWith = []string{}
	}
	if makesFree == nil {
		makesFree = []string{}
	}
	for _, p := range players {
		d = d.Push(Card{
			name,
			kind,
			cost,
			freeWith,
			makesFree,
			p,
		})
	}
	return d
}

func (c Card) Compare(other card.Comparer) (int, bool) {
	oc, ok := other.(Card)
	if !ok {
		return 0, false
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
