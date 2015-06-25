package love_letter

import (
	"fmt"

	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type CharPrincess struct{}

func (p CharPrincess) Name() string { return "Princess" }
func (p CharPrincess) Number() int  { return Princess }
func (p CharPrincess) Text() string {
	return "You are eliminated if you discard the Princess"
}
func (p CharPrincess) Colour() string { return render.Yellow }

func (p CharPrincess) Play(g *Game, player int, args ...string) error {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played %s",
		g.RenderName(player),
		RenderCard(Princess),
	)))
	return nil
}
