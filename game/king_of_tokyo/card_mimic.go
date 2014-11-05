package king_of_tokyo

import "fmt"

type CardMimic struct{}

func (c CardMimic) Name() string {
	return "Mimic"
}

func (c CardMimic) Description() string {
	return fmt.Sprintf(
		"{{b}}Choose a card any monster has in play{{_b}} and put a mimic counter on it. {{b}}This card counts as a duplicate of that card as if it just had been bought.{{_b}} Spend %s at the start of your turn to change the power you are mimicking.",
		RenderEnergy(1),
	)
}

func (c CardMimic) Cost() int {
	return 8
}

func (c CardMimic) Kind() int {
	return CardKindKeep
}
