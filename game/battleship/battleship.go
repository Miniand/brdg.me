package battleship

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	STATE_PLACING = iota
	STATE_SHOOTING
)

const (
	LOCATION_EMPTY = iota
	SHIP_CARRIER
	SHIP_BATTLESHIP
	SHIP_CRUISER
	SHIP_SUBMARINE
	SHIP_DESTROYER
	LOCATION_HIT
	LOCATION_MISS
)

const (
	X_1 = iota
	X_2
	X_3
	X_4
	X_5
	X_6
	X_7
	X_8
	X_9
	X_10
)

const (
	Y_A = iota
	Y_B
	Y_C
	Y_D
	Y_E
	Y_F
	Y_G
	Y_H
	Y_I
	Y_J
)

const (
	DIRECTION_UP = iota
	DIRECTION_RIGHT
	DIRECTION_DOWN
	DIRECTION_LEFT
)

var ships = []int{
	SHIP_CARRIER,
	SHIP_BATTLESHIP,
	SHIP_CRUISER,
	SHIP_SUBMARINE,
	SHIP_DESTROYER,
}

var shipSizes = map[int]int{
	SHIP_CARRIER:    5,
	SHIP_BATTLESHIP: 4,
	SHIP_CRUISER:    3,
	SHIP_SUBMARINE:  3,
	SHIP_DESTROYER:  2,
}

var shipNames = map[int]string{
	SHIP_CARRIER:    "carrier",
	SHIP_BATTLESHIP: "battleship",
	SHIP_CRUISER:    "cruiser",
	SHIP_SUBMARINE:  "submarine",
	SHIP_DESTROYER:  "destroyer",
}

var directionNames = map[int]string{
	DIRECTION_UP:    "up",
	DIRECTION_DOWN:  "down",
	DIRECTION_LEFT:  "left",
	DIRECTION_RIGHT: "right",
}

var directionStrings = map[string]int{
	"up":    DIRECTION_UP,
	"down":  DIRECTION_DOWN,
	"left":  DIRECTION_LEFT,
	"right": DIRECTION_RIGHT,
	"north": DIRECTION_UP,
	"south": DIRECTION_DOWN,
	"west":  DIRECTION_LEFT,
	"east":  DIRECTION_RIGHT,
}

var tileOutputsSelf = map[int]string{
	LOCATION_EMPTY:  `{{c "blue"}}░░{{_c}}`,
	SHIP_CARRIER:    "██",
	SHIP_BATTLESHIP: "██",
	SHIP_CRUISER:    "██",
	SHIP_SUBMARINE:  "██",
	SHIP_DESTROYER:  "██",
	LOCATION_HIT:    `{{c "red"}}{{b}}╬╬{{_b}}{{_c}}`,
	LOCATION_MISS:   `{{c "gray"}}XX{{_c}}`,
}

var tileOutputsEnemy = map[int]string{
	LOCATION_EMPTY:  tileOutputsSelf[LOCATION_EMPTY],
	SHIP_CARRIER:    tileOutputsSelf[LOCATION_EMPTY],
	SHIP_BATTLESHIP: tileOutputsSelf[LOCATION_EMPTY],
	SHIP_CRUISER:    tileOutputsSelf[LOCATION_EMPTY],
	SHIP_SUBMARINE:  tileOutputsSelf[LOCATION_EMPTY],
	SHIP_DESTROYER:  tileOutputsSelf[LOCATION_EMPTY],
	LOCATION_HIT:    tileOutputsSelf[LOCATION_HIT],
	LOCATION_MISS:   tileOutputsSelf[LOCATION_MISS],
}

type Game struct {
	Players       []string
	CurrentPlayer int
	Log           *log.Log
	State         int
	Boards        [2][10][10]int
	LeftToPlace   [2][]int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlaceCommand{},
		ShootCommand{},
	}
}

func (g *Game) Name() string {
	return "Battleship"
}

