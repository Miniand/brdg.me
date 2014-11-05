package king_of_tokyo

type CardSkyscraper struct{}

func (c CardSkyscraper) Name() string {
	return "Skyscraper"
}

func (c CardSkyscraper) Description() string {
	return RenderVPChange(4)
}

func (c CardSkyscraper) Cost() int {
	return 6
}

func (c CardSkyscraper) Kind() int {
	return CardKindDiscard
}
