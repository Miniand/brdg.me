package seven_wonders_duel

import "github.com/Miniand/brdg.me/render"

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
