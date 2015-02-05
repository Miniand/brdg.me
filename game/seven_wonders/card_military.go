package seven_wonders

import (
	"fmt"
	"strings"
)

type CardMilitary struct {
	Card
	Strength int
}

func NewCardMilitary(
	name string,
	cost Cost,
	strength int,
	freeWith, makesFree []string,
) CardMilitary {
	return CardMilitary{
		NewCard(name, CardKindMilitary, cost, freeWith, makesFree),
		strength,
	}
}

func (c CardMilitary) SuppString() string {
	return RenderResourceColour(
		strings.TrimSpace(strings.Repeat(
			fmt.Sprintf("%s ", ResourceSymbols[AttackStrength]),
			c.Strength,
		)),
		AttackStrength,
		true,
	)
}
