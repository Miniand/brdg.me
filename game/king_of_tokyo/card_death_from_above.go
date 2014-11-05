package king_of_tokyo

import "fmt"

type CardDeathFromAbove struct{}

func (c CardDeathFromAbove) Name() string {
	return "Death From Above"
}

func (c CardDeathFromAbove) Description() string {
	return fmt.Sprintf(
		"%s {{b}}and take control of Tokyo{{_b}} if you don't already have it",
		RenderVPChange(2),
	)
}

func (c CardDeathFromAbove) Cost() int {
	return 5
}

func (c CardDeathFromAbove) Kind() int {
	return CardKindDiscard
}
