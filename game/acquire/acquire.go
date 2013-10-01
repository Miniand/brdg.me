package acquire

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/render"
	"strings"
)

const (
	TILE_EMPTY = iota
	TILE_PLACED
	TILE_CORP_WORLDWIDE
	TILE_CORP_SACKSON
	TILE_CORP_FESTIVAL
	TILE_CORP_IMPERIAL
	TILE_CORP_AMERICAN
	TILE_CORP_CONTINENTAL
	TILE_CORP_TOWER
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

const (
	INIT_SHARES      = 25
	INIT_CASH        = 6000
	INIT_TILES       = 6
	START_VALUE_LOW  = 200
	START_VALUE_MED  = 300
	START_VALUE_HIGH = 400
)

var CorpColours = map[int]string{
	TILE_CORP_WORLDWIDE:   "magenta",
	TILE_CORP_SACKSON:     "cyan",
	TILE_CORP_FESTIVAL:    "green",
	TILE_CORP_IMPERIAL:    "yellow",
	TILE_CORP_AMERICAN:    "blue",
	TILE_CORP_CONTINENTAL: "red",
	TILE_CORP_TOWER:       "black",
}

var CorpNames = map[int]string{
	TILE_CORP_WORLDWIDE:   "Worldwide",
	TILE_CORP_SACKSON:     "Sackson",
	TILE_CORP_FESTIVAL:    "Festival",
	TILE_CORP_IMPERIAL:    "Imperial",
	TILE_CORP_AMERICAN:    "American",
	TILE_CORP_CONTINENTAL: "Continental",
	TILE_CORP_TOWER:       "Tower",
}

var CorpShortNames = map[int]string{
	TILE_CORP_WORLDWIDE:   "WO",
	TILE_CORP_SACKSON:     "SA",
	TILE_CORP_FESTIVAL:    "FE",
	TILE_CORP_IMPERIAL:    "IM",
	TILE_CORP_AMERICAN:    "AM",
	TILE_CORP_CONTINENTAL: "CO",
	TILE_CORP_TOWER:       "TO",
}

var CorpStartValues = map[int]int{
	TILE_CORP_WORLDWIDE:   START_VALUE_LOW,
	TILE_CORP_SACKSON:     START_VALUE_LOW,
	TILE_CORP_FESTIVAL:    START_VALUE_MED,
	TILE_CORP_IMPERIAL:    START_VALUE_MED,
	TILE_CORP_AMERICAN:    START_VALUE_MED,
	TILE_CORP_CONTINENTAL: START_VALUE_HIGH,
	TILE_CORP_TOWER:       START_VALUE_HIGH,
}

type Game struct {
	Players      []string
	Board        map[int]map[int]int
	PlayerCash   map[int]int
	PlayerShares map[int]map[int]int
	PlayerTiles  map[int]card.Deck
	BankShares   map[int]int
	BankTiles    card.Deck
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

func RegisterGobTypes() {
	gob.Register(card.SuitRankCard{})
}

func (g *Game) Encode() ([]byte, error) {
	RegisterGobTypes()
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(g)
	return buf.Bytes(), err
}

func (g *Game) Decode(data []byte) error {
	RegisterGobTypes()
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(g)
}

func (g *Game) Start(players []string) error {
	g.Players = players
	// Initialise board
	g.Board = map[int]map[int]int{}
	for _, r := range Rows() {
		g.Board[r] = map[int]int{}
	}
	// Initialise player supplies
	g.PlayerCash = map[int]int{}
	g.PlayerShares = map[int]map[int]int{}
	g.BankTiles = Tiles().Shuffle()
	g.PlayerTiles = map[int]card.Deck{}
	for p, _ := range g.Players {
		g.PlayerCash[p] = INIT_CASH
		g.PlayerShares[p] = map[int]int{}
		g.PlayerTiles[p], g.BankTiles = g.BankTiles.PopN(INIT_TILES)
	}
	// Initialise shares
	g.BankShares = map[int]int{}
	for _, c := range Corps() {
		g.BankShares[c] = INIT_SHARES
	}
	// Testing values
	g.Board[BOARD_ROW_H][BOARD_COL_11] = TILE_PLACED
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
		output = fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, TileText(row, col))
	case TILE_PLACED:
		output = `{{b}}{{c "gray"}}XX{{_c}}{{_b}}`
	default:
		output = fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`, CorpColours[t],
			CorpShortNames[t])
	}
	return
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNumber(player)
	if err != nil {
		return "", err
	}
	output := bytes.NewBufferString("{{b}}Board:{{_b}}\n\n")
	// Board
	cells := [][]string{}
	for _, r := range Rows() {
		row := []string{}
		for _, c := range Cols() {
			cellOutput := g.RenderTile(r, c)
			// We embolden the tile if the player has it in their hand
			if _, n := g.PlayerTiles[pNum].Remove(card.SuitRankCard{
				Suit: r,
				Rank: c,
			}, 1); n > 0 {
				cellOutput = fmt.Sprintf("{{b}}%s{{_b}}", cellOutput)
			}
			row = append(row, cellOutput)
		}
		cells = append(cells, row)
	}
	boardOutput, err := render.Table(cells, 0, 1)
	if err != nil {
		return "", err
	}
	output.WriteString(boardOutput)
	// Hand
	handTiles := []string{}
	for _, t := range g.PlayerTiles[pNum].Sort() {
		tCard := t.(card.SuitRankCard)
		handTiles = append(handTiles, TileText(tCard.Suit, tCard.Rank))
	}
	output.WriteString(fmt.Sprintf(
		"\n\n{{b}}Your tiles: {{c \"gray\"}}%s{{_c}}{{_b}}\n",
		strings.Join(handTiles, " ")))
	output.WriteString(fmt.Sprintf(
		"{{b}}Your cash:  $%d{{_b}}", g.PlayerCash[pNum]))
	// Corp table
	cells = [][]string{
		[]string{
			"{{b}}Corporation{{_b}}",
			"{{b}}Size{{_b}}",
			"{{b}}Value{{_b}}",
			"{{b}}You own{{_b}}",
			"{{b}}Remaining{{_b}}",
		},
	}
	for _, c := range Corps() {
		cells = append(cells, []string{
			fmt.Sprintf(`{{b}}{{c "%s"}}%s (%s){{_c}}{{_b}}`, CorpColours[c],
				CorpNames[c], CorpShortNames[c]),
			fmt.Sprintf("%d", g.CorpSize(c)),
			fmt.Sprintf("$%d", g.CorpValue(c)),
			fmt.Sprintf("%d shares", g.PlayerShares[pNum][c]),
			fmt.Sprintf("%d shares", g.BankShares[c]),
		})
	}
	corpOutput, err := render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	output.WriteString("\n\n")
	output.WriteString(corpOutput)
	return output.String(), nil
}

func (g *Game) PlayerNumber(player string) (int, error) {
	for pNum, p := range g.Players {
		if p == player {
			return pNum, nil
		}
	}
	return 0, errors.New("Could not find player")
}

func (g *Game) CorpSize(corp int) (n int) {
	for _, r := range Rows() {
		for _, c := range Cols() {
			if g.Board[r][c] == corp {
				n += 1
			}
		}
	}
	return
}

func (g *Game) CorpValue(corp int) int {
	return CorpValue(g.CorpSize(corp), corp)
}

func CorpValue(size, corp int) int {
	if size <= 0 {
		return 0
	}
	multiplier := size - 2
	if size >= 41 {
		multiplier = 8
	} else if size >= 31 {
		multiplier = 7
	} else if size >= 21 {
		multiplier = 6
	} else if size >= 11 {
		multiplier = 5
	} else if size >= 6 {
		multiplier = 4
	}
	return CorpStartValues[corp] + multiplier*100
}

func CorpMajorityShareholderBonus(size, corp int) int {
	return CorpMinorityShareholderBonus(size, corp) * 2
}

func CorpMinorityShareholderBonus(size, corp int) int {
	return CorpValue(size, corp) * 5
}

func TileText(row, col int) string {
	return fmt.Sprintf("%d%c", 1+col, 'A'+row)
}

func Rows() []int {
	rows := []int{}
	for r := BOARD_ROW_A; r <= BOARD_ROW_I; r++ {
		rows = append(rows, r)
	}
	return rows
}

func Cols() []int {
	cols := []int{}
	for c := BOARD_COL_1; c <= BOARD_COL_12; c++ {
		cols = append(cols, c)
	}
	return cols
}

func Corps() []int {
	corps := []int{}
	for c := TILE_CORP_WORLDWIDE; c <= TILE_CORP_TOWER; c++ {
		corps = append(corps, c)
	}
	return corps
}

func Tiles() card.Deck {
	d := card.Deck{}
	for _, r := range Rows() {
		for _, c := range Cols() {
			d = d.Push(card.SuitRankCard{
				Suit: r,
				Rank: c,
			})
		}
	}
	return d
}
