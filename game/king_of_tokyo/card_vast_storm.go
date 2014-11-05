package king_of_tokyo

import "fmt"

type CardVastStorm struct{}

func (c CardVastStorm) Name() string {
	return "Vast Storm"
}

func (c CardVastStorm) Description() string {
	return fmt.Sprintf(
		"{{b}}%s. All other monsters lose %s for every %s{{_b}} they have.",
		RenderVPChange(2),
		RenderEnergy(1),
		RenderEnergy(2),
	)
}

func (c CardVastStorm) Cost() int {
	return 6
}

func (c CardVastStorm) Kind() int {
	return CardKindDiscard
}
