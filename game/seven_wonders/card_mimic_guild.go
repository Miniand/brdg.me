package seven_wonders

type CardMimicGuild struct {
	Card
}

func (c CardMimicGuild) VictoryPoints(player int, g *Game) int {
	vp := 0
	for _, dir := range DirNeighbours {
		for _, c := range g.Cards[g.NumFromPlayer(player, dir)] {
			if vper, ok := c.(VictoryPointer); ok {
				cVP := vper.VictoryPoints(player, g)
				if cVP > vp {
					vp = cVP
				}
			}
		}
	}
	return vp
}
