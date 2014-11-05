package king_of_tokyo

import "fmt"

type CardMetamorph struct{}

func (c CardMetamorph) Name() string {
	return "Metamorph"
}

func (c CardMetamorph) Description() string {
	return fmt.Sprintf(
		"At the end of your turn you can {{b}}discard any Keep cards you have to receive the %s{{_b}} they were purchased for.",
		EnergySymbol,
	)
}

func (c CardMetamorph) Cost() int {
	return 3
}

func (c CardMetamorph) Kind() int {
	return CardKindKeep
}
