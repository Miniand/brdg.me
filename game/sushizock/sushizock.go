package sushizock

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players         []string
	CurrentPlayer   int
	Log             *log.Log
	BlueTiles       Tiles
	RedTiles        Tiles
	PlayerBlueTiles []Tiles
	PlayerRedTiles  []Tiles
	RolledDice      []int
	KeptDice        []int
	RemainingRolls  int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		RollCommand{},
		TakeCommand{},
		StealCommand{},
	}
}

func (g *Game) Name() string {
	return "Sushizock im Gockelwok"
}

func (g *Game) Identifier() string {
	return "sushizock"
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

func (g *Game) RenderForPlayer(player string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	// Dice
	diceCounts := g.DiceCounts()
	diceNumbers := make([]string, len(g.RolledDice))
	for i, _ := range g.RolledDice {
		diceNumbers[i] = fmt.Sprintf(`{{c "gray"}}%d{{_c}}`, i+1)
	}
	dice := append(BoldStrings(DiceStrings(g.RolledDice)),
		DiceStrings(g.KeptDice)...)
	cells := [][]string{
		dice,
		diceNumbers,
	}
	table, err := render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(fmt.Sprintf(
		"{{b}}Dice{{_b}}\n%s\n\n", table))
	// Tiles
	blueTilesStrs := g.BlueTiles.Strings()
	if diceCounts[DiceSushi] > 0 && diceCounts[DiceSushi] <= len(blueTilesStrs) {
		blueTilesStrs[diceCounts[DiceSushi]-1] = fmt.Sprintf(`{{b}}%s{{_b}}`,
			blueTilesStrs[diceCounts[DiceSushi]-1])
	}
	redTilesStrs := g.RedTiles.Strings()
	if diceCounts[DiceBones] > 0 && diceCounts[DiceBones] <= len(redTilesStrs) {
		redTilesStrs[diceCounts[DiceBones]-1] = fmt.Sprintf(`{{b}}%s{{_b}}`,
			redTilesStrs[diceCounts[DiceBones]-1])
	}
	cells = [][]string{
		blueTilesStrs,
		redTilesStrs,
	}
	table, err = render.Table(cells, 0, 1)
	if err != nil {
		return "", err
	}
	buf.WriteString(fmt.Sprintf(
		"{{b}}Tiles{{_b}}\n%s\n\n", table))
	// Players
	cells = [][]string{
		{`{{b}}Player{{_b}}`, `{{b}}Blue{{_b}}`, `{{b}}Red{{_b}}`},
	}
	for pNum, p := range g.Players {
		blueText := `{{c "gray"}}none{{_c}}`
		redText := blueText
		bLen := len(g.PlayerBlueTiles[pNum])
		if bLen > 0 {
			blueText = fmt.Sprintf(`%s {{c "gray"}}(%d tiles){{_c}}`,
				g.PlayerBlueTiles[pNum][bLen-1].String(), bLen)
		}
		rLen := len(g.PlayerRedTiles[pNum])
		if rLen > 0 {
			redText = fmt.Sprintf(`%s {{c "gray"}}(%d tiles){{_c}}`,
				g.PlayerRedTiles[pNum][rLen-1].String(), rLen)
		}
		cells = append(cells, []string{
			render.PlayerName(pNum, p),
			blueText,
			redText,
		})
	}
	table, err = render.Table(cells, 0, 2)
	if err != nil {
		return "", err
	}
	buf.WriteString(table)
	return buf.String(), nil
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 || len(players) > 5 {
		return errors.New("must be between 2 and 5 players")
	}
	g.Log = log.New()
	g.Players = players
	g.BlueTiles = ShuffleTiles(BlueTiles())
	g.RedTiles = ShuffleTiles(RedTiles())
	g.PlayerBlueTiles = make([]Tiles, len(g.Players))
	g.PlayerRedTiles = make([]Tiles, len(g.Players))
	for p, _ := range g.Players {
		g.PlayerBlueTiles[p] = Tiles{}
		g.PlayerRedTiles[p] = Tiles{}
	}
	g.StartTurn()
	return nil
}

func (g *Game) StartTurn() {
	g.RolledDice = RollDice(5)
	g.KeptDice = []int{}
	g.RemainingRolls = 2
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	g.StartTurn()
}

func (g *Game) Dice() []int {
	dice := []int{}
	dice = append(dice, g.RolledDice...)
	dice = append(dice, g.KeptDice...)
	return dice
}

func (g *Game) DiceCounts() map[int]int {
	return DiceCounts(g.Dice())
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return len(g.BlueTiles) == 0 && len(g.RedTiles) == 0
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	highestPlayers := []string{}
	highestScore := 0
	for pNum, p := range g.Players {
		pScore := Score(g.PlayerBlueTiles[pNum], g.PlayerRedTiles[pNum])
		if len(highestPlayers) == 0 || pScore > highestScore {
			highestPlayers = []string{}
			highestScore = pScore
		}
		if pScore == highestScore {
			highestPlayers = append(highestPlayers, p)
		}
	}
	return highestPlayers
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) CanTake(player int) bool {
	return g.CanTakeBlue(player) || g.CanTakeRed(player)
}

func (g *Game) CanTakeBlue(player int) bool {
	if player != g.CurrentPlayer {
		return false
	}
	diceCounts := g.DiceCounts()
	return diceCounts[DiceSushi] > 0 &&
		len(g.BlueTiles) >= diceCounts[DiceSushi]
}

func (g *Game) CanTakeRed(player int) bool {
	if player != g.CurrentPlayer {
		return false
	}
	diceCounts := g.DiceCounts()
	return diceCounts[DiceBones] > 0 &&
		len(g.RedTiles) >= diceCounts[DiceBones]
}

func (g *Game) TakeBlue(player int) error {
	if !g.CanTakeBlue(player) {
		return errors.New("unable to take blue at the moment")
	}
	t, remaining := g.BlueTiles.Remove(g.DiceCounts()[DiceSushi] - 1)
	g.PlayerBlueTiles[player] = append(g.PlayerBlueTiles[player], t)
	g.BlueTiles = remaining
	g.NextPlayer()
	return nil
}

func (g *Game) TakeRed(player int) error {
	if !g.CanTakeRed(player) {
		return errors.New("unable to take red at the moment")
	}
	t, remaining := g.RedTiles.Remove(g.DiceCounts()[DiceBones] - 1)
	g.PlayerRedTiles[player] = append(g.PlayerRedTiles[player], t)
	g.RedTiles = remaining
	g.NextPlayer()
	return nil
}

func (g *Game) PlayerNum(player string) (int, error) {
	pNum, err := helper.StringInStrings(player, g.Players)
	if err != nil {
		return 0, fmt.Errorf("could not find player '%s'", player)
	}
	return pNum, nil
}

func (g *Game) CanRoll(player int) bool {
	return g.CurrentPlayer == player && g.RemainingRolls > 0 &&
		len(g.RolledDice) > 1
}

func (g *Game) RollDice(player int, dice []int) error {
	if !g.CanRoll(player) {
		return errors.New("unable to roll at the moment")
	}
	rollMap := map[int]bool{}
	for _, d := range dice {
		if d < 1 || d > len(g.RolledDice) {
			return fmt.Errorf("%d is not a valid die number", d)
		}
		rollMap[d-1] = true
	}
	if len(rollMap) == len(g.RolledDice) {
		return fmt.Errorf("you must keep at least one die")
	}
	for i, d := range g.RolledDice {
		if !rollMap[i] {
			g.KeptDice = append(g.KeptDice, d)
		}
	}
	g.RolledDice = RollDice(len(rollMap))
	g.RemainingRolls -= 1
	if g.RemainingRolls == 0 || len(g.RolledDice) == 1 {
		g.KeptDice = append(g.KeptDice, g.RolledDice...)
		g.RolledDice = []int{}
		g.RemainingRolls = 0
		if !g.CanTake(player) && !g.CanSteal(player) {
			g.TakeWorst()
		}
	}
	return nil
}

func (g *Game) CanSteal(player int) bool {
	return g.CanStealBlue(player) || g.CanStealRed(player)
}

func (g *Game) AnotherPlayerHasBlue(player int) bool {
	for p, _ := range g.Players {
		if p != player && len(g.PlayerBlueTiles[p]) > 0 {
			return true
		}
	}
	return false
}

func (g *Game) AnotherPlayerHasRed(player int) bool {
	for p, _ := range g.Players {
		if p != player && len(g.PlayerRedTiles[p]) > 0 {
			return true
		}
	}
	return false
}

func (g *Game) CanStealBlue(player int) bool {
	return player == g.CurrentPlayer && g.AnotherPlayerHasBlue(player) &&
		g.DiceCounts()[DiceBlueChopsticks] >= 3
}

func (g *Game) CanStealRed(player int) bool {
	return player == g.CurrentPlayer && g.AnotherPlayerHasRed(player) &&
		g.DiceCounts()[DiceRedChopsticks] >= 3
}

func (g *Game) CanStealBlueN(player int) bool {
	return player == g.CurrentPlayer && g.AnotherPlayerHasBlue(player) &&
		g.DiceCounts()[DiceBlueChopsticks] >= 4
}

func (g *Game) CanStealRedN(player int) bool {
	return player == g.CurrentPlayer && g.AnotherPlayerHasRed(player) &&
		g.DiceCounts()[DiceRedChopsticks] >= 4
}

func (g *Game) StealRed(player, targetPlayer int) error {
	if !g.CanStealRed(player) {
		return errors.New("can't steal at the moment")
	}
	if player == targetPlayer {
		return errors.New("can't steal from yourself")
	}
	if len(g.PlayerRedTiles[targetPlayer]) == 0 {
		return errors.New("they don't have any red tiles to steal")
	}
	t, remaining := g.PlayerRedTiles[targetPlayer].Remove(
		len(g.PlayerRedTiles[targetPlayer]) - 1)
	g.PlayerRedTiles[player] = append(g.PlayerRedTiles[player], t)
	g.PlayerRedTiles[targetPlayer] = remaining
	g.NextPlayer()
	return nil
}

func (g *Game) StealBlue(player, targetPlayer int) error {
	if !g.CanStealBlue(player) {
		return errors.New("can't steal at the moment")
	}
	if player == targetPlayer {
		return errors.New("can't steal from yourself")
	}
	if len(g.PlayerBlueTiles[targetPlayer]) == 0 {
		return errors.New("they don't have any blue tiles to steal")
	}
	t, remaining := g.PlayerBlueTiles[targetPlayer].Remove(
		len(g.PlayerBlueTiles[targetPlayer]) - 1)
	g.PlayerBlueTiles[player] = append(g.PlayerBlueTiles[player], t)
	g.PlayerBlueTiles[targetPlayer] = remaining
	g.NextPlayer()
	return nil
}

func (g *Game) StealRedN(player, targetPlayer, n int) error {
	if !g.CanStealRed(player) {
		return errors.New("can't steal at the moment")
	}
	if player == targetPlayer {
		return errors.New("can't steal from yourself")
	}
	if len(g.PlayerRedTiles[targetPlayer]) == 0 {
		return errors.New("they don't have any red tiles to steal")
	}
	index := len(g.PlayerRedTiles[targetPlayer]) - n
	if index < 0 || index >= len(g.PlayerRedTiles[targetPlayer]) {
		return fmt.Errorf(
			"invalid tile number, you need to pick something between 1 and %d",
			len(g.PlayerRedTiles[targetPlayer]))
	}
	t, remaining := g.PlayerRedTiles[targetPlayer].Remove(index)
	g.PlayerRedTiles[player] = append(g.PlayerRedTiles[player], t)
	g.PlayerRedTiles[targetPlayer] = remaining
	g.NextPlayer()
	return nil
}

func (g *Game) StealBlueN(player, targetPlayer, n int) error {
	if !g.CanStealBlue(player) {
		return errors.New("can't steal at the moment")
	}
	if player == targetPlayer {
		return errors.New("can't steal from yourself")
	}
	if len(g.PlayerBlueTiles[targetPlayer]) == 0 {
		return errors.New("they don't have any blue tiles to steal")
	}
	index := len(g.PlayerBlueTiles[targetPlayer]) - n
	if index < 0 || index >= len(g.PlayerBlueTiles[targetPlayer]) {
		return fmt.Errorf(
			"invalid tile number, you need to pick something between 1 and %d",
			len(g.PlayerBlueTiles[targetPlayer]))
	}
	t, remaining := g.PlayerBlueTiles[targetPlayer].Remove(index)
	g.PlayerBlueTiles[player] = append(g.PlayerBlueTiles[player], t)
	g.PlayerBlueTiles[targetPlayer] = remaining
	g.NextPlayer()
	return nil
}

func (g *Game) TakeWorst() {
	g.NextPlayer()
}

func BoldStrings(strs []string) []string {
	bolded := make([]string, len(strs))
	for i, s := range strs {
		bolded[i] = fmt.Sprintf(`{{b}}%s{{_b}}`, s)
	}
	return bolded
}
