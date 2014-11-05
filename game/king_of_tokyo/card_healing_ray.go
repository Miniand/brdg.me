package king_of_tokyo

import "fmt"

type CardHealingRay struct{}

func (c CardHealingRay) Name() string {
	return "Healing Ray"
}

func (c CardHealingRay) Description() string {
	return fmt.Sprintf(
		"{{b}}You can heal other monsters with your %s results.{{_b}} They must pay you %s for each damage you heal (or their remaining %s if they don't have enough).",
		RenderDie(DieHeal),
		RenderEnergy(2),
		EnergySymbol,
	)
}

func (c CardHealingRay) Cost() int {
	return 4
}

func (c CardHealingRay) Kind() int {
	return CardKindKeep
}
