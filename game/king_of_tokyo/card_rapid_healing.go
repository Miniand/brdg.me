package king_of_tokyo

import "fmt"

type CardRapidHealing struct{}

func (c CardRapidHealing) Name() string {
	return "Rapid Healing"
}

func (c CardRapidHealing) Description() string {
	return fmt.Sprintf(
		"Spend %s at any time to {{b}}heal 1 damage.{{_b}}",
		RenderEnergy(2),
	)
}

func (c CardRapidHealing) Cost() int {
	return 3
}

func (c CardRapidHealing) Kind() int {
	return CardKindKeep
}
