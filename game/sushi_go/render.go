package sushi_go

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/render"
)

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	buf := bytes.Buffer{}
	buf.WriteString("{{b}}Hand:{{_b}}\n")
	for i, c := range g.Hands[pNum] {
		buf.WriteString(fmt.Sprintf(
			"\n{{c \"gray\"}}(%d){{_c}} %s",
			i+1,
			RenderCard(c),
		))
	}
	return buf.String(), nil
}

func RenderCard(c int) string {
	return render.Markup(CardStrings[c], CardColours[c], c != CardPlayed)
}
