package king_of_tokyo

type CardApartmentBuilding struct{}

func (c CardApartmentBuilding) Name() string {
	return "Apartment Building"
}

func (c CardApartmentBuilding) Description() string {
	return RenderVPChange(3)
}

func (c CardApartmentBuilding) Cost() int {
	return 5
}

func (c CardApartmentBuilding) Kind() int {
	return CardKindDiscard
}
