package seven_wonders

import (
	"encoding/gob"

	"github.com/Miniand/brdg.me/game/cost"
)

type CardVP struct {
	Card
	VP int
}

func init() {
	gob.Register(CardVP{})
}

func NewCardCivilian(
	name string,
	cost cost.Cost,
	vp int,
	freeWith, makesFree []string,
) CardVP {
	return CardVP{
		NewCard(name, CardKindCivilian, cost, freeWith, makesFree),
		vp,
	}
}

func NewCardWonderVP(
	name string,
	cost cost.Cost,
	vp int,
) CardVP {
	return CardVP{
		NewCard(name, CardKindWonder, cost, nil, nil),
		vp,
	}
}

func (c CardVP) SuppString() string {
	return RenderVP(c.VP)
}

func (c CardVP) VictoryPoints(player int, g *Game) int {
	return c.VP
}
