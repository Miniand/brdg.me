package king_of_tokyo

import "fmt"

type CardFreezeTime struct{}

func (c CardFreezeTime) Name() string {
	return "Freeze Time"
}

func (c CardFreezeTime) Description() string {
	oneDie := RenderDie(Die1)
	return fmt.Sprintf(
		"On a turn where you score %s%s%s, {{b}}you can take another turn{{_b}} with one less die.",
		oneDie, oneDie, oneDie,
	)
}

func (c CardFreezeTime) Cost() int {
	return 5
}

func (c CardFreezeTime) Kind() int {
	return CardKindKeep
}
