package king_of_tokyo

type CardRegeneration struct{}

func (c CardRegeneration) Name() string {
	return "Regeneration"
}

func (c CardRegeneration) Description() string {
	return "When you heal, {{b}}heal 1 extra damage.{{_b}}"
}

func (c CardRegeneration) Cost() int {
	return 4
}

func (c CardRegeneration) Kind() int {
	return CardKindKeep
}
