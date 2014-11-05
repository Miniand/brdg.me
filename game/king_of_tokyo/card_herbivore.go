package king_of_tokyo

import "fmt"

type CardHerbivore struct{}

func (c CardHerbivore) Name() string {
	return "Herbivore"
}

func (c CardHerbivore) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} on your turn if you don't damage anyone.",
		RenderVP(1),
	)
}

func (c CardHerbivore) Cost() int {
	return 5
}

func (c CardHerbivore) Kind() int {
	return CardKindKeep
}
