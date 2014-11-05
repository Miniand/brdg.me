package king_of_tokyo

type CardCommuterTrain struct{}

func (c CardCommuterTrain) Name() string {
	return "Commuter Train"
}

func (c CardCommuterTrain) Description() string {
	return RenderVPChange(2)
}

func (c CardCommuterTrain) Cost() int {
	return 4
}

func (c CardCommuterTrain) Kind() int {
	return CardKindDiscard
}
