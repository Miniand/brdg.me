package king_of_tokyo

import "fmt"

type CardWings struct{}

func (c CardWings) Name() string {
	return "Wings"
}

func (c CardWings) Description() string {
	return fmt.Sprintf(
		`{{b}}Spend %s to negate damage to you{{_b}} for a turn.`,
		RenderEnergyChange(2),
	)
}

func (c CardWings) Cost() int {
	return 6
}

func (c CardWings) Kind() int {
	return CardKindKeep
}
