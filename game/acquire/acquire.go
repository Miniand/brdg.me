package acquire

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/render"
)

const (
	TILE_EMPTY = iota
	TILE_CORP_AMERICAN
	TILE_CORP_CONTINENTAL
	TILE_CORP_FESTIVAL
	TILE_CORP_IMPERIAL
	TILE_CORP_SACKSON
	TILE_CORP_TOWER
	TILE_CORP_WORLDWIDE
)

const (
	BOARD_ROW_A = iota
	BOARD_ROW_B
	BOARD_ROW_C
	BOARD_ROW_D
	BOARD_ROW_E
	BOARD_ROW_F
	BOARD_ROW_G
	BOARD_ROW_H
	BOARD_ROW_I
)

const (
	BOARD_COL_1 = iota
	BOARD_COL_2
	BOARD_COL_3
	BOARD_COL_4
	BOARD_COL_5
	BOARD_COL_6
	BOARD_COL_7
	BOARD_COL_8
	BOARD_COL_9
	BOARD_COL_10
	BOARD_COL_11
	BOARD_COL_12
)

var CorpColours = map[int]string{
	TILE_CORP_AMERICAN:    "blue",
	TILE_CORP_CONTINENTAL: "red",
	TILE_CORP_FESTIVAL:    "green",
	TILE_CORP_IMPERIAL:    "yellow",
	TILE_CORP_SACKSON:     "magenta",
	TILE_CORP_TOWER:       "cyan",
	TILE_CORP_WORLDWIDE:   "gray",
}

var CorpNames = map[int]string{
	TILE_CORP_AMERICAN:    "American",
	TILE_CORP_CONTINENTAL: "Continental",
	TILE_CORP_FESTIVAL:    "Festival",
	TILE_CORP_IMPERIAL:    "Imperial",
	TILE_CORP_SACKSON:     "Sackson",
	TILE_CORP_TOWER:       "Tower",
	TILE_CORP_WORLDWIDE:   "Worldwide",
}

var CorpShortNames = map[int]string{
	TILE_CORP_AMERICAN:    "AM",
	TILE_CORP_CONTINENTAL: "CO",
	TILE_CORP_FESTIVAL:    "FE",
	TILE_CORP_IMPERIAL:    "IM",
	TILE_CORP_SACKSON:     "SA",
	TILE_CORP_TOWER:       "TO",
	TILE_CORP_WORLDWIDE:   "WO",
}

type Game struct {
	Players []string
	Board   [BOARD_ROW_I + 1][BOARD_COL_12 + 1]int
}

func (g *Game) Name() string {
	return "Acquire"
}

func (g *Game) Identifier() string {
	return "acquire"
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Encode() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(g)
	return buf.Bytes(), err
}

func (g *Game) Decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(g)
}

func (g *Game) Start(players []string) error {
	g.Players = players
	g.Board[BOARD_ROW_B][BOARD_COL_6] = TILE_CORP_AMERICAN
	g.Board[BOARD_ROW_C][BOARD_COL_1] = TILE_CORP_CONTINENTAL
	g.Board[BOARD_ROW_A][BOARD_COL_4] = TILE_CORP_FESTIVAL
	g.Board[BOARD_ROW_D][BOARD_COL_7] = TILE_CORP_IMPERIAL
	g.Board[BOARD_ROW_F][BOARD_COL_3] = TILE_CORP_SACKSON
	g.Board[BOARD_ROW_I][BOARD_COL_12] = TILE_CORP_TOWER
	g.Board[BOARD_ROW_G][BOARD_COL_10] = TILE_CORP_WORLDWIDE
	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return g.Players
}

func (g *Game) RenderTile(row, col int) (output string) {
	t := g.Board[row][col]
	switch t {
	case TILE_EMPTY:
		output = TileText(row, col)
	default:
		output = fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`, CorpColours[t],
			CorpShortNames[t])
	}
	return
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	cells := [][]string{}
	for r := BOARD_ROW_A; r <= BOARD_ROW_I; r++ {
		row := []string{}
		for c := BOARD_COL_1; c <= BOARD_COL_12; c++ {
			row = append(row, g.RenderTile(r, c))
		}
		cells = append(cells, row)
	}
	for letter := 'A'; letter <= 'I'; letter++ {
		for number := 1; number <= 12; number++ {
		}
	}
	return render.Table(cells, 0, 1)
}

func TileText(row, col int) string {
	return fmt.Sprintf("%d%c", 1+col, 'A'+row)
}
