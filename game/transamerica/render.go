package transamerica

import (
	"bytes"
	"fmt"
	"strings"
)

func (b *Board) RenderForPlayer(player int) string {
	var minX, minY, maxX, maxY byte
	first := true
	for l := range b.Nodes {
		if first || l.X < minX {
			minX = l.X
		}
		if first || l.Y < minY {
			minY = l.Y
		}
		if first || l.X > maxX {
			maxX = l.X
		}
		if first || l.Y > maxY {
			maxY = l.X
		}
		first = false
	}
	lines := []string{}
	for y := minY; y <= maxY; y++ {
		line := &bytes.Buffer{}
		if y%2 == 0 {
			line.WriteString("  ")
		}
		for x := minX; x <= maxX; x++ {
			loc := Loc{x, y}
			if n, ok := b.Nodes[loc]; ok {
				if n.City != "" {
					line.WriteString(fmt.Sprintf(
						`{{bg "%s"}}%s{{_bg}}`,
						n.City,
						loc,
					))
				} else {
					line.WriteString(loc.String())
				}
			} else {
				line.WriteString("  ")
			}
			line.WriteString("  ")
		}
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n\n")
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	b := &Board{
		Nodes: America,
	}
	return b.RenderForPlayer(0), nil
}
