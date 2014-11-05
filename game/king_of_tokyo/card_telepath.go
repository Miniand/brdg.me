package king_of_tokyo

import "fmt"

type CardTelepath struct{}

func (c CardTelepath) Name() string {
	return "Telepath"
}

func (c CardTelepath) Description() string {
	return fmt.Sprintf(
		"Spend %s to {{b}}get 1 extra reroll.{{_b}}",
		RenderEnergy(1),
	)
}

func (c CardTelepath) Cost() int {
	return 4
}

func (c CardTelepath) Kind() int {
	return CardKindKeep
}
