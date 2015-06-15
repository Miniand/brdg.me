package bang_dice

import "fmt"

type CharSlabTheKiller struct{}

func (c CharSlabTheKiller) Name() string {
	return "Slab the Killer"
}

func (c CharSlabTheKiller) Description() string {
	return fmt.Sprintf(
		"Once per turn, you can use a %s to double a %s or %s.",
		RenderDie(DieBeer),
		RenderDie(Die1),
		RenderDie(Die2),
	)
}

func (c CharSlabTheKiller) StartingLife() int {
	return 8
}
