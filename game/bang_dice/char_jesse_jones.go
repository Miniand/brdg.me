package bang_dice

import "fmt"

type CharJesseJones struct{}

func (c CharJesseJones) Name() string {
	return "Jesse Jones"
}

func (c CharJesseJones) Description() string {
	return fmt.Sprintf(
		"If you have four life points or less, you game two if you use %s for yourself.",
		RenderDie(DieBeer),
	)
}

func (c CharJesseJones) StartingLife() int {
	return 9
}