func (g *Game) Identifier() string {
	return "battleship"
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
	p, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	}
	output := bytes.Buffer{}
	if g.State == STATE_PLACING {
		if len(g.LeftToPlace[p]) > 0 {
			output.WriteString(
				"{{b}}Ships left to place (ship size in brackets):{{_b}}")
			for _, s := range g.LeftToPlace[p] {
				output.WriteString(fmt.Sprintf("\n%s (%d)", shipNames[s],
					shipSizes[s]))
			}
		} else {
			output.WriteString(
				"{{b}}Waiting for your opponent to place their ships{{_b}}")
		}
		output.WriteString("\n\n")
	} else {
		output.WriteString("{{b}}Enemy board:{{_b}}\n\n")
		tiles := tileOutputsEnemy
		if g.IsFinished() {
			tiles = tileOutputsSelf
		}
		output.WriteString(RenderBoard(g.Boards[OtherPlayer(p)], tiles))
		output.WriteString("\n\n")
	}
	output.WriteString("{{b}}Your board:{{_b}}\n\n")
	output.WriteString(RenderBoard(g.Boards[p], tileOutputsSelf))
	return output.String(), nil
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Can only play with 2 players")
	}
	g.Players = players
	g.Log = log.New()
	g.InitBoards()
	g.InitLeftToPlace()
	return nil
}

func (g *Game) InitBoards() {
	g.Boards = [2][10][10]int{}
	for p := 0; p < 2; p++ {
		g.Boards[p] = [10][10]int{}
		for y := Y_A; y <= Y_J; y++ {
			g.Boards[p][y] = [10]int{}
		}
	}
}

func (g *Game) InitLeftToPlace() {
	g.LeftToPlace = [2][]int{}
	for p := 0; p < 2; p++ {
		g.LeftToPlace[p] = append(g.LeftToPlace[p], ships...)
	}
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	if g.State == STATE_PLACING {
		return false
	}
	for pNum, _ := range g.Players {
		if g.PlayerHitsRemaining(pNum) == 0 {
			return true
		}
	}
	return false
}

func (g *Game) PlayerHitsRemaining(player int) int {
	remaining := 0
	for _, r := range g.Boards[player] {
		for _, c := range r {
			if c >= SHIP_CARRIER && c <= SHIP_DESTROYER {
				remaining += 1
			}
		}
	}
	return remaining
}

func (g *Game) PlayerShipHitsRemaining(player, ship int) int {
	remaining := 0
	for _, r := range g.Boards[player] {
		for _, c := range r {
			if c == ship {
				remaining += 1
			}
		}
	}
	return remaining
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winner := 0
	if g.PlayerHitsRemaining(0) == 0 {
		winner = 1
	}
	return []string{g.Players[winner]}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	players := []string{}
	if g.State == STATE_PLACING {
		for pNum, p := range g.Players {
			if len(g.LeftToPlace[pNum]) > 0 {
				players = append(players, p)
			}
		}
	} else {
		players = append(players, g.Players[g.CurrentPlayer])
	}
	return players
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) CanPlace(player string) bool {
	if g.IsFinished() || g.State != STATE_PLACING {
		return false
	}
	return g.IsPlayersTurn(player)
}

func (g *Game) PlaceShip(player, ship, y, x, dir int) error {
	if !g.CanPlace(g.Players[player]) {
		return errors.New("You are not allowed to place a ship at the moment")
	}
	foundAt := -1
	for i, ps := range g.LeftToPlace[player] {
		if ps == ship {
			foundAt = i
			break
		}
	}
	if foundAt == -1 {
		return errors.New("You don't have any of that type of shift to place")
	}
	// Try to place
	locs := LocationsInDirection(y, x, dir, shipSizes[ship]-1)
	for _, l := range locs {
		if !IsValidLocation(l[0], l[1]) {
			return errors.New(
				"Can't place there because it would go off the board")
		}
		if g.Boards[player][l[0]][l[1]] != LOCATION_EMPTY {
			return errors.New(
				"Can't place there because there's a ship in the way")
		}
	}
	for _, l := range locs {
		g.Boards[player][l[0]][l[1]] = ship
	}
	// Remove from array and state change if needed
	if len(g.LeftToPlace[player]) == 1 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s finished placing their ships", g.RenderPlayerName(player))))
		g.LeftToPlace[player] = []int{}
		if len(g.LeftToPlace[OtherPlayer(player)]) == 0 {
			g.State = STATE_SHOOTING
		}
	} else {
		g.LeftToPlace[player] = append(g.LeftToPlace[player][:foundAt],
			g.LeftToPlace[player][foundAt+1:]...)
	}
	return nil
}

func (g *Game) CanShoot(player string) bool {
	if g.IsFinished() || g.State != STATE_SHOOTING {
		return false
	}
	return g.IsPlayersTurn(player)
}

