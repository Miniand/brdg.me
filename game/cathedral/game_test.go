package cathedral

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func parseBoard(input string) (board [10][10]Tile, err error) {
	lines := strings.Split(input, "\n")
	board = [10][10]Tile{}
	i := 0
	for _, l := range lines {
		if l == "" {
			continue
		}
		if i >= 10 {
			err = errors.New("out of range")
			return
		}
		if len(l) != 20 {
			err = errors.New("expected row to be 20 chars long")
		}
		board[i] = [10]Tile{}
		for j := 0; j < 10; j++ {
			board[i][j], err = parseTile(l[j*2 : (j+1)*2])
			if err != nil {
				return
			}
		}
		i++
	}
	if i != 10 {
		err = errors.New("there wasn't 10 rows")
	}
	return
}

func parseTile(input string) (t Tile, err error) {
	t = EmptyTile
	if len(input) != 2 {
		err = errors.New("input should be len 2")
		return
	}
	if input == ".." {
		t.Player = NoPlayer
		return
	}
	switch input[0] {
	case 'R':
		t.Player = 0
	case 'G':
		t.Player = 1
	case 'C':
		t.Player = PlayerCathedral
	default:
		err = errors.New("tile should start with '.', 'R', 'G' or 'C'")
		return
	}
	switch input[1] {
	case '.':
		t.Owner, t.Player = t.Player, NoPlayer
	default:
		t.Type, err = strconv.Atoi(string(input[1]))
	}
	return
}

func TestGame_Encode(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:2]))
	data, err := g.Encode()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestParseBoard(t *testing.T) {
	board, err := parseBoard(`
G.G3................
G.G3................
G.G3................
G3G3................
....................
....................
....................
....................
....................
....................`)
	assert.NoError(t, err)
	assert.Equal(t, Tile{
		PlayerType: PlayerType{
			Player: NoPlayer,
		},
		Owner: 1,
	}, board[0][0])
	assert.Equal(t, Tile{
		PlayerType: PlayerType{
			Player: 1,
			Type:   3,
		},
		Owner: NoPlayer,
	}, board[0][1])
}
