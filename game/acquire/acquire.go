package acquire

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/render"
	"regexp"
	"strconv"
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
	TURN_PHASE_PLAY_TILE = iota
	TURN_PHASE_PLACE_CORP
	TURN_PHASE_MERGER
	TURN_PHASE_BUY_TILES
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
	Players       []string
	CurrentPlayer int
	TurnPhase     int
	GameEnded     bool
	Board         map[int]map[int]int
	PlayerCash    map[int]int
	PlayerShares  map[int]map[int]int
	PlayerTiles   map[int]card.Deck
	BankShares    map[int]int
	BankTiles     card.Deck
}

func (g *Game) Name() string {
	return "Acquire"
}

func (g *Game) Identifier() string {
	return "acquire"
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
	}
}

func RegisterGobTypes() {
	gob.Register(Tile{})
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
	if len(players) < 3 || len(players) > 6 {
		return errors.New("Acquire is between 3 and 6 players")
	}
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
	return g.GameEnded
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) RenderTile(t Tile) (output string) {
	val := g.Board[t.Row][t.Column]
	switch val {
	case TILE_EMPTY:
		output = fmt.Sprintf(`{{c "gray"}}%s{{_c}}`, TileText(t))
	case TILE_PLACED:
		output = `{{b}}{{c "gray"}}XX{{_c}}{{_b}}`
	default:
		output = fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`, CorpColours[val],
			CorpShortNames[val])
	}
	return
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	output := bytes.NewBufferString("{{b}}Board:{{_b}}\n\n")
	// Board
	cells := [][]string{}
	for _, r := range Rows() {
		row := []string{}
		for _, c := range Cols() {
			cellOutput := g.RenderTile(Tile{r, c})
			// We embolden the tile if the player has it in their hand
			if _, n := g.PlayerTiles[pNum].Remove(Tile{
				Row:    r,
				Column: c,
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
	for _, tRaw := range g.PlayerTiles[pNum].Sort() {
		t := tRaw.(Tile)
		handTiles = append(handTiles, TileText(t))
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
			"{{b}}1st bonus{{_b}}",
			"{{b}}2nd bonus{{_b}}",
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
			fmt.Sprintf("$%d", g.Corp1stBonus(c)),
			fmt.Sprintf("$%d", g.Corp2ndBonus(c)),
		})
	}
	corpOutput, err := render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	output.WriteString("\n\n")
	output.WriteString(corpOutput)
	if g.IsFinished() {
		// Player table
		cells = [][]string{
			[]string{
				"{{b}}Player{{_b}}",
			},
		}
		for pNum, p := range g.Players {
			cells = append(cells, []string{
				fmt.Sprintf("{{b}}%s{{_b}}", render.PlayerName(pNum, p)),
			})
		}
		playerOutput, err := render.Table(cells, 0, 2)
		if err != nil {
			return "", err
		}
		output.WriteString("\n\n")
		output.WriteString(playerOutput)
	}
	return output.String(), nil
}

func (g *Game) PlayerNum(player string) (int, error) {
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

func (g *Game) Corp1stBonus(corp int) int {
	return Corp1stBonus(g.CorpSize(corp), corp)
}

func (g *Game) Corp2ndBonus(corp int) int {
	return Corp2ndBonus(g.CorpSize(corp), corp)
}

func (g *Game) PlayTile(playerNum int, t Tile) error {
	if g.IsFinished() || g.CurrentPlayer != playerNum ||
		g.TurnPhase != TURN_PHASE_PLAY_TILE {
		return errors.New("You are not allowed to play a tile at the moment")
	}
	// Check the tile is in the relevant player's hand
	newPlayerTiles, n := g.PlayerTiles[playerNum].Remove(t, 1)
	if n == 0 {
		return errors.New("You don't have that tile in your hand")
	}
	g.PlayerTiles[playerNum] = newPlayerTiles
	// Check if the tile is adjacent to others

	return nil
}

func IsValidLocation(t Tile) bool {
	return t.Row >= BOARD_ROW_A && t.Row <= BOARD_ROW_I &&
		t.Column >= BOARD_COL_1 && t.Column <= BOARD_COL_12
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

func Corp1stBonus(size, corp int) int {
	return Corp2ndBonus(size, corp) * 2
}

func Corp2ndBonus(size, corp int) int {
	return CorpValue(size, corp) * 5
}

func TileText(t Tile) string {
	return fmt.Sprintf("%d%c", 1+t.Column, 'A'+t.Row)
}

func ParseTileText(text string) (t Tile, err error) {
	matches := regexp.MustCompile(`\b(\d{1,2})([A-I])\b`).FindStringSubmatch(
		strings.ToUpper(text))
	if matches == nil {
		err = errors.New(
			"Invalid tile, it must be between 1A and 12H")
		return
	}
	t.Column, err = strconv.Atoi(matches[1])
	if err != nil {
		return
	}
	t.Column = t.Column - 1 // 0 index
	t.Row = int(matches[2][0] - 'A')
	if !IsValidLocation(t) {
		err = errors.New("Invalid tile, it must be between 1A and 12H")
	}
	return
}

func AdjacentTiles(t Tile) card.Deck {
	d := card.Deck{}
	// Left
	left := Tile{t.Row - 1, t.Column}
	if IsValidLocation(left) {
		d = d.Push(left)
	}
	// Right
	right := Tile{t.Row + 1, t.Column}
	if IsValidLocation(right) {
		d = d.Push(right)
	}
	// Up
	up := Tile{t.Row, t.Column - 1}
	if IsValidLocation(up) {
		d = d.Push(up)
	}
	// Down
	down := Tile{t.Row, t.Column + 1}
	if IsValidLocation(down) {
		d = d.Push(down)
	}
	return d
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
			d = d.Push(Tile{r, c})
		}
	}
	return d
}
