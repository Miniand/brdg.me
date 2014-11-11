package bang_dice

import "fmt"

type CharBlackJack struct{}

func (c CharBlackJack) Name() string {
	return "Black Jack"
}

func (c CharBlackJack) Description() string {
	return fmt.Sprintf(
		"You may re-roll %s. (Not if you roll three ore more!)",
		DieStrings[DieDynamite],
	)
}

func (c CharBlackJack) StartingLife() int {
	return 8
}
