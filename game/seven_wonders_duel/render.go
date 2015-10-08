package seven_wonders_duel

import (
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

const (
	CardTypeText      = "##"
	ProgressTokenText = `{{b}}{{c "green"}}@{{_c}}{{_b}}`
	ExtraTurnText     = `{{b}}{{c "blue"}}&{{_c}}{{_b}}`
	WonderText        = `{{b}}{{c "yellow"}}WOND{{_c}}{{_b}}`
)

var CardColours = map[int]string{
	CardTypeRaw:          render.Black,
	CardTypeManufactured: render.Gray,
	CardTypeCivilian:     render.Blue,
	CardTypeScientific:   render.Green,
	CardTypeCommercial:   render.Yellow,
	CardTypeMilitary:     render.Red,
	CardTypeGuild:        render.Magenta,
	CardTypeWonder:       render.Cyan,
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func RenderCoins(amount int) string {
	return render.Markup(fmt.Sprintf("$%d", amount), render.Yellow, true)
}

func RenderVP(amount int) string {
	return render.Markup(fmt.Sprintf("%dVP", amount), render.Green, true)
}

func RenderCardType(cardType int) string {
	return render.Markup(CardTypeText, CardColours[cardType], true)
}
