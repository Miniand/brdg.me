package acquire

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	TILE_DISCARDED = iota - 1
	TILE_EMPTY
	TILE_UNINCORPORATED
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
	TURN_PHASE_FOUND_CORP
	TURN_PHASE_MERGER_CHOOSE
	TURN_PHASE_MERGER
	TURN_PHASE_BUY_SHARES
)

const (
	INIT_SHARES        = 25
	INIT_CASH          = 6000
	INIT_TILES         = 6
	START_VALUE_LOW    = 200
	START_VALUE_MED    = 300
	START_VALUE_HIGH   = 400
	CORP_SAFE_SIZE     = 11
	CORP_END_GAME_SIZE = 41
	TILE_REGEXP        = `\b(1[012]|[1-9])([A-I])\b`
	MAX_BUY_PER_TURN   = 3
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
	Players             []string
	CurrentPlayer       int
	TurnPhase           int
	GameEnded           bool
	FinalTurn           bool
	Board               map[int]map[int]int
	PlayerCash          map[int]int
	PlayerShares        map[int]map[int]int
	PlayerTiles         map[int]card.Deck
	BankShares          map[int]int
	BankTiles           card.Deck
	PlayedTile          Tile
	MergerCurrentPlayer int
	MergerFromCorp      int
	MergerIntoCorp      int
	BoughtShares        int
	Log                 *log.Log
}

func (g *Game) Name() string {
	return "Acquire"
}

