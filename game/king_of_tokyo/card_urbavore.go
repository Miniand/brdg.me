package king_of_tokyo

import "fmt"

type CardUrbavore struct{}

func (c CardUrbavore) Name() string {
	return "Urbavore"
}

func (c CardUrbavore) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain 1 extra %s{{_b}} when beginning the turn in Tokyo. {{b}}Deal 1 extra damage{{_b}} when dealing any damage from Tokyo.",
		VPSymbol,
	)
}

func (c CardUrbavore) Cost() int {
	return 4
}

func (c CardUrbavore) Kind() int {
	return CardKindKeep
}
