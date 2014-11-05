package king_of_tokyo

import "fmt"

type CardGourmet struct{}

func (c CardGourmet) Name() string {
	return "Gourmet"
}

func (c CardGourmet) Description() string {
	oneDie := RenderDie(Die1)
	return fmt.Sprintf(
		"When scoring %s%s%s, {{b}}gain 2 extra %s.{{_b}}",
		oneDie, oneDie, oneDie,
		VPSymbol,
	)
}

func (c CardGourmet) Cost() int {
	return 4
}

func (c CardGourmet) Kind() int {
	return CardKindKeep
}
