package battleship

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	STATE_PLACING = iota
	STATE_SHOOTING
)

const (
	SHIP_CARRIER = iota
	SHIP_BATTLESHIP
	SHIP_CRUISER
	SHIP_SUBMARINE
	SHIP_DESTROYER
)

const (
	X_0 = iota
	X_1
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

type Game struct {
	Players       []string
	CurrentPlayer int
	Log           *log.Log
	State         int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
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
	return "", nil
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("Can only play with 2 players")
	}
	g.Players = players
	g.Log = log.New()
	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return false
}

func (g *Game) Winners() []string {
	if g.IsFinished() {
		return []string{}
	}
	return []string{}
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

func (g *Game) CanPlace(player string) bool {
	if g.IsFinished() || g.State != STATE_PLACING {
		return false
	}
	return g.IsPlayersTurn(player)
}

func (g *Game) PlaceShip(player, ship, y, x, dir int) error {
	return nil
}

func (g *Game) CanShoot(player string) bool {
	if g.IsFinished() || g.State != STATE_SHOOTING {
		return false
	}
	return g.IsPlayersTurn(player)
}

func (g *Game) Shoot(player, y, x int) error {
	return nil
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

func ParseShip(s string) (int, error) {
	return SHIP_BATTLESHIP, nil
}

func ParseLocation(s string) (y, x int, err error) {
	return
}

func ParseDirection(s string) (int, error) {
	return DIRECTION_UP, nil
}
