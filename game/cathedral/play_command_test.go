package cathedral

import (
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestPlay_Capture(t *testing.T) {
	g, err := parseGame(`
............G1......
..C1......G1G1G1G1..
C1C1C1..G1G1....G1..
..C1......G1R2..G1..
..C1......G2....G1..
..........G2G2G2....
....................
....................
..R1R1..............
....................
`)
	assert.NoError(t, err)
	assert.NoError(t, helper.Cmd(g, helper.Mick, "play 9 f9 down"))
	assertBoard(t, `
............G1G.G.G.
..C1......G1G1G1G1G.
C1C1C1..G1G1G.G.G1G.
..C1......G1G.G.G1G.
..C1......G2G.G.G1G.
..........G2G2G2G9G.
................G9G9
....................
..R1R1..............
....................
`, g.Board)
}
