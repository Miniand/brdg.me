package seven_wonders

import "encoding/gob"

func init() {
	gob.Register(CardMimicGuild{})
}

type CardMimicGuild struct {
	Card
}

func (c CardMimicGuild) VictoryPoints(player int, g *Game) int {
	vp := 0
	for _, dir := range DirNeighbours {
		for _, c := range g.Cards[g.NumFromPlayer(player, dir)] {
			if carder, ok := c.(Carder); !ok ||
				carder.GetCard().Kind != CardKindGuild {
				continue
			}
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

func (c CardMimicGuild) SuppString() string {
	return "Mimic a neighbouring guild card"
}
