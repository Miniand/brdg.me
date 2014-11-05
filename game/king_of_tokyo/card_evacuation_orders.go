package king_of_tokyo

import "fmt"

type CardEvacuationOrders struct{}

func (c CardEvacuationOrders) Name() string {
	return "Evacuation Orders"
}

func (c CardEvacuationOrders) Description() string {
	return fmt.Sprintf(
		"{{b}}All other monsters lose %s{{_b}}.",
		RenderVP(5),
	)
}

func (c CardEvacuationOrders) Cost() int {
	return 7
}

func (c CardEvacuationOrders) Kind() int {
	return CardKindDiscard
}
