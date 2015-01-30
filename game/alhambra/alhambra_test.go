package alhambra

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

var parseGridRegexp = regexp.MustCompile("(?ms)^\n?(.*)\n?$")

func parseGrid(t *testing.T, input string) Grid {
	g := Grid{}
	matches := parseGridRegexp.FindStringSubmatch(input)
	if matches == nil {
		t.Fatal("Unable to match input")
	}
	lines := strings.Split(matches[1], "\n")

	for y := 1; y < len(lines); y += 2 {
		for x := 1; x < len(lines[y]); x += 2 {
			if lines[y][x] == ' ' {
				continue
			}
			tileType, err := helper.MatchStringInStringMap(
				string(lines[y][x]),
				TileAbbrs,
			)
			if err != nil {
				t.Fatalf("Could not understand tile %c", lines[y][x])
			}
			if tileType == TileTypeEmpty {
				continue
			}
			t := Tile{
				Type:  tileType,
				Walls: map[int]bool{},
			}
			if isWall(charAt(x, y-1, lines)) {
				t.Walls[DirUp] = true
			}
			if isWall(charAt(x, y+1, lines)) {
				t.Walls[DirDown] = true
			}
			if isWall(charAt(x-1, y, lines)) {
				t.Walls[DirLeft] = true
			}
			if isWall(charAt(x+1, y, lines)) {
				t.Walls[DirRight] = true
			}
			g[Vect{(x - 1) / 2, (y - 1) / 2}] = t
		}
	}

	return g
}

func charAt(x, y int, lines []string) byte {
	if y < 0 || y >= len(lines) || x < 0 || x >= len(lines[y]) {
		return ' '
	}
	return lines[y][x]
}

func isWall(char byte) bool {
	return char == '|' || char == '-' || char == '+'
}

func TestParseGrid(t *testing.T) {
	assert.Equal(t, Grid{}, parseGrid(t, ""))

	assert.Equal(t, Grid{
		Vect{0, 0}: Tile{
			Type: TileTypeFountain,
			Walls: map[int]bool{
				DirUp:   true,
				DirDown: true,
				DirLeft: true,
			},
		},
		Vect{0, 1}: Tile{
			Type: TileTypeSeraglio,
			Walls: map[int]bool{
				DirUp: true,
			},
		},
	}, parseGrid(t, `
+-+
|F 
+-+
 S
`))
}
