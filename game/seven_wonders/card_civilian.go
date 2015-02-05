package seven_wonders

import "encoding/gob"

type CardCivilian struct {
	Card
	VP int
}

func init() {
	gob.Register(CardCivilian{})
}

func NewCardCivilian(
	name string,
	cost Cost,
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
