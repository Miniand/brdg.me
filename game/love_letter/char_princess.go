package love_letter

import "errors"

type CharPrincess struct{}

func (p CharPrincess) Name() string { return "Princess" }
func (p CharPrincess) Number() int  { return Princess }
func (p CharPrincess) Text() string {
	return "You are eliminated if you discard the Princess"
}

func (p CharPrincess) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}

func (p CharPrincess) HandleDiscard() {
	panic("not implemented")
}
