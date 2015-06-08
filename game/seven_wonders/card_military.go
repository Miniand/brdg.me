package seven_wonders

import (
	"encoding/gob"
	"strconv"

	"github.com/Miniand/brdg.me/game/cost"
)

func init() {
	gob.Register(CardMilitary{})
}

type CardMilitary struct {
	Card
	Strength int
}

func NewCardMilitary(
	name string,
	cost cost.Cost,
	strength int,
	freeWith, makesFree []string,
) CardMilitary {
	return CardMilitary{
		NewCard(name, CardKindMilitary, cost, freeWith, makesFree),
		strength,
	}
}

func (c CardMilitary) SuppString() string {
	return RenderResourceWithSymbol(strconv.Itoa(c.Strength), AttackStrength)
}

func (c CardMilitary) AttackStrength() int {
	return c.Strength
}
