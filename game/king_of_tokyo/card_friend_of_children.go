package king_of_tokyo

import "fmt"

type CardFriendOfChildren struct{}

func (c CardFriendOfChildren) Name() string {
	return "Friend of Children"
}

func (c CardFriendOfChildren) Description() string {
	return fmt.Sprintf(
		"When you gain any %s {{b}}gain 1 extra %s.{{_b}}",
		EnergySymbol, EnergySymbol,
	)
}

func (c CardFriendOfChildren) Cost() int {
	return 3
}

func (c CardFriendOfChildren) Kind() int {
	return CardKindKeep
}
