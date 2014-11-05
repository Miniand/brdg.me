package king_of_tokyo

import "fmt"

type CardEaterOfTheDead struct{}

func (c CardEaterOfTheDead) Name() string {
	return "Eater of the Dead"
}

func (c CardEaterOfTheDead) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} every time a monster's %s goes to 0.",
		RenderVP(3),
		HealthSymbol,
	)
}

func (c CardEaterOfTheDead) Cost() int {
	return 4
}

func (c CardEaterOfTheDead) Kind() int {
	return CardKindKeep
}
