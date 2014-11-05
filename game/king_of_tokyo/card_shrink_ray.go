package king_of_tokyo

import "fmt"

type CardShrinkRay struct{}

func (c CardShrinkRay) Name() string {
	return "Shrink Ray"
}

func (c CardShrinkRay) Description() string {
	healDie := RenderDie(DieHeal)
	return fmt.Sprintf(
		"When you deal damage to Monsters, give them a shrink counter. {{b}}A monster rolls one less die for each shrink counter.{{_b}} You can get rid of a shrink counter with a %s (that %s does not also heal damage).",
		healDie, healDie,
	)
}

func (c CardShrinkRay) Cost() int {
	return 6
}

func (c CardShrinkRay) Kind() int {
	return CardKindKeep
}
