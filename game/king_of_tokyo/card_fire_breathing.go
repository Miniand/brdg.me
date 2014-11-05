package king_of_tokyo

type CardFireBreathing struct{}

func (c CardFireBreathing) Name() string {
	return "Fire Breathing"
}

func (c CardFireBreathing) Description() string {
	return "{{b}}Your neighbours take 1 extra damage{{_b}} when you deal damage."
}

func (c CardFireBreathing) Cost() int {
	return 4
}

func (c CardFireBreathing) Kind() int {
	return CardKindKeep
}
