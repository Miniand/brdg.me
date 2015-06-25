package love_letter

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
)

type CharHandmaid struct{}

func (p CharHandmaid) Name() string { return "Handmaid" }
func (p CharHandmaid) Number() int  { return Handmaid }
func (p CharHandmaid) Text() string {
	return "Immune to the effects of other players' cards until next turn"
}

func (p CharHandmaid) Play(g *Game, player int, args ...string) error {
	if len(args) > 0 {
		return errors.New("Handmaid doesn't accept any arguments when playing")
	}

	g.DiscardCard(player, Handmaid)

	g.Protected[player] = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s and is immune to the effects of other players' cards until the start of their next turn",
		g.RenderName(player),
		RenderCard(Handmaid),
	)))

	return nil
}
