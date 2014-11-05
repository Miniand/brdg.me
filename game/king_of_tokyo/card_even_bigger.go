package king_of_tokyo

import "fmt"

type CardEvenBigger struct{}

func (c CardEvenBigger) Name() string {
	return "Even Bigger"
}

func (c CardEvenBigger) Description() string {
	return fmt.Sprintf(
		"{{b}}Your maximum %s is increased by 2.  Gain %s{{_b}} when you get this card.",
		HealthSymbol,
		RenderHealth(2),
	)
}

func (c CardEvenBigger) Cost() int {
	return 4
}

func (c CardEvenBigger) Kind() int {
	return CardKindKeep
}
