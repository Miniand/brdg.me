package king_of_tokyo

import "fmt"

type CardEnergyHoarder struct{}

func (c CardEnergyHoarder) Name() string {
	return "Energy Hoarder"
}

func (c CardEnergyHoarder) Description() string {
	return fmt.Sprintf(
		"{{b}}You gain %s{{_b}} for every %s you have at the end of your turn.",
		RenderVP(1),
		RenderEnergy(6),
	)
}

func (c CardEnergyHoarder) Cost() int {
	return 3
}

func (c CardEnergyHoarder) Kind() int {
	return CardKindKeep
}
