package seven_wonders

import (
	"encoding/gob"

	"github.com/Miniand/brdg.me/game/cost"
)

type CardCivilian struct {
	Card
	VP int
}

func init() {
	gob.Register(CardCivilian{})
}

func NewCardCivilian(
	name string,
	cost cost.Cost,
	vp int,
	freeWith, makesFree []string,
) CardCivilian {
	return CardCivilian{
		NewCard(name, CardKindCivilian, cost, freeWith, makesFree),
		vp,
	}
}

func (c CardCivilian) SuppString() string {
	return RenderVP(c.VP)
}

func (c CardCivilian) VictoryPoints(player int, g *Game) int {
	return c.VP
}
