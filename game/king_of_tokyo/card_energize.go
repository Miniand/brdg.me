package king_of_tokyo

type CardEnergize struct{}

func (c CardEnergize) Name() string {
	return "Energize"
}

func (c CardEnergize) Description() string {
	return RenderVPChange(9)
}

func (c CardEnergize) Cost() int {
	return 8
}

func (c CardEnergize) Kind() int {
	return CardKindDiscard
}
