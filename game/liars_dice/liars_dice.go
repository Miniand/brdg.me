package liars_dice

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/die"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	START_DICE_COUNT = 5
)

type Game struct {
	Players       []string
	CurrentPlayer int
	PlayerDice    [][]int
	BidQuantity   int
	BidValue      int
	BidPlayer     int
	Log           *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		BidCommand{},
		CallCommand{},
	}
}

func (g *Game) Name() string {
	return "Liar's Dice"
}

func (g *Game) Identifier() string {
	return "liars_dice"
}

func (g *Game) GameLog() *log.Log {
	return g.Log
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

func (g *Game) PlayerNum(player string) (int, error) {
	for pNum, pName := range g.Players {
		if pName == player {
			return pNum, nil
		}
	}
	return 0, errors.New("Could not find player")
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	playerNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	currentBidText := `{{c "gray"}}first bid{{_c}}`
	if g.BidQuantity != 0 {
		currentBidText = RenderBid(g.BidQuantity, g.BidValue)
	}
	buf.WriteString(fmt.Sprintf("Current bid: %s\n", currentBidText))
	if len(g.PlayerDice[playerNum]) > 0 {
		buf.WriteString(fmt.Sprintf("Your dice: {{l}}%s{{_l}}\n\n",
			strings.Join(die.RenderDice(g.PlayerDice[playerNum]), " ")))
	}
	cells := [][]string{
		[]string{"{{b}}Player{{_b}}", "{{b}}Remaining dice{{_b}}"},
	}
	for pNum, p := range g.Players {
		cells = append(cells, []string{
			render.PlayerName(pNum, p),
			fmt.Sprintf("%d", len(g.PlayerDice[pNum])),
		})
	}
	table := render.Table(cells, 0, 1)
	buf.WriteString(table)
	return buf.String(), nil
}

func (g *Game) Start(players []string) error {
	// Set players
	if len(players) < 2 || len(players) > 6 {
		return errors.New("Liar's Dice must be between 2 and 6 players")
	}
	g.Players = players
	g.Log = log.New()
	// Set a random first player
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.CurrentPlayer = r.Int() % len(g.Players)
	// Initialise dice
	g.PlayerDice = make([][]int, len(g.Players))
	for pNum, _ := range g.Players {
		g.PlayerDice[pNum] = make([]int, START_DICE_COUNT)
	}
	// Kick off the first round
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.BidQuantity = 0
	g.RollDice()
}

func (g *Game) RollDice() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for pNum, _ := range g.Players {
		for d, _ := range g.PlayerDice[pNum] {
			g.PlayerDice[pNum][d] = (r.Int() % 6) + 1
		}
	}
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return len(g.ActivePlayers()) < 2
}

func (g *Game) Winners() []string {
	if g.IsFinished() {
		return []string{g.Players[g.ActivePlayers()[0]]}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) ActivePlayers() (players []int) {
	for pNum, _ := range g.Players {
		if len(g.PlayerDice[pNum]) > 0 {
			players = append(players, pNum)
		}
	}
	return
}

func (g *Game) NextActivePlayer(from int) int {
	next := (from + 1) % len(g.Players)
	for len(g.PlayerDice[next]) == 0 && next != from {
		next = (next + 1) % len(g.Players)
	}
	return next
}

func (g *Game) EliminatedPlayerList() (eliminated []string) {
	for pNum, p := range g.Players {
		if len(g.PlayerDice[pNum]) == 0 {
			eliminated = append(eliminated, p)
		}
	}
	return
}

func RenderBid(quantity int, value int) string {
	return fmt.Sprintf("%d {{l}}%s{{_l}}", quantity, die.Render(value))
}
