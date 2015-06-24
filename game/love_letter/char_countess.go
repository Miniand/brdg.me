package love_letter

import "errors"

type CharCountess struct{}

func (p CharCountess) Name() string { return "Countess" }
func (p CharCountess) Number() int  { return Countess }
func (p CharCountess) Text() string {
	return "If you have the King or Prince in your hand, you must discard the Countess"
}

func (p CharCountess) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
