package love_letter

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CharCountess struct{}

func (p CharCountess) Name() string { return "Countess" }
func (p CharCountess) Number() int  { return Countess }
func (p CharCountess) Text() string {
	return "If you have the King or Prince in your hand, you must discard the Countess"
}

func (p CharCountess) Play(g *Game, player int, args ...string) error {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded %s, they might have been forced to if they also had %s or %s",
		g.RenderName(player),
		RenderCard(Countess),
		RenderCard(King),
		RenderCard(Prince),
	)))
	return nil
}
