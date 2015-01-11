package king_of_tokyo

type CardMadeInALab struct{}

func (c CardMadeInALab) Name() string {
	return "Made in a Lab"
}

func (c CardMadeInALab) Description() string {
	return "When purchasing cards you can {{b}}peek at and purchase the top card{{_b}} of the deck."
}

func (c CardMadeInALab) Cost() int {
	return 2
}

func (c CardMadeInALab) Kind() int {
	return CardKindKeep
}

func (c CardMadeInALab) ModifyBuyable(
	game *Game,
	player int,
	buyable []CardBase,
) []CardBase {
	if player != game.CurrentPlayer || game.Phase != PhaseBuy ||
		len(game.Deck) == 0 {
		return buyable
	}
	extraCard := game.Deck[0]
	// We only allow peeking the first card, even with mimic, so check if it's
	// already there.
	for _, c := range buyable {
		if c == extraCard {
			return buyable
		}
	}
	return append(buyable, game.Deck[0])
}
