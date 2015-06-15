package bang_dice

import "fmt"

type CharKitCarlson struct{}

func (c CharKitCarlson) Name() string {
	return "Kit Carlson"
}

func (c CharKitCarlson) Description() string {
	return fmt.Sprintf(
		"For each %s you may discard one arrow from any player.",
		RenderDie(DieGatling),
	)
}

func (c CharKitCarlson) StartingLife() int {
	return 7
}
