package seven_wonders

import (
	"bytes"
	"fmt"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/render"
)

var ResourceColours = map[int]string{
	GoodCoin:    render.Yellow,
	GoodWood:    render.Green,
	GoodStone:   render.Gray,
	GoodOre:     render.Black,
	GoodClay:    render.Red,
	GoodPapyrus: render.Cyan,
	GoodTextile: render.Magenta,
	GoodGlass:   render.Blue,

	CardKindRaw:          render.Black,
	CardKindManufactured: render.Gray,
	CardKindCivilian:     render.Blue,
	CardKindScientific:   render.Green,
	CardKindCommercial:   render.Yellow,
	CardKindMilitary:     render.Red,
	CardKindGuild:        render.Magenta,

	TokenVictory: render.Red,
	TokenDefeat:  render.Gray,

	WonderStage: render.Yellow,
}

const (
	CardSymbol = "##"
)

var ResourceSymbols = map[int]string{
	GoodCoin:    "Co",
	GoodWood:    "Wo",
	GoodStone:   "St",
	GoodOre:     "Or",
	GoodClay:    "Cl",
	GoodPapyrus: "Pa",
	GoodTextile: "Te",
	GoodGlass:   "Gl",

	CardKindRaw:          CardSymbol,
	CardKindManufactured: CardSymbol,
	CardKindCivilian:     CardSymbol,
	CardKindScientific:   CardSymbol,
	CardKindCommercial:   CardSymbol,
	CardKindMilitary:     CardSymbol,
	CardKindGuild:        CardSymbol,

	TokenVictory: "V",
	TokenDefeat:  "X",

	WonderStage: "▲",
}

func String(in interface{}) string {
	if str, ok := in.(string); ok {
		return str
	}
	if str, ok := in.(fmt.Stringer); ok {
		return str.String()
	}
	return fmt.Sprintf("%v", in)
}

func RenderResource(str interface{}, resource int, bold bool) string {
	return render.Markup(String(str), ResourceColours[resource], bold)
}

func RenderResourceSymbol(resource int) string {
	return RenderResource(ResourceSymbols[resource], resource, true)
}

func RenderCardName(c Carder) string {
	ca := c.GetCard()
	return RenderResource(
		fmt.Sprintf("%s %s", ResourceSymbols[ca.Kind], ca.Name),
		ca.Kind,
		true,
	)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	output := bytes.NewBuffer([]byte{})
	for _, d := range []card.Deck{
		DeckAge1(),
		DeckAge2(),
		DeckAge3(),
		DeckGuild(),
	} {
		for _, c := range d {
			output.WriteString(fmt.Sprintf("%s\n", RenderCardName(c.(Carder))))
		}
	}
	return output.String(), nil
}
