package cathedral

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func parseGame(input string) (*Game, error) {
	board, err := parseBoard(input)
	if err != nil {
		return nil, err
	}
	g := &Game{}
	if err := g.Start(helper.Players[:2]); err != nil {
		return nil, err
	}
	g.Board = board
	for _, t := range g.Board {
		if t.Player != NoPlayer {
			player := t.Player
			if player == PlayerCathedral {
				player = 1
			}
			g.PlayedPieces[player][t.Type] = true
		}
	}
	return g, nil
}

func parseBoard(input string) (board Board, err error) {
	lines := strings.Split(input, "\n")
	board = Board{}
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
		for j := 0; j < 10; j++ {
			board[Loc{j, i}], err = parseTile(l[j*2 : (j+1)*2])
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
	case 'G':
		t.Player = 0
	case 'R':
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

func outputBoard(b Board) string {
	rowStrs := []string{}
	for _, row := range LocsByRow {
		rowStr := bytes.NewBuffer([]byte{})
		for _, l := range row {
			t := b[l]
			player := t.Player
			if player == NoPlayer {
				player = t.Owner
			}
			b := byte('.')
			switch player {
			case 0:
				b = 'G'
			case 1:
				b = 'R'
			case 2:
				b = 'C'
			}
			rowStr.WriteByte(b)
			b = '.'
			if t.Player != NoPlayer {
				b = '0' + byte(t.Type)
			}
			rowStr.WriteByte(b)
		}
		rowStrs = append(rowStrs, rowStr.String())
	}
	return strings.Join(rowStrs, "\n")
}

func assertBoardOutputsEqual(t *testing.T, expected, actual string) bool {
	trimmedExpected := strings.TrimSpace(expected)
	trimmedActual := strings.TrimSpace(actual)
	t.Logf(
		"Expected:\n%s\n\nActual:\n%s",
		trimmedExpected,
		trimmedActual,
	)
	return assert.Equal(t, trimmedExpected, trimmedActual)
}

func assertBoardsEqual(t *testing.T, expected, actual Board) bool {
	return assertBoardOutputsEqual(t, outputBoard(expected), outputBoard(actual))
}

func assertBoard(t *testing.T, expected string, actual Board) bool {
	return assertBoardOutputsEqual(t, expected, outputBoard(actual))
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
		Owner: 0,
	}, board[Loc{0, 0}])
	assert.Equal(t, Tile{
		PlayerType: PlayerType{
			Player: 0,
			Type:   3,
		},
		Owner: NoPlayer,
	}, board[Loc{1, 0}])
}

func TestOutputBoard(t *testing.T) {
	b1, err := parseBoard(`
G.G3................
G.G3................
G.G3....C1..........
G3G3..C1C1C1........
........C1..........
........C1..........
....................
..........R5R5......
............R5......
....................`)
	assert.NoError(t, err)
	b2, err := parseBoard(outputBoard(b1))
	assert.NoError(t, err)
	assertBoardsEqual(t, b1, b2)
}
