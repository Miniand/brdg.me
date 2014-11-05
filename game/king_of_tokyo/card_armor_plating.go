package king_of_tokyo

type CardArmorPlating struct{}

func (c CardArmorPlating) Name() string {
	return "Armor Plating"
}

func (c CardArmorPlating) Description() string {
	return "{{b}}Ignore damage when it is only 1 point.{{_b}}"
}

func (c CardArmorPlating) Cost() int {
	return 4
}

func (c CardArmorPlating) Kind() int {
	return CardKindKeep
}
