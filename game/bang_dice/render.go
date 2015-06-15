package bang_dice

import "github.com/Miniand/brdg.me/render"

func (g *Game) RenderForPlayer(string) (string, error) {
	cells := [][]interface{}{}
	return render.Table(cells, 0, 2), nil
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
