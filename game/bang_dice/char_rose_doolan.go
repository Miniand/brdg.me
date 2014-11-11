package bang_dice

import "fmt"

type CharRoseDoolan struct{}

func (c CharRoseDoolan) Name() string {
	return "Rose Doolan"
}

func (c CharRoseDoolan) Description() string {
	return fmt.Sprintf(
		"You may use %s or %s for players sitting one place further",
		DieStrings[Die1],
		DieStrings[Die2],
	)
}

func (c CharRoseDoolan) StartingLife() int {
	return 9
}
