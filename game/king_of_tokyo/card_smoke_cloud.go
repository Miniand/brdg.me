package king_of_tokyo

type CardSmokeCloud struct{}

func (c CardSmokeCloud) Name() string {
	return "Smoke Cloud"
}

func (c CardSmokeCloud) Description() string {
	return "Place 3 Smoke counters on this card. {{b}}Spend 1 Smoke counter for an extra reroll.{{_b}} Discard this card when all Smoke counters are spent."
}

func (c CardSmokeCloud) Cost() int {
	return 4
}

func (c CardSmokeCloud) Kind() int {
	return CardKindKeep
}
