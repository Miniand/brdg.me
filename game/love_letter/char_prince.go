package love_letter

import "errors"

type CharPrince struct{}

func (p CharPrince) Name() string { return "Prince" }
func (p CharPrince) Number() int  { return Prince }
func (p CharPrince) Text() string {
	return "Choose a player to discard and draw a new card"
}

func (p CharPrince) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