func (g *Game) Identifier() string {
	return "acquire"
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		EndCommand{},
		PlayCommand{},
		MergeCommand{},
		SellCommand{},
		TradeCommand{},
		KeepCommand{},
		BuyCommand{},
		DoneCommand{},
		FoundCommand{},
	}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
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
	if len(players) < 2 || len(players) > 6 {
		return errors.New("Acquire is between 2 and 6 players")
	}
	g.Players = players
	g.Log = log.New()
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
	// Draw tiles to the board for the number of players
	var rawT card.Card
	tileStrings := []string{}
	for _, _ = range g.Players {
		rawT, g.BankTiles = g.BankTiles.Pop()
		t := rawT.(Tile)
		g.Board[t.Row][t.Column] = TILE_UNINCORPORATED
		tileStrings = append(tileStrings, fmt.Sprintf(
			"{{b}}%s{{_b}}", TileText(t)))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Started the game with %s on the board",
		render.CommaList(tileStrings))))
	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.GameEnded
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	highestCash := 0
	highestPlayers := []string{}
	for pNum, p := range g.Players {
		if g.PlayerCash[pNum] > highestCash {
			highestCash = g.PlayerCash[pNum]
			highestPlayers = []string{}
		}
		if g.PlayerCash[pNum] == highestCash {
			highestPlayers = append(highestPlayers, p)
		}
	}
	return highestPlayers
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	if g.TurnPhase == TURN_PHASE_MERGER {
		return []string{g.Players[g.MergerCurrentPlayer]}
	}
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) RenderTile(t Tile) (output string) {
	val := g.Board[t.Row][t.Column]
	switch val {
	case TILE_DISCARDED:
		output = "  "
	case TILE_EMPTY:
		output = `{{c "gray"}}--{{_c}}`
	case TILE_UNINCORPORATED:
		output = `{{c "gray"}}##{{_c}}`
	default:
		output = fmt.Sprintf(`{{b}}%s{{_b}}`, RenderCorpShort(val))
	}
	return
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	output := bytes.NewBufferString("")
	output.WriteString("{{b}}Board:{{_b}}\n\n")
	// Board
	cells := [][]interface{}{}
	for _, r := range Rows() {
		row := []interface{}{}
		for _, c := range Cols() {
			cellOutput := g.RenderTile(Tile{r, c})
			// We embolden the tile if the player has it in their hand
			t := Tile{
				Row:    r,
				Column: c,
			}
			if _, n := g.PlayerTiles[pNum].Remove(t, 1); n > 0 {
				cellOutput = fmt.Sprintf(`{{c "gray"}}{{b}}%s{{_b}}{{_c}}`,
					TileText(t))
			}
			row = append(row, cellOutput)
		}
		cells = append(cells, row)
	}
	boardOutput := render.Table(cells, 0, 1)
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
		"{{b}}Your cash:  $%d{{_b}}\n", g.PlayerCash[pNum]))
	output.WriteString(fmt.Sprintf(
		"{{b}}Tiles left: %d{{_b}}", len(g.BankTiles)))
	// Corp table
	cells = [][]interface{}{
		[]interface{}{
			"{{b}}Corporation{{_b}}",
			"{{b}}Size{{_b}}",
			"{{b}}Value{{_b}}",
			"{{b}}Shares{{_b}}",
			"{{b}}Major{{_b}}",
			"{{b}}Minor{{_b}}",
		},
	}
	for _, c := range Corps() {
		cells = append(cells, []interface{}{
			fmt.Sprintf(`{{b}}%s{{_b}}`, RenderCorpWithShort(c)),
			fmt.Sprintf("%d", g.CorpSize(c)),
			fmt.Sprintf("$%d", g.CorpValue(c)),
			fmt.Sprintf("%d left", g.BankShares[c]),
			fmt.Sprintf("$%d", g.Corp1stBonus(c)),
			fmt.Sprintf("$%d", g.Corp2ndBonus(c)),
		})
	}
	corpOutput := render.Table(cells, 0, 2)
	output.WriteString("\n\n")
	output.WriteString(corpOutput)
	// Player table
	playerHeadings := []interface{}{
		"{{b}}Player{{_b}}",
		"{{b}}Cash{{_b}}",
	}
	for _, corp := range Corps() {
		playerHeadings = append(playerHeadings, fmt.Sprintf(
			"{{b}}%s{{_b}}", RenderCorpShort(corp)))
	}
	cells = [][]interface{}{
		playerHeadings,
	}
	for pNum, p := range g.Players {
		row := []interface{}{
			fmt.Sprintf("{{b}}%s{{_b}}", render.PlayerName(pNum, p)),
			fmt.Sprintf("$%d", g.PlayerCash[pNum]),
		}
		for _, corp := range Corps() {
			row = append(row, fmt.Sprintf("%d", g.PlayerShares[pNum][corp]))
		}
		cells = append(cells, row)
	}
	playerOutput := render.Table(cells, 0, 2)
	output.WriteString("\n\n")
	output.WriteString(playerOutput)
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
	// VALIDATE
	if g.IsFinished() || g.CurrentPlayer != playerNum ||
		g.TurnPhase != TURN_PHASE_PLAY_TILE {
		return errors.New("You are not allowed to play a tile at the moment")
	}
	// Check the tile is in the relevant player's hand
	newPlayerTiles, n := g.PlayerTiles[playerNum].Remove(t, 1)
	if n == 0 {
		return errors.New("You don't have that tile in your hand")
	}
	// Check the tile is not next to two or more safe corps
	if g.IsJoiningSafeCorps(t) {
		return errors.New(fmt.Sprintf(
			"You are not allowed to play %s as it would join two safe corporations",
			TileText(t)))
	}
	// Check if it's creating a new corp but there's none left
	if g.WouldFoundCorp(t) && len(g.InactiveCorps()) == 0 {
		return errors.New(fmt.Sprintf(
			"You are not allowed to play %s as there are no corporations available to found",
			TileText(t)))
	}
	// Remove the tile from the player's hand and place it
	g.PlayedTile = t
	g.PlayerTiles[playerNum] = newPlayerTiles
	g.Board[t.Row][t.Column] = TILE_UNINCORPORATED
	// Check for special actions based on adjacent tiles
	adjacentCorps := g.AdjacentCorps(t)
	if len(adjacentCorps) > 1 {
		// We have a merger
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"{{b}}%s{{_b}} played {{b}}%s{{_b}} which triggered a merger",
			g.RenderPlayer(playerNum), TileText(t))))
		potentialMergers := g.PotentialMergers(t)
		if len(potentialMergers) > 1 {
			g.TurnPhase = TURN_PHASE_MERGER_CHOOSE
		} else {
			g.ChooseMerger(t, potentialMergers[0][0], potentialMergers[0][1])
		}
	} else if len(adjacentCorps) == 1 {
		// Extending an existing corp
		g.SetAreaOnBoard(t, adjacentCorps[0])
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`{{b}}%s{{_b}} played {{b}}%s{{_b}} which increased the size of {{b}}%s{{_b}} to {{b}}%d{{_b}}`,
			g.RenderPlayer(playerNum), TileText(t),
			RenderCorp(adjacentCorps[0]), g.CorpSize(adjacentCorps[0]))))
		g.BuySharesPhase()
	} else if g.AdjacentToUnincorporated(t) {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`{{b}}%s{{_b}} played {{b}}%s{{_b}} to found a new corporation`,
			g.RenderPlayer(playerNum), TileText(t))))
		g.TurnPhase = TURN_PHASE_FOUND_CORP
	} else {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`{{b}}%s{{_b}} played {{b}}%s{{_b}}`,
			g.RenderPlayer(playerNum), TileText(t))))
		// Nothing adjacent
		g.BuySharesPhase()
	}
	return nil
}