func (g *Game) Shoot(player, y, x int) error {
	if !g.CanShoot(g.Players[player]) {
		return errors.New("You are not allowed to shoot at the moment")
	}
	if !IsValidLocation(y, x) {
		return errors.New("That is not a valid location on the board")
	}
	switch g.Boards[OtherPlayer(player)][y][x] {
	case LOCATION_HIT, LOCATION_MISS:
		return errors.New("You have already shot there previously")
	case LOCATION_EMPTY:
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s shot at %s and missed", g.RenderPlayerName(player),
			LocationName(y, x))))
		g.Boards[OtherPlayer(player)][y][x] = LOCATION_MISS
	default:
		ship := g.Boards[OtherPlayer(player)][y][x]
		g.Boards[OtherPlayer(player)][y][x] = LOCATION_HIT
		if g.PlayerShipHitsRemaining(OtherPlayer(player), ship) == 0 {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"{{b}}%s shot at %s and sunk a %s!{{_b}}",
				g.RenderPlayerName(player), LocationName(y, x),
				shipNames[ship])))
		} else {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"%s shot at %s and hit a ship", g.RenderPlayerName(player),
				LocationName(y, x))))
		}
	}
	g.NextPlayer()
	return nil
}

func (g *Game) RenderPlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = OtherPlayer(g.CurrentPlayer)
}

func (g *Game) PlayerFromString(s string) (int, error) {
	for pNum, p := range g.Players {
		if s == p {
			return pNum, nil
		}
	}
	return 0, errors.New("Could not find player with that name")
}

func (g *Game) IsPlayersTurn(player string) bool {
	for _, p := range g.WhoseTurn() {
		if p == player {
			return true
		}
	}
	return false
}

func OtherPlayer(p int) int {
	return (p + 1) % 2
}

func ParseShip(s string) (int, error) {
	if len(s) < 3 {
		return 0, errors.New("Ship name must be at least 3 characters long")
	}
	lower := strings.ToLower(strings.TrimSpace(s))
	for _, ship := range ships {
		if strings.HasPrefix(strings.ToLower(shipNames[ship]), lower) {
			return ship, nil
		}
	}
	return 0, errors.New("Could not find ship matching that name")
}

func ParseLocation(s string) (y, x int, err error) {
	if len(s) < 2 {
		return 0, 0, errors.New(
			"Location must be a letter immediately followed by a number, such as B10")
	}
	upper := strings.ToUpper(strings.TrimSpace(s))
	y = int(upper[0] - 'A')
	if y < Y_A || y > Y_J {
		return 0, 0, errors.New(
			"The first character of a location must be a letter between A and J")
	}
	x, err = strconv.Atoi(s[1:])
	if err != nil {
		return 0, 0, errors.New(
			"The location letter should be immediately followed by a number")
	}
	x -= 1 // Zero indexed
	if x < X_1 || x > X_10 {
		return 0, 0, errors.New(
			"The number in the location must be between 1 and 10")
	}
	return
}

func IsValidLocation(y, x int) bool {
	return y >= Y_A && y <= Y_J && x >= X_1 && x <= X_10
}

func ParseDirection(s string) (int, error) {
	if len(s) < 1 {
		return 0, errors.New("Direction must be at least 1 character long")
	}
	lower := strings.ToLower(strings.TrimSpace(s))
	for dirStr, dir := range directionStrings {
		if strings.HasPrefix(dirStr, lower) {
			return dir, nil
		}
	}
	return 0, errors.New("Could not find direction matching that name, please use up, down, left, right, north, south, west, east")
}

func LocationName(y, x int) string {
	return fmt.Sprintf("%c%d", 'A'+y, x+1)
}

func DirectionModifiers(dir int) (yMod, xMod int) {
	switch dir {
	case DIRECTION_UP:
		yMod = -1
	case DIRECTION_DOWN:
		yMod = 1
	case DIRECTION_LEFT:
		xMod = -1
	case DIRECTION_RIGHT:
		xMod = 1
	}
	return
}

func LocationsInDirection(y, x, dir, dist int) [][2]int {
	yMod, xMod := DirectionModifiers(dir)
	locs := [][2]int{}
	for i := 0; i <= dist; i++ {
		locs = append(locs, [2]int{y + i*yMod, x + i*xMod})
	}
	return locs
}

func RenderBoard(board [10][10]int, tiles map[int]string) string {
	output := bytes.Buffer{}
	output.WriteString("  1 2 3 4 5 6 7 8 9 10")
	for y, row := range board {
		output.WriteString(fmt.Sprintf("\n%c ", y+'A'))
		for _, cell := range row {
			output.WriteString(tiles[cell])
		}
		output.WriteString(fmt.Sprintf(" %c", y+'A'))
	}
	output.WriteString("\n  1 2 3 4 5 6 7 8 9 10")
	return output.String()
}
