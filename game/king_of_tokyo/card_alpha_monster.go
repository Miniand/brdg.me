package king_of_tokyo

import "fmt"

type CardAlphaMonster struct{}

func (c CardAlphaMonster) Name() string {
	return "Alpha Monster"
}

func (c CardAlphaMonster) Description() string {
	return fmt.Sprintf(
		"{{b}}Gain %s{{_b}} when you attack.",
		RenderVP(1),
	)
}

func (c CardAlphaMonster) Cost() int {
	return 5
}

func (c CardAlphaMonster) Kind() int {
	return CardKindKeep
}
