package king_of_tokyo

import "fmt"

type CardPoisonSpit struct{}

func (c CardPoisonSpit) Name() string {
	return "Poison Spit"
}

func (c CardPoisonSpit) Description() string {
	healthDie := RenderDie(DieHeal)
	return fmt.Sprintf(
		"When you deal damage to Monsters, give them a poison counter. {{b}}Monsters take 1 damage for each poison counter they have at the end of their turn.{{_b}} You can get rid of a poison counter with a %s (that %s does not also heal a damage).",
		healthDie, healthDie,
	)
}

func (c CardPoisonSpit) Cost() int {
	return 4
}

func (c CardPoisonSpit) Kind() int {
	return CardKindKeep
}
