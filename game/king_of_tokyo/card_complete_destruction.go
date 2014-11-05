package king_of_tokyo

import "fmt"

type CardCompleteDestruction struct{}

func (c CardCompleteDestruction) Name() string {
	return "Complete Destruction"
}

func (c CardCompleteDestruction) Description() string {
	return fmt.Sprintf(
		"If you roll %s%s%s%s%s%s {{b}}gain %s{{_b}} in addition to the regular results.",
		RenderDie(Die1),
		RenderDie(Die2),
		RenderDie(Die3),
		RenderDie(DieHeal),
		RenderDie(DieAttack),
		RenderDie(DieEnergy),
		RenderVP(9),
	)
}

func (c CardCompleteDestruction) Cost() int {
	return 3
}

func (c CardCompleteDestruction) Kind() int {
	return CardKindKeep
}
