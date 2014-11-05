package king_of_tokyo

import "fmt"

type CardRootingForTheUnderdog struct{}

func (c CardRootingForTheUnderdog) Name() string {
	return "Rooting for the Underdog"
}

func (c CardRootingForTheUnderdog) Description() string {
	return fmt.Sprintf(
		"At the end of a turn, if you have the fewest %s, {{b}}gain %s.{{_b}}",
		VPSymbol,
		RenderVP(1),
	)
}

func (c CardRootingForTheUnderdog) Cost() int {
	return 3
}

func (c CardRootingForTheUnderdog) Kind() int {
	return CardKindKeep
}
