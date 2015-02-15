package seven_wonders

import (
	"encoding/gob"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/cost"
)

func init() {
	gob.Register(Card{})
}

type Carder interface {
	card.Comparer
	GetCard() Card
}

type PreActionExecuteHandler interface {
	HandlePreActionExecute(player int, g *Game)
}

type PostActionExecuteHandler interface {
	HandlePostActionExecute(player int, g *Game)
}

type VictoryPointer interface {
	VictoryPoints(player int, g *Game) int
}

type Attacker interface {
	AttackStrength() int
}

type Sciencer interface {
	ScienceField(player int, g *Game) int
}

type PlayFinalCarder interface {
	PlayFinalCard() bool
}

type CardForPlayers struct {
	Card    string
	Players []int
}

type WonderStager interface {
	WonderStages() int
}

type Card struct {
	Name      string
	Kind      int
	Cost      cost.Cost
	FreeWith  []string
	MakesFree []string
}

func NewCard(
	name string,
	kind int,
	c cost.Cost,
	freeWith, makesFree []string,
) Card {
	if c == nil {
		c = cost.Cost{}
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
		c,
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
