package king_of_tokyo

import "fmt"

type CardPsychicProbe struct{}

func (c CardPsychicProbe) Name() string {
	return "Psychic Probe"
}

func (c CardPsychicProbe) Description() string {
	return fmt.Sprintf(
		"You can {{b}}reroll a die of each other monster once each turn.{{_b}} If the reroll is %s discard this card.",
		RenderDie(DieHeal),
	)
}

func (c CardPsychicProbe) Cost() int {
	return 3
}

func (c CardPsychicProbe) Kind() int {
	return CardKindKeep
}
