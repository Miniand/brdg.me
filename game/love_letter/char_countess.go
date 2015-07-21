package love_letter

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type CharCountess struct{}

func (p CharCountess) Name() string { return "Countess" }
func (p CharCountess) Number() int  { return Countess }
func (p CharCountess) Text() string {
	return "Discard the Countess if you have the King or Prince in your hand"
}
func (p CharCountess) Colour() string { return render.Red }

func (p CharCountess) Play(g *Game, player int, args ...string) error {
	g.DiscardCard(player, Countess)

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded %s, they might have been forced to if they also had %s or %s",
		g.RenderName(player),
		RenderCard(Countess),
		RenderCard(King),
		RenderCard(Prince),
	)))
	return nil
}
