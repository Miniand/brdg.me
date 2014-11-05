package king_of_tokyo

import "fmt"

type CardItHasAChild struct{}

func (c CardItHasAChild) Name() string {
	return "It Has a Child!"
}

func (c CardItHasAChild) Description() string {
	return fmt.Sprintf(
		"If you are eliminated discard all your cards and lose all your %s. {{b}}Heal to %s and start again.{{_b}}",
		VPSymbol,
		RenderHealth(10),
	)
}

func (c CardItHasAChild) Cost() int {
	return 7
}

func (c CardItHasAChild) Kind() int {
	return CardKindKeep
}