func (g *Game) IsValidPlay(t Tile) bool {
	return !g.IsJoiningSafeCorps(t) && !(g.WouldFoundCorp(t) &&
		len(g.InactiveCorps()) == 0)
}

func (g *Game) WouldFoundCorp(t Tile) bool {
	adjacentCorps := g.AdjacentCorps(t)
	return len(adjacentCorps) == 0 && g.AdjacentToUnincorporated(t)
}

func (g *Game) BuySharesPhase() {
	if g.PlayerCanAffordShares(g.CurrentPlayer) || g.CanEnd(g.CurrentPlayer) {
		g.TurnPhase = TURN_PHASE_BUY_SHARES
		g.BoughtShares = 0
	} else {
		// Too poor
		g.NextPlayer()
	}
}

func (g *Game) PlayerCanAffordShares(playerNum int) bool {
	if g.PlayerCash[playerNum] < START_VALUE_LOW {
		return false
	}
	for _, corp := range Corps() {
		corpValue := g.CorpValue(corp)
		if g.BankShares[corp] > 0 && corpValue > 0 && corpValue <=
			g.PlayerCash[playerNum] {
			return true
		}
	}
	return false
}

func (g *Game) PayShareholderBonuses(corp int) {
	buf := bytes.NewBuffer([]byte{})
	// Shareholder bonuses
	majors := []int{}
	majorCount := 0
	minors := []int{}
	minorCount := 0
	stockMarketMessage := ""
	if len(g.Players) == 2 {
		// Special rule, play against stock market
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		stockMarketShares := r.Int()%6 + 1
		stockMarketMessage = fmt.Sprintf(
			"\n{{b}}The stock market{{_b}} rolled {{b}}%d{{_b}}, has {{b}}%d{{_b}} shares",
			stockMarketShares,
			stockMarketShares,
		)
		majors = append(majors, -1)
		majorCount = stockMarketShares
	}
	for pNum, _ := range g.Players {
		count := g.PlayerShares[pNum][corp]
		if count > 0 {
			if count > majorCount {
				// Shuffle down
				minors = majors
				minorCount = majorCount
				majors = []int{}
				majorCount = count
			}
			if count == majorCount {
				majors = append(majors, pNum)
			} else {
				if count > minorCount {
					minors = []int{}
					minorCount = count
				}
				if count == minorCount {
					minors = append(minors, pNum)
				}
			}
		}
	}
	majorBonus := float64(g.Corp1stBonus(corp))
	minorBonus := float64(g.Corp2ndBonus(corp))
	if len(majors) > 1 || len(minors) == 0 {
		majorBonus += minorBonus
	}
	buf.WriteString(fmt.Sprintf(
		`Paying shareholder bonuses for {{b}}%s{{_b}} (size {{b}}%d{{_b}}), major bonus is {{b}}$%d{{_b}}, minor bonus is {{b}}$%d{{_b}}.  Player share counts are as follows:%s`,
		RenderCorp(corp),
		g.CorpSize(corp),
		g.Corp1stBonus(corp),
		g.Corp2ndBonus(corp),
		stockMarketMessage,
	))
	for pNum, _ := range g.Players {
		buf.WriteString(fmt.Sprintf("\n{{b}}%s{{_b}}: {{b}}%d{{_b}}",
			g.RenderPlayer(pNum), g.PlayerShares[pNum][corp]))
	}

	// Pay major
	aveMajorBonus := int(math.Ceil(majorBonus/100/float64(len(majors)))) * 100
	for _, pNum := range majors {
		buf.WriteString(fmt.Sprintf(
			"\nPaid {{b}}%s{{_b}} a major bonus of {{b}}$%d{{_b}}",
			g.RenderPlayer(pNum), aveMajorBonus))
		if pNum >= 0 {
			g.PlayerCash[pNum] += aveMajorBonus
		}
	}
	// Pay minor if needed
	if len(majors) == 1 && len(minors) > 0 {
		aveMinorBonus := int(math.Ceil(minorBonus/100/float64(len(minors)))) *
			100
		for _, pNum := range minors {
			buf.WriteString(fmt.Sprintf(
				"\nPaid {{b}}%s{{_b}} a minor bonus of {{b}}$%d{{_b}}",
				g.RenderPlayer(pNum), aveMinorBonus))
			if pNum >= 0 {
				g.PlayerCash[pNum] += aveMinorBonus
			}
		}
	}
	g.Log.Add(log.NewPublicMessage(buf.String()))
}

