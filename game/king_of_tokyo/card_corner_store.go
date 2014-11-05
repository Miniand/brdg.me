package king_of_tokyo

type CardCornerStore struct{}

func (c CardCornerStore) Name() string {
	return "Corner Store"
}

func (c CardCornerStore) Description() string {
	return RenderVPChange(1)
}

func (c CardCornerStore) Cost() int {
	return 3
}

func (c CardCornerStore) Kind() int {
	return CardKindDiscard
}
