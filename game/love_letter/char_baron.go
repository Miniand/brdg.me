package love_letter

import "errors"

type CharBaron struct{}

func (p CharBaron) Name() string { return "Baron" }
func (p CharBaron) Number() int  { return Baron }
func (p CharBaron) Text() string {
	return "Compare hands with another player, lowest card is eliminated"
}

func (p CharBaron) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
