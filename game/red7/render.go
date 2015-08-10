package red7

import (
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func RenderCard(card int) string {
	suit, rank := CardValues(card)
	return render.Markup(fmt.Sprintf(
		"%s%d",
		SuitAbbr[suit],
		RankVal[rank],
	), SuitCol[suit], true)
}
