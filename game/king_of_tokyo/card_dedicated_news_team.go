package king_of_tokyo

import "fmt"

type CardDedicatedNewsTeam struct{}

func (c CardDedicatedNewsTeam) Name() string {
	return "Dedicated News Team"
}

func (c CardDedicatedNewsTeam) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} whenever you buy a card.",
		RenderVP(1),
	)
}

func (c CardDedicatedNewsTeam) Cost() int {
	return 3
}

func (c CardDedicatedNewsTeam) Kind() int {
	return CardKindKeep
}
