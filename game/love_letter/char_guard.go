package love_letter

import "errors"

type CharGuard struct{}

func (p CharGuard) Name() string { return "Guard" }
func (p CharGuard) Number() int  { return Guard }
func (p CharGuard) Text() string {
	return "Guess another player's card to eliminate them, except for Guard"
}

func (p CharGuard) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
