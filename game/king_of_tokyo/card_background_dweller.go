package king_of_tokyo

import "fmt"

type CardBackgroundDweller struct{}

func (c CardBackgroundDweller) Name() string {
	return "Background Dweller"
}

func (c CardBackgroundDweller) Description() string {
	return fmt.Sprintf(
		"{{b}}You can always reroll any %s{{_b}} you have.",
		RenderDie(Die3),
	)
}

func (c CardBackgroundDweller) Cost() int {
	return 4
}

func (c CardBackgroundDweller) Kind() int {
	return CardKindKeep
}