func (g *Game) ChooseMerger(at Tile, from, into int) error {
	found := false
	for _, merger := range g.PotentialMergers(at) {
		if merger[0] == from && merger[1] == into {
			found = true
		}
	}
	if !found {
		return errors.New(fmt.Sprintf(
			"A merger from %s into %s is not available at the moment",
			CorpShortNames[from], CorpShortNames[into]))
	}
	g.TurnPhase = TURN_PHASE_MERGER
	g.MergerCurrentPlayer = g.CurrentPlayer
	g.MergerFromCorp = from
	g.MergerIntoCorp = into
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} (size {{b}}%d{{_b}}) is merging into {{b}}%s{{_b}} (size {{b}}%d{{_b}})`,
		RenderCorp(from), g.CorpSize(from),
		RenderCorp(into), g.CorpSize(into))))
	if g.Board[at.Row][at.Column] <= TILE_UNINCORPORATED {
		g.Board[at.Row][at.Column] = TILE_UNINCORPORATED
		g.SetAreaOnBoard(at, into)
	}
	g.PayShareholderBonuses(from)
	// Go to next player if merging player has no shares
	if g.PlayerShares[g.MergerCurrentPlayer][g.MergerFromCorp] == 0 {
		g.NextMergerPhasePlayer()
	}
	return nil
}

func (g *Game) SellSharesAction(playerNum, corp, amount int) error {
	if g.TurnPhase != TURN_PHASE_MERGER || g.MergerCurrentPlayer != playerNum {
		return errors.New("It's not your turn to sell shares")
	}
	if corp != g.MergerFromCorp {
		return errors.New("You can't sell shares in that corp")
	}
	if err := g.SellShares(playerNum, corp, amount); err != nil {
		return err
	}
	if g.PlayerShares[playerNum][corp] == 0 {
		g.NextMergerPhasePlayer()
	}
	return nil
}

func (g *Game) SellShares(playerNum, corp, amount int) error {
	if amount > g.PlayerShares[playerNum][corp] {
		return errors.New(fmt.Sprintf(`You only have %d shares`,
			g.PlayerShares[playerNum][corp]))
	}
	corpValue := g.CorpValue(corp)
	total := corpValue * amount
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} sold {{b}}%d{{_b}} shares in {{b}}%s{{_b}} for {{b}}$%d{{_b}} ({{b}}$%d{{_b}} per share)`,
		g.RenderPlayer(playerNum), amount, RenderCorp(corp), total, corpValue)))
	g.PlayerCash[playerNum] += total
	g.PlayerShares[playerNum][corp] -= amount
	g.BankShares[corp] += amount
	return nil
}

