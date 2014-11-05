package king_of_tokyo

import "fmt"

type CardHerdCuller struct{}

func (c CardHerdCuller) Name() string {
	return "Herd Culler"
}

func (c CardHerdCuller) Description() string {
	return fmt.Sprintf(
		"You can {{b}}change one of your dice to a %s{{_b}} each turn.",
		RenderDie(Die1),
	)
}

func (c CardHerdCuller) Cost() int {
	return 3
}

func (c CardHerdCuller) Kind() int {
	return CardKindKeep
}
