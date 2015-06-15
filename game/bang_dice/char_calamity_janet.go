package bang_dice

import "fmt"

type CharCalamityJanet struct{}

func (c CharCalamityJanet) Name() string {
	return "Calamity Janet"
}

func (c CharCalamityJanet) Description() string {
	return fmt.Sprintf(
		"You can use %s as %s and vice-versa.",
		RenderDie(Die1),
		RenderDie(Die2),
	)
}

func (c CharCalamityJanet) StartingLife() int {
	return 8
}
