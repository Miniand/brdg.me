package king_of_tokyo

type CardHeal struct{}

func (c CardHeal) Name() string {
	return "Heal"
}

func (c CardHeal) Description() string {
	return "{{b}}Heal 2 damage{{_b}}."
}

func (c CardHeal) Cost() int {
	return 3
}

func (c CardHeal) Kind() int {
	return CardKindDiscard
}