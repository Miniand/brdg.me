package love_letter

import "errors"

type CharKing struct{}

func (p CharKing) Name() string { return "King" }
func (p CharKing) Number() int  { return King }
func (p CharKing) Text() string {
	return "Trade your hand with another player"
}

func (p CharKing) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
