package king_of_tokyo

type CardFireBlast struct{}

func (c CardFireBlast) Name() string {
	return "Fire Blast"
}

func (c CardFireBlast) Description() string {
	return "{{b}}Deal 2 damage to all other monsters.{{_b}}"
}

func (c CardFireBlast) Cost() int {
	return 3
}

func (c CardFireBlast) Kind() int {
	return CardKindDiscard
}
