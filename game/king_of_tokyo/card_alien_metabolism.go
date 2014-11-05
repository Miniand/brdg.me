package king_of_tokyo

import "fmt"

type CardAlienMetabolism struct{}

func (c CardAlienMetabolism) Name() string {
	return "Alien Metabolism"
}

func (c CardAlienMetabolism) Description() string {
	return fmt.Sprintf(
		`{{b}}Buying cards costs you 1 less %s.`,
		EnergySymbol,
	)
}

func (c CardAlienMetabolism) Cost() int {
	return 3
}

func (c CardAlienMetabolism) Kind() int {
	return CardKindKeep
}
