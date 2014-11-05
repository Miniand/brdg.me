package king_of_tokyo

import "fmt"

type CardSolarPowered struct{}

func (c CardSolarPowered) Name() string {
	return "Solar Powered"
}

func (c CardSolarPowered) Description() string {
	return fmt.Sprintf(
		"At the end of your turn {{b}}gain %s if you have no %s{{_b}}.",
		RenderEnergy(1),
		EnergySymbol,
	)
}

func (c CardSolarPowered) Cost() int {
	return 2
}

func (c CardSolarPowered) Kind() int {
	return CardKindKeep
}
