package king_of_tokyo

import "fmt"

type CardJetFighters struct{}

func (c CardJetFighters) Name() string {
	return "Jet Fighters"
}

func (c CardJetFighters) Description() string {
	return fmt.Sprintf(
		"{{b}}%s and take 4 damage.{{_b}}",
		RenderVPChange(5),
	)
}

func (c CardJetFighters) Cost() int {
	return 5
}

func (c CardJetFighters) Kind() int {
	return CardKindDiscard
}
