package king_of_tokyo

import (
	"bytes"
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

const (
	VPSymbol     = `{{b}}{{c "blue"}}★{{_c}}{{_b}}`
	EnergySymbol = `{{b}}{{c "green"}}⚡{{_c}}{{_b}}`
	HealthSymbol = `{{b}}{{c "red"}}♥{{_c}}{{_b}}`
	AttackSymbol = `{{b}}{{c "yellow"}}☀{{_c}}{{_b}}`
)

var LocationStrings = map[int]string{
	LocationOutside:   render.Markup("outside Tokyo", "gray", false),
	LocationTokyoCity: render.Markup("Tokyo City", "yellow", true),
	LocationTokyoBay:  render.Markup("Tokyo Bay", "yellow", true),
}

func RenderVP(num int) string {
	return fmt.Sprintf("{{b}}%d{{_b}}%s", num, VPSymbol)
}

func RenderVPChange(num int) string {
	return fmt.Sprintf("{{b}}%+d{{_b}}%s", num, VPSymbol)
}

func RenderEnergy(num int) string {
	return fmt.Sprintf("{{b}}%d{{_b}}%s", num, EnergySymbol)
}

func RenderEnergyChange(num int) string {
	return fmt.Sprintf("{{b}}%+d{{_b}}%s", num, EnergySymbol)
}

func RenderHealth(num int) string {
	return fmt.Sprintf("{{b}}%d{{_b}}%s", num, HealthSymbol)
}

func RenderCardKind(kind int) string {
	col := "blue"
	text := "keep"
	if kind == CardKindDiscard {
		col = "red"
		text = "disc"
	}
	return render.Markup(text, col, true)
}

func RenderDie(kind int) string {
	var face string
	switch kind {
	case DieEnergy:
		face = EnergySymbol
	case DieAttack:
		face = AttackSymbol
	case DieHeal:
		face = HealthSymbol
	default:
		face = render.Markup(kind+1, "magenta", false)
	}
	return fmt.Sprintf("{{b}}[%s]{{_b}}", face)
}

func RenderCardTable(cards []CardBase) string {
	cells := [][]interface{}{}
	for _, c := range cards {
		cells = append(cells, [][]interface{}{
			{
				RenderEnergy(c.Cost()),
				fmt.Sprintf("%s (%s)", render.Bold(c.Name()), RenderCardKind(c.Kind())),
			},
			{
				"",
				render.Unbounded(render.Markup(c.Description(), "gray", false)),
			},
		}...)
	}
	return render.Table(cells, 0, 2)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	// Current roll
	diceStrs := []interface{}{render.Bold("Current roll:")}
	diceNums := []interface{}{""}
	for i, d := range g.CurrentRoll {
		diceStrs = append(diceStrs, RenderDie(d))
		diceNums = append(diceNums, render.Centred(render.Colour(i+1, "gray")))
	}
	cells := [][]interface{}{
		diceStrs,
		diceNums,
	}
	buf.WriteString(render.Table(cells, 0, 2))
	buf.WriteString("\n\n")
	// Player table
	cells = [][]interface{}{}
	for p, _ := range g.Players {
		cells = append(cells, []interface{}{
			g.RenderName(p),
			render.Centred(RenderHealth(g.Boards[p].Health)),
			render.Centred(RenderEnergy(g.Boards[p].Energy)),
			render.Centred(RenderVP(g.Boards[p].VP)),
			LocationStrings[g.PlayerLocation(p)],
		})
	}
	buf.WriteString(render.Table(cells, 0, 2))
	buf.WriteString("\n\n")
	// Shop
	buf.WriteString("{{b}}Available cards:{{_b}}\n")
	buf.WriteString(RenderCardTable(g.Buyable))
	// Player boards
	for p, _ := range g.Players {
		if len(g.Boards[p].Cards) > 0 {
			buf.WriteString(fmt.Sprintf(
				"\n\n{{b}}%s cards:{{_b}}\n",
				g.RenderName(p),
			))
			buf.WriteString(RenderCardTable(g.Boards[p].Cards))
		}
	}
	return buf.String(), nil
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
