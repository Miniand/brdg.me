package king_of_tokyo

import (
	"fmt"
	"strconv"
)

const (
	VPSymbol     = `{{b}}{{c "blue"}}★{{_c}}{{_b}}`
	EnergySymbol = `{{b}}{{c "green"}}e{{_c}}{{_b}}`
	HealthSymbol = `{{b}}{{c "red"}}h{{_c}}{{_b}}`
	AttackSymbol = `{{b}}{{c "yellow"}}a{{_c}}{{_b}}`
)

func RenderVP(num int) string {
	return fmt.Sprintf("{{b}}%d{{_b}} %s", num, VPSymbol)
}

func RenderVPChange(num int) string {
	return fmt.Sprintf("{{b}}%+d{{_b}} %s", num, VPSymbol)
}

func RenderEnergy(num int) string {
	return fmt.Sprintf("{{b}}%d{{_b}} %s", num, EnergySymbol)
}

func RenderEnergyChange(num int) string {
	return fmt.Sprintf("{{b}}%+d{{_b}} %s", num, EnergySymbol)
}

func RenderHealth(num int) string {
	return fmt.Sprintf("{{b}}%d{{_b}} %s", num, HealthSymbol)
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
		face = strconv.Itoa(kind + 1)
	}
	return fmt.Sprintf("[%s]", face)
}