func (g *Game) TradeShares(playerNum, from, into, amount int) error {
	if g.TurnPhase != TURN_PHASE_MERGER || g.MergerCurrentPlayer != playerNum {
		return errors.New("It's not your turn to trade shares")
	}
	if from != g.MergerFromCorp {
		return errors.New("You can't trade shares from that corp")
	}
	if into != g.MergerIntoCorp {
		return errors.New("You can't trade shares into that corp")
	}
	if amount%2 != 0 {
		return errors.New(
			"You can only trade multiples of 2, trades are 2 for 1")
	}
	if amount > g.PlayerShares[playerNum][from] {
		return errors.New(fmt.Sprintf(`You only have %d shares`,
			g.PlayerShares[playerNum][from]))
	}
	if amount/2 > g.BankShares[into] {
		return errors.New(fmt.Sprintf(`The bank only has %d left in %s`,
			g.BankShares[into], CorpNames[into]))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} traded {{b}}%d{{_b}} shares in {{b}}%s{{_b}} ({{b}}$%d{{_b}} per share) for {{b}}%d{{_b}} shares in {{b}}%s{{_b}} ({{b}}$%d{{_b}} per share)`,
		g.RenderPlayer(playerNum), amount, RenderCorp(from), g.CorpValue(from),
		amount/2, RenderCorp(into), g.CorpValue(into))))
	g.PlayerShares[playerNum][from] -= amount
	g.BankShares[from] += amount
	g.PlayerShares[playerNum][into] += amount / 2
	g.BankShares[into] -= amount / 2
	if g.PlayerShares[playerNum][from] == 0 {
		g.NextMergerPhasePlayer()
	}
	return nil
}

func (g *Game) KeepShares(playerNum int) error {
	if g.TurnPhase != TURN_PHASE_MERGER || g.MergerCurrentPlayer != playerNum {
		return errors.New("It's not your turn to keep shares")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} kept remaining {{b}}%d{{_b}} shares in {{b}}%s{{_b}}`,
		g.RenderPlayer(playerNum), g.PlayerShares[playerNum][g.MergerFromCorp],
		RenderCorp(g.MergerFromCorp))))
	g.NextMergerPhasePlayer()
	return nil
}

