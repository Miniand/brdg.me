package king_of_tokyo

type CardJets struct{}

func (c CardJets) Name() string {
	return "Jets"
}

func (c CardJets) Description() string {
	return "{{b}}You suffer no damage{{_b}} when yielding Tokyo."
}

func (c CardJets) Cost() int {
	return 5
}

func (c CardJets) Kind() int {
	return CardKindKeep
}
