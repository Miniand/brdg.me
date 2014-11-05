package king_of_tokyo

import "fmt"

type CardGasRefinery struct{}

func (c CardGasRefinery) Name() string {
	return "Gas Refinery"
}

func (c CardGasRefinery) Description() string {
	return fmt.Sprintf(
		"{{b}}%s and deal 3 damage to all other monsters.{{_b}}",
		RenderVPChange(2),
	)
}

func (c CardGasRefinery) Cost() int {
	return 6
}

func (c CardGasRefinery) Kind() int {
	return CardKindDiscard
}
