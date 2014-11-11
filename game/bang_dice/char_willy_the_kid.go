package bang_dice

import "fmt"

type CharWillyTheKid struct{}

func (c CharWillyTheKid) Name() string {
	return "Willy the Kid"
}

func (c CharWillyTheKid) Description() string {
	return fmt.Sprintf(
		"You only need %s%s to use the Gatling Gun.",
		DieStrings[DieGatling],
		DieStrings[DieGatling],
	)
}

func (c CharWillyTheKid) StartingLife() int {
	return 8
}
