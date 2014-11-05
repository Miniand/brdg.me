package king_of_tokyo

import "fmt"

type CardTanks struct{}

func (c CardTanks) Name() string {
	return "Tanks"
}

func (c CardTanks) Description() string {
	return fmt.Sprintf(
		"{{b}}%s and take 3 damage.{{_b}}",
		RenderVPChange(4),
	)
}

func (c CardTanks) Cost() int {
	return 4
}

func (c CardTanks) Kind() int {
	return CardKindDiscard
}
