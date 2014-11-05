package king_of_tokyo

import "fmt"

type CardCamouflage struct{}

func (c CardCamouflage) Name() string {
	return "Camouflage"
}

func (c CardCamouflage) Description() string {
	return fmt.Sprintf(
		"If you take damage roll a die for each damage point. {{b}}On a %s you do not take that damage point.{{_b}}",
		RenderDie(DieHeal),
	)
}

func (c CardCamouflage) Cost() int {
	return 3
}

func (c CardCamouflage) Kind() int {
	return CardKindKeep
}
