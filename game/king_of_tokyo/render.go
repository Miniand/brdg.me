package king_of_tokyo

import (
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

const (
	VPSymbol     = `{{b}}{{c "blue"}}★{{_c}}{{_b}}`
	EnergySymbol = `{{b}}{{c "green"}}⚡{{_c}}{{_b}}`
	HealthSymbol = `{{b}}{{c "red"}}♥{{_c}}{{_b}}`
	AttackSymbol = `{{b}}{{c "yellow"}}☀{{_c}}{{_b}}`
)

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

func (g *Game) RenderForPlayer(player string) (string, error) {
	cells := [][]interface{}{}
	for _, c := range Deck() {
		cells = append(cells, [][]interface{}{
			{
				RenderEnergy(c.Cost()),
				fmt.Sprintf("%s (%s)", render.Bold(c.Name()), RenderCardKind(c.Kind())),
			},
			{
				render.Unbounded(render.Markup(c.Description(), "gray", false)),
			},
			{},
		}...)
	}
	return render.Table(cells, 0, 2), nil
}