func (g *Game) BuyShares(playerNum, corp, amount int) error {
	if g.TurnPhase != TURN_PHASE_BUY_SHARES || g.CurrentPlayer != playerNum {
		return errors.New("It's not your turn to buy shares")
	}
	if amount > MAX_BUY_PER_TURN-g.BoughtShares {
		return errors.New(fmt.Sprintf("You can only buy %d more this turn",
			MAX_BUY_PER_TURN-g.BoughtShares))
	}
	if amount > g.BankShares[corp] {
		return errors.New(fmt.Sprintf(
			"%s does not have %d left in the bank, has %d remaining",
			CorpNames[corp], amount, g.BankShares[corp]))
	}
	corpValue := g.CorpValue(corp)
	if corpValue == 0 {
		return errors.New(fmt.Sprintf(
			"Cannot buy shares in %s because they aren't active on the board",
			CorpNames[corp]))
	}
	total := corpValue * amount
	if g.PlayerCash[playerNum] < total {
		return errors.New(fmt.Sprintf("That would cost $%d, you only have $%d",
			total, g.PlayerCash[playerNum]))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} bought {{b}}%d{{_b}} shares in {{b}}%s{{_b}} for {{b}}$%d{{_b}} ({{b}}$%d{{_b}} per share)`,
		g.RenderPlayer(playerNum), amount, RenderCorp(corp), total, corpValue)))
	g.PlayerCash[playerNum] -= total
	g.BankShares[corp] -= amount
	g.PlayerShares[playerNum][corp] += amount
	g.BoughtShares += amount
	if g.BoughtShares == MAX_BUY_PER_TURN ||
		!g.PlayerCanAffordShares(playerNum) {
		g.NextPlayer()
	}
	return nil
}

func (g *Game) FoundCorp(playerNum, corp int) error {
	if g.TurnPhase != TURN_PHASE_FOUND_CORP || g.CurrentPlayer != playerNum {
		return errors.New("It's not your turn to buy shares")
	}
	if g.CorpSize(corp) > 0 {
		return errors.New(fmt.Sprintf("%s is already active on the board",
			CorpNames[corp]))
	}
	g.SetAreaOnBoard(g.PlayedTile, corp)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`{{b}}%s{{_b}} founded {{b}}%s{{_b}} at {{b}}%s{{_b}} (size {{b}}%d{{_b}})`,
		g.RenderPlayer(playerNum), RenderCorp(corp), TileText(g.PlayedTile),
		g.CorpSize(corp))))
	if g.BankShares[corp] > 0 {
		// Free share for founder
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`{{b}}%s{{_b}} received a free founder share in {{b}}%s{{_b}}`,
			g.RenderPlayer(playerNum), RenderCorp(corp))))
		g.BankShares[corp] -= 1
		g.PlayerShares[playerNum][corp] += 1
	}
	g.BuySharesPhase()
	return nil
}

func (g *Game) SetAreaOnBoard(t Tile, val int) {
	origVal := g.TileAt(t)
	if origVal == val {
		return
	}
	g.Board[t.Row][t.Column] = val
	for _, rawAdjT := range AdjacentTiles(t) {
		adjT := rawAdjT.(Tile)
		if g.TileAt(adjT) == origVal {
			g.SetAreaOnBoard(adjT, val)
		}
	}
}

func (g *Game) ConvertCorp(from, to int) {
	for _, r := range Rows() {
		for _, c := range Cols() {
			if g.Board[r][c] == from {
				g.Board[r][c] = to
			}
		}
	}
}

func (g *Game) NextMergerPhasePlayer() {
	g.MergerCurrentPlayer = (g.MergerCurrentPlayer + 1) % len(g.Players)
	if g.MergerCurrentPlayer == g.CurrentPlayer {
		g.ConvertCorp(g.MergerFromCorp, g.MergerIntoCorp)
		// Check if we have more mergers to do
		if potentialMergers := g.PotentialMergers(
			g.PlayedTile); len(potentialMergers) > 0 {
			if len(potentialMergers) > 1 {
				g.TurnPhase = TURN_PHASE_MERGER_CHOOSE
			} else {
				g.ChooseMerger(g.PlayedTile, potentialMergers[0][0],
					potentialMergers[0][1])
			}
		} else {
			g.BuySharesPhase()
		}
	} else if g.PlayerShares[g.MergerCurrentPlayer][g.MergerFromCorp] == 0 {
		g.NextMergerPhasePlayer()
	}
}

func (g *Game) NextPlayer() {
	if g.FinalTurn {
		g.EndGame()
	} else {
		// Draw tiles if needed
		g.DiscardUnplayableTiles(g.CurrentPlayer)
		if !g.DrawTiles(g.CurrentPlayer) {
			g.Log.Add(log.NewPublicMessage(
				"There aren't enough tiles left in the draw pile, it is the end of the game"))
			g.EndGame()
			return
		}
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		g.TurnPhase = TURN_PHASE_PLAY_TILE
		for {
			// If player can't play any tiles, discard all and draw 6 more.
			for _, tRaw := range g.PlayerTiles[g.CurrentPlayer] {
				if g.IsValidPlay(tRaw.(Tile)) {
					return
				}
			}
			// No valid plays, discard all and draw 6 more.
			discardStrs := []string{}
			for _, tRaw := range g.PlayerTiles[g.CurrentPlayer] {
				t := tRaw.(Tile)
				g.Board[t.Row][t.Column] = TILE_DISCARDED
				discardStrs = append(discardStrs, render.Bold(TileText(t)))
			}
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"%s can't play any tiles, discarding %s and drawing new tiles",
				g.RenderPlayer(g.CurrentPlayer),
				render.CommaList(discardStrs),
			)))
			g.PlayerTiles[g.CurrentPlayer] = card.Deck{}
			if !g.DrawTiles(g.CurrentPlayer) {
				g.Log.Add(log.NewPublicMessage(
					"There aren't enough tiles left in the draw pile, it is the end of the game"))
				g.EndGame()
				return
			}
		}
	}
}

func (g *Game) EndGame() {
	g.GameEnded = true
	// Pay out remaining corps on the board
	g.Log.Add(log.NewPublicMessage(
		"{{b}}It is the end of the game, now paying shareholder bonuses and selling all shares for active corporations.{{_b}}"))
	for _, corp := range Corps() {
		if g.CorpSize(corp) > 0 {
			g.PayShareholderBonuses(corp)
			for playerNum, _ := range g.Players {
				if g.PlayerShares[playerNum][corp] > 0 {
					g.SellShares(playerNum, corp,
						g.PlayerShares[playerNum][corp])
				}
			}
		}
	}
	buf := bytes.NewBufferString(fmt.Sprintf(
		"Final player cash is as follows:"))
	for playerNum, _ := range g.Players {
		buf.WriteString(fmt.Sprintf("\n{{b}}%s{{_b}}: {{b}}$%d{{_b}}",
			g.RenderPlayer(playerNum), g.PlayerCash[playerNum]))
	}
	g.Log.Add(log.NewPublicMessage(buf.String()))
}

func (g *Game) DiscardUnplayableTiles(playerNum int) {
	newHand := g.PlayerTiles[playerNum]
	discarded := []string{}
	for _, tRaw := range g.PlayerTiles[playerNum] {
		t := tRaw.(Tile)
		if g.IsJoiningSafeCorps(t) {
			newHand, _ = newHand.Remove(t, -1)
			g.Board[t.Row][t.Column] = TILE_DISCARDED
			discarded = append(discarded, fmt.Sprintf("{{b}}%s{{_b}}",
				TileText(t)))
		}
	}
	if len(discarded) > 0 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`{{b}}%s{{_b}} discarded %s and drew replacement tiles`,
			g.RenderPlayer(playerNum), render.CommaList(discarded))))
	}
	g.PlayerTiles[playerNum] = newHand
}

func (g *Game) DrawTiles(playerNum int) bool {
	drawNum := INIT_TILES - len(g.PlayerTiles[playerNum])
	if drawNum > len(g.BankTiles) {
		return false
	}
	if drawNum > 0 {
		var drawnTiles card.Deck
		drawnTiles, g.BankTiles = g.BankTiles.PopN(drawNum)
		tileStr := []string{}
		for _, t := range drawnTiles {
			tileStr = append(tileStr, render.Markup(
				TileText(t.(Tile)),
				render.Gray,
				true,
			))
		}
		g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
			"You drew %s",
			render.CommaList(tileStr),
		), []string{g.Players[playerNum]}))
		g.PlayerTiles[playerNum] =
			g.PlayerTiles[playerNum].PushMany(drawnTiles)
	}
	return true
}

func (g *Game) IsJoiningSafeCorps(t Tile) bool {
	safeCorps := 0
	for _, c := range g.AdjacentCorps(t) {
		if c != TILE_UNINCORPORATED && g.CorpSize(c) >= CORP_SAFE_SIZE {
			safeCorps += 1
			if safeCorps > 1 {
				return true
			}
		}
	}
	return false
}

func (g *Game) PotentialMergers(t Tile) [][2]int {
	potentialMergers := [][2]int{}
	from := []int{}
	fromSize := 0
	into := []int{}
	intoSize := 0
	// Figure out potential mergers based on sizes
	for _, corp := range g.AdjacentCorps(t) {
		if corp != TILE_UNINCORPORATED {
			size := g.CorpSize(corp)
			if size > intoSize {
				// Biggest so far, shuffle down
				from = into
				fromSize = intoSize
				into = []int{}
				intoSize = size
			}
			if size == intoSize {
				// Matches the current biggest size, potentially merged into
				into = append(into, corp)
			} else {
				if size > fromSize {
					// The new biggest from size
					from = []int{}
					fromSize = size
				}
				if size == fromSize {
					// The same at the other from sizes, potentially merged from
					from = append(from, corp)
				}
			}
		}
	}
	// Organise into relevant mergers, [2]int{from, into}
	if len(from)+len(into) < 2 {
		// No potential mergers
		return potentialMergers
	}
	if len(into) > 1 {
		// Potential mergers between same sized corps, calculate permutations
		for _, corp1 := range into {
			for _, corp2 := range into {
				if corp1 != corp2 {
					potentialMergers = append(potentialMergers,
						[2]int{corp1, corp2})
				}
			}
		}
		return potentialMergers
	}
	// One or more from
	for _, fromCorp := range from {
		potentialMergers = append(potentialMergers, [2]int{fromCorp, into[0]})
	}
	return potentialMergers
}

func (g *Game) AdjacentCorps(t Tile) []int {
	adjacentCorpMap := map[int]bool{}
	for _, adjTRaw := range AdjacentTiles(t) {
		adjT := adjTRaw.(Tile)
		if g.Board[adjT.Row][adjT.Column] > TILE_UNINCORPORATED {
			adjacentCorpMap[g.Board[adjT.Row][adjT.Column]] = true
		}
	}
	adjacentCorps := []int{}
	for c, _ := range adjacentCorpMap {
		adjacentCorps = append(adjacentCorps, c)
	}
	return adjacentCorps
}

func (g *Game) AdjacentToUnincorporated(t Tile) bool {
	for _, adjTRaw := range AdjacentTiles(t) {
		adjT := adjTRaw.(Tile)
		if g.Board[adjT.Row][adjT.Column] == TILE_UNINCORPORATED {
			return true
		}
	}
	return false
}

func (g *Game) ActiveCorps() []int {
	active := []int{}
	for _, corp := range Corps() {
		if g.CorpSize(corp) > 0 {
			active = append(active, corp)
		}
	}
	return active
}

func (g *Game) InactiveCorps() []int {
	active := []int{}
	for _, corp := range Corps() {
		if g.CorpSize(corp) == 0 {
			active = append(active, corp)
		}
	}
	return active
}

func (g *Game) TileAt(t Tile) int {
	return g.Board[t.Row][t.Column]
}

func (g *Game) CanEnd(playerNum int) bool {
	if g.FinalTurn || g.IsFinished() || g.CurrentPlayer != playerNum {
		return false
	}
	allSafe := true
	oneActive := false
	for _, corp := range Corps() {
		size := g.CorpSize(corp)
		if size >= CORP_END_GAME_SIZE {
			return true
		}
		oneActive = oneActive || size > 0
		allSafe = allSafe && (size == 0 || size >= CORP_SAFE_SIZE)
	}
	return oneActive && allSafe
}

func (g *Game) RenderPlayer(playerNum int) string {
	if playerNum == -1 {
		return "{{b}}the stock market{{_b}}"
	}
	return render.PlayerName(playerNum, g.Players[playerNum])
}

func RenderInCorpColour(corp int, text string) string {
	return fmt.Sprintf(`{{c "%s"}}%s{{_c}}`, CorpColours[corp], text)
}

func RenderCorp(corp int) string {
	return RenderInCorpColour(corp, CorpNames[corp])
}

func RenderCorpShort(corp int) string {
	return RenderInCorpColour(corp, CorpShortNames[corp])
}

func RenderCorpWithShort(corp int) string {
	return RenderInCorpColour(corp, fmt.Sprintf("%s (%s)", CorpNames[corp],
		CorpShortNames[corp]))
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
	matches := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, TILE_REGEXP)).
		FindStringSubmatch(strings.ToUpper(text))
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

func FindCorp(search string) (int, error) {
	search = strings.ToLower(search)
	for _, corp := range Corps() {
		if search == strings.ToLower(CorpNames[corp]) {
			return corp, nil
		}
		if search == strings.ToLower(CorpShortNames[corp]) {
			return corp, nil
		}
	}
	return 0, errors.New(fmt.Sprintf(
		"Could not find a corporation with the name %s", search))
}
