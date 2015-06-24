package love_letter

import "errors"

type CharHandmaid struct{}

func (p CharHandmaid) Name() string { return "Handmaid" }
func (p CharHandmaid) Number() int  { return Handmaid }
func (p CharHandmaid) Text() string {
	return "Immune to the effects of other players' cards until next turn"
}

func (p CharHandmaid) Play(g *Game, player int, args ...string) error {
	return errors.New("not implemented")
}
