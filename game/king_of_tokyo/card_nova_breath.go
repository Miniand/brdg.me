package king_of_tokyo

type CardNovaBreath struct{}

func (c CardNovaBreath) Name() string {
	return "Nova Breath"
}

func (c CardNovaBreath) Description() string {
	return "{{b}}Your attacks damage all other monsters.{{_b}}"
}

func (c CardNovaBreath) Cost() int {
	return 7
}

func (c CardNovaBreath) Kind() int {
	return CardKindKeep
}
