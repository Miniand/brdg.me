package king_of_tokyo

type CardAcidAttack struct{}

func (c CardAcidAttack) Name() string {
	return "Acid Attack"
}

func (c CardAcidAttack) Description() string {
	return "{{b}}Deal 1 extra damage each turn{{_b}} (even when you don't otherwise attack)."
}

func (c CardAcidAttack) Cost() int {
	return 6
}

func (c CardAcidAttack) Kind() int {
	return CardKindKeep
}
