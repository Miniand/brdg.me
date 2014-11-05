package king_of_tokyo

import "fmt"

type CardStretchy struct{}

func (c CardStretchy) Name() string {
	return "Stretchy"
}

func (c CardStretchy) Description() string {
	return fmt.Sprintf(
		"You can spend %s to {{b}}change one of your dice to any result.{{_b}}",
		RenderEnergy(2),
	)
}

func (c CardStretchy) Cost() int {
	return 3
}

func (c CardStretchy) Kind() int {
	return CardKindKeep
}
