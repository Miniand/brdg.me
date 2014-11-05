package king_of_tokyo

import "fmt"

type CardParasiticTentacles struct{}

func (c CardParasiticTentacles) Name() string {
	return "Parasitic Tentacles"
}

func (c CardParasiticTentacles) Description() string {
	return fmt.Sprintf(
		"{{b}}You can purchase cards from other monsters.{{_b}} Pay them the %s cost.",
		EnergySymbol,
	)
}

func (c CardParasiticTentacles) Cost() int {
	return 4
}

func (c CardParasiticTentacles) Kind() int {
	return CardKindKeep
}
