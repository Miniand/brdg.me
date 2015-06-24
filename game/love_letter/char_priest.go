package love_letter

import "errors"

type CharPriest struct{}

func (p CharPriest) Name() string { return "Priest" }
func (p CharPriest) Number() int  { return Priest }
func (p CharPriest) Text() string {
	return "Look at another player's hand"
}

func (p CharPriest) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
