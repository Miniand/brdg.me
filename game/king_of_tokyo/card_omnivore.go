package king_of_tokyo

import "fmt"

type CardOmnivore struct{}

func (c CardOmnivore) Name() string {
	return "Omnivore"
}

func (c CardOmnivore) Description() string {
	return fmt.Sprintf(
		"{{b}}Once on each turn you can score %s%s%s for %s.{{_b}} You can also use these dice in other combinations.",
		RenderDie(Die1),
		RenderDie(Die2),
		RenderDie(Die3),
		RenderVP(2),
	)
}

func (c CardOmnivore) Cost() int {
	return 4
}

func (c CardOmnivore) Kind() int {
	return CardKindKeep
}
