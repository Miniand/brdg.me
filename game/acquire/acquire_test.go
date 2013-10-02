package acquire

import (
	"github.com/Miniand/brdg.me/game/card"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func (g *Game) parseGameBoard(b string, t *testing.T) {
	g.BankTiles = card.Deck{}
	rows := regexp.MustCompile(`\n`).Split(strings.TrimSpace(b), -1)
	if len(rows) != 9 {
		t.Fatal("Must be 9 rows")
	}
	colReg := regexp.MustCompile(`\s+`)
	for rowN, row := range rows {
		cols := colReg.Split(row, -1)
		if len(cols) != 12 {
			t.Fatal("Must be 12 cols")
		}
		for colN, cell := range cols {
			val, err := strconv.Atoi(cell)
			if err != nil {
				t.Fatal(err)
			}
			if val < TILE_EMPTY || val > TILE_CORP_TOWER {
				t.Fatal("Invalid number")
			}
			g.Board[rowN][colN] = val
			if val == TILE_EMPTY {
				g.BankTiles = g.BankTiles.Push(Tile{
					Row:    rowN,
					Column: colN,
				})
			}
		}
	}
	for playerN, _ := range g.Players {
		g.PlayerTiles[playerN], g.BankTiles = g.BankTiles.PopN(INIT_TILES)
	}
}

func TestStart(t *testing.T) {
	g := Game{}
	if err := g.Start([]string{"Mick", "Steve", "BJ"}); err != nil {
		t.Fatal(err)
	}
}

func checkCorpValues(corp int, expected map[int]int, t *testing.T) {
	for size, expectedValue := range expected {
		actual := CorpValue(size, corp)
		if actual != expectedValue {
			t.Fatal("Corp", corp, "size", size, "expected", expectedValue,
				"got", actual)
		}
	}
}

func TestCorpValue(t *testing.T) {
	low := map[int]int{
		2:  200,
		3:  300,
		4:  400,
		5:  500,
		6:  600,
		10: 600,
		11: 700,
		20: 700,
		21: 800,
		30: 800,
		31: 900,
		40: 900,
		41: 1000,
	}
	med := map[int]int{
		2:  300,
		3:  400,
		4:  500,
		5:  600,
		6:  700,
		10: 700,
		11: 800,
		20: 800,
		21: 900,
		30: 900,
		31: 1000,
		40: 1000,
		41: 1100,
	}
	high := map[int]int{
		2:  400,
		3:  500,
		4:  600,
		5:  700,
		6:  800,
		10: 800,
		11: 900,
		20: 900,
		21: 1000,
		30: 1000,
		31: 1100,
		40: 1100,
		41: 1200,
	}
	checkCorpValues(TILE_CORP_WORLDWIDE, low, t)
	checkCorpValues(TILE_CORP_SACKSON, low, t)
	checkCorpValues(TILE_CORP_FESTIVAL, med, t)
	checkCorpValues(TILE_CORP_IMPERIAL, med, t)
	checkCorpValues(TILE_CORP_AMERICAN, med, t)
	checkCorpValues(TILE_CORP_CONTINENTAL, high, t)
	checkCorpValues(TILE_CORP_TOWER, high, t)
}

func checkCorp1stBonuses(corp int, expected map[int]int,
	t *testing.T) {
	for size, expectedValue := range expected {
		actual := Corp1stBonus(size, corp)
		if actual != expectedValue {
			t.Fatal("Corp", corp, "size", size, "expected", expectedValue,
				"got", actual)
		}
	}
}

func TestCorp1stBonuses(t *testing.T) {
	low := map[int]int{
		2:  2000,
		3:  3000,
		4:  4000,
		5:  5000,
		6:  6000,
		10: 6000,
		11: 7000,
		20: 7000,
		21: 8000,
		30: 8000,
		31: 9000,
		40: 9000,
		41: 10000,
	}
	med := map[int]int{
		2:  3000,
		3:  4000,
		4:  5000,
		5:  6000,
		6:  7000,
		10: 7000,
		11: 8000,
		20: 8000,
		21: 9000,
		30: 9000,
		31: 10000,
		40: 10000,
		41: 11000,
	}
	high := map[int]int{
		2:  4000,
		3:  5000,
		4:  6000,
		5:  7000,
		6:  8000,
		10: 8000,
		11: 9000,
		20: 9000,
		21: 10000,
		30: 10000,
		31: 11000,
		40: 11000,
		41: 12000,
	}
	checkCorp1stBonuses(TILE_CORP_WORLDWIDE, low, t)
	checkCorp1stBonuses(TILE_CORP_SACKSON, low, t)
	checkCorp1stBonuses(TILE_CORP_FESTIVAL, med, t)
	checkCorp1stBonuses(TILE_CORP_IMPERIAL, med, t)
	checkCorp1stBonuses(TILE_CORP_AMERICAN, med, t)
	checkCorp1stBonuses(TILE_CORP_CONTINENTAL, high, t)
	checkCorp1stBonuses(TILE_CORP_TOWER, high, t)
}

func checkCorp2ndBonuses(corp int, expected map[int]int,
	t *testing.T) {
	for size, expectedValue := range expected {
		actual := Corp2ndBonus(size, corp)
		if actual != expectedValue {
			t.Fatal("Corp", corp, "size", size, "expected", expectedValue,
				"got", actual)
		}
	}
}

func TestCorp2ndBonuses(t *testing.T) {
	low := map[int]int{
		2:  1000,
		3:  1500,
		4:  2000,
		5:  2500,
		6:  3000,
		10: 3000,
		11: 3500,
		20: 3500,
		21: 4000,
		30: 4000,
		31: 4500,
		40: 4500,
		41: 5000,
	}
	med := map[int]int{
		2:  1500,
		3:  2000,
		4:  2500,
		5:  3000,
		6:  3500,
		10: 3500,
		11: 4000,
		20: 4000,
		21: 4500,
		30: 4500,
		31: 5000,
		40: 5000,
		41: 5500,
	}
	high := map[int]int{
		2:  2000,
		3:  2500,
		4:  3000,
		5:  3500,
		6:  4000,
		10: 4000,
		11: 4500,
		20: 4500,
		21: 5000,
		30: 5000,
		31: 5500,
		40: 5500,
		41: 6000,
	}
	checkCorp2ndBonuses(TILE_CORP_WORLDWIDE, low, t)
	checkCorp2ndBonuses(TILE_CORP_SACKSON, low, t)
	checkCorp2ndBonuses(TILE_CORP_FESTIVAL, med, t)
	checkCorp2ndBonuses(TILE_CORP_IMPERIAL, med, t)
	checkCorp2ndBonuses(TILE_CORP_AMERICAN, med, t)
	checkCorp2ndBonuses(TILE_CORP_CONTINENTAL, high, t)
	checkCorp2ndBonuses(TILE_CORP_TOWER, high, t)
}

func checkParseTileText(text string, expectedRow, expectedCol int,
	t *testing.T) {
	tile, err := ParseTileText(text)
	if err != nil {
		t.Fatal(err)
	}
	if tile.Row != expectedRow {
		t.Fatal("Expected row", expectedRow, "got", tile.Row)
	}
	if tile.Column != expectedCol {
		t.Fatal("Expected col", expectedCol, "got", tile.Column)
	}
}

func TestParseTileText(t *testing.T) {
	checkParseTileText("1A", BOARD_ROW_A, BOARD_COL_1, t)
	checkParseTileText("12i", BOARD_ROW_I, BOARD_COL_12, t)
	checkParseTileText("6c", BOARD_ROW_C, BOARD_COL_6, t)
}

func TestAdjacentTiles(t *testing.T) {
	adj := AdjacentTiles(Tile{BOARD_ROW_A, BOARD_COL_1})
	if len(adj) != 2 {
		t.Fatal("Expected there to only be two")
	}
	if _, n := adj.Remove(Tile{BOARD_ROW_B, BOARD_COL_1}, -1); n != 1 {
		t.Fatal("Expected 1B to be adjacent")
	}
	if _, n := adj.Remove(Tile{BOARD_ROW_A, BOARD_COL_2}, -1); n != 1 {
		t.Fatal("Expected 2A to be adjacent")
	}
	adj = AdjacentTiles(Tile{BOARD_ROW_F, BOARD_COL_4})
	if len(adj) != 4 {
		t.Fatal("Expected there to only be four")
	}
	if _, n := adj.Remove(Tile{BOARD_ROW_F, BOARD_COL_3}, -1); n != 1 {
		t.Fatal("Expected 3F to be adjacent")
	}
	if _, n := adj.Remove(Tile{BOARD_ROW_F, BOARD_COL_5}, -1); n != 1 {
		t.Fatal("Expected 5F to be adjacent")
	}
	if _, n := adj.Remove(Tile{BOARD_ROW_E, BOARD_COL_4}, -1); n != 1 {
		t.Fatal("Expected 4E to be adjacent")
	}
	if _, n := adj.Remove(Tile{BOARD_ROW_G, BOARD_COL_4}, -1); n != 1 {
		t.Fatal("Expected 4G to be adjacent")
	}
}

func TestIsJoiningSafeCorps(t *testing.T) {
	g := Game{}
	if err := g.Start([]string{"Mick", "Steve", "BJ"}); err != nil {
		t.Fatal(err)
	}
	g.parseGameBoard(`
0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 4 4 0 2 2 0 0 0 0
0 0 0 4 4 0 2 2 0 5 0 0
0 0 0 4 4 0 2 2 0 5 0 0
0 0 0 4 4 0 2 2 0 5 0 0
0 0 0 4 4 0 2 2 0 0 0 0
0 0 0 4 4 0 2 2 0 0 0 0
0 0 0 0 0 3 0 0 0 0 0 0
0 0 0 0 3 3 3 0 0 0 0 0
`, t)
	if !g.IsJoiningSafeCorps(Tile{
		Row:    BOARD_ROW_B,
		Column: BOARD_COL_6,
	}) {
		t.Fatal("6B should be joining safe corps")
	}
	if !g.IsJoiningSafeCorps(Tile{
		Row:    BOARD_ROW_G,
		Column: BOARD_COL_6,
	}) {
		t.Fatal("6G should be joining safe corps")
	}
	if g.IsJoiningSafeCorps(Tile{
		Row:    BOARD_ROW_H,
		Column: BOARD_COL_5,
	}) {
		t.Fatal("5H shouldn't be joining safe corps")
	}
	if g.IsJoiningSafeCorps(Tile{
		Row:    BOARD_ROW_C,
		Column: BOARD_COL_9,
	}) {
		t.Fatal("9C shouldn't be joining safe corps")
	}
}

func checkPotentialMergers(expected, actual [][2]int, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fatal("Expected length", len(expected), "got", len(actual))
	}
	for _, e := range expected {
		found := false
		for _, a := range actual {
			if e == a {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Could not find %#v in %#v", e, actual)
		}
	}
}

func TestPotentialMergers(t *testing.T) {
	g := Game{}
	if err := g.Start([]string{"Mick", "Steve", "BJ"}); err != nil {
		t.Fatal(err)
	}
	g.parseGameBoard(`
0 0 0 7 0 0 0 0 0 0 0 0
0 0 0 7 0 0 0 0 0 0 0 0
0 6 6 0 0 0 0 0 0 5 0 0
0 0 0 4 4 0 2 2 0 5 0 0
0 0 0 4 4 0 2 2 0 5 0 0
0 0 0 4 4 0 2 2 0 0 0 0
0 0 0 4 4 0 2 2 0 0 0 0
0 0 0 0 0 3 0 0 0 0 0 0
0 0 0 0 3 3 3 0 0 0 0 0
`, t)
	checkPotentialMergers([][2]int{
		[2]int{2, 4},
		[2]int{4, 2},
	}, g.PotentialMergers(Tile{
		Row:    BOARD_ROW_G,
		Column: BOARD_COL_6,
	}), t)
	checkPotentialMergers([][2]int{
		[2]int{3, 4},
	}, g.PotentialMergers(Tile{
		Row:    BOARD_ROW_H,
		Column: BOARD_COL_5,
	}), t)
	checkPotentialMergers([][2]int{
		[2]int{6, 4},
		[2]int{7, 4},
	}, g.PotentialMergers(Tile{
		Row:    BOARD_ROW_C,
		Column: BOARD_COL_4,
	}), t)
}
