package seven_wonders

import "strings"

type CardDrawDiscard struct {
	Card
	VP int
}

func (c CardDrawDiscard) VictoryTokens(player int, g *Game) int {
	return c.VP
}

func (c CardDrawDiscard) SuppString() string {
	parts := []string{"Build a discarded card for free"}
	if c.VP != 0 {
		parts = append(parts, RenderVP(c.VP))
	}
	return strings.Join(parts, " and ")
}

func (c CardDrawDiscard) HandlePostActionExecute(player int, g *Game) {
	if len(g.Discard) > 0 {
		g.ToResolve = append(g.ToResolve, ResolveDrawDiscard{player})
	}
}
