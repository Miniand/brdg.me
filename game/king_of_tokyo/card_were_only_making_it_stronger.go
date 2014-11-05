package king_of_tokyo

import "fmt"

type CardWereOnlyMakingItStronger struct{}

func (c CardWereOnlyMakingItStronger) Name() string {
	return "We're Only Making It Stronger!"
}

func (c CardWereOnlyMakingItStronger) Description() string {
	return fmt.Sprintf(
		"When you lose %s or more {{b}}gain %s.{{_b}}",
		RenderHealth(2),
		RenderEnergy(1),
	)
}

func (c CardWereOnlyMakingItStronger) Cost() int {
	return 3
}

func (c CardWereOnlyMakingItStronger) Kind() int {
	return CardKindKeep
}
