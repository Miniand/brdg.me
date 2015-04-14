package farkle

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/die"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players       []string
	FirstPlayer   int
	Player        int
	Scores        map[int]int
	TurnScore     int
	RemainingDice []int
	TakenThisRoll bool
	Log           *log.Log
}

var DiceColours = map[int]string{
	1: render.Cyan,
	2: render.Green,
	3: render.Red,
	4: render.Blue,
	5: render.Yellow,
	6: render.Magenta,
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		TakeCommand{},
		RollCommand{},
		DoneCommand{},
	}
}

func (g *Game) Name() string {
	return "Farkle"
}

func (g *Game) Identifier() string {
	return "farkle"
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

func (g *Game) RenderForPlayer(player string) (string, error) {
	playerNum, err := g.GetPlayerNum(player)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	cells := [][]interface{}{}
	if playerNum == g.Player {
		cells = append(cells,
			[]interface{}{"{{b}}Remaining dice{{_b}}", RenderDice(g.RemainingDice)},
			[]interface{}{"{{b}}Score this turn{{_b}}", strconv.Itoa(g.TurnScore)})
	}
	cells = append(cells,
		[]interface{}{"{{b}}Your score{{_b}}", strconv.Itoa(g.Scores[playerNum])})
	t := render.Table(cells, 0, 1)
	buf.WriteString(t)
	buf.WriteString("\n\n")
	cells = [][]interface{}{
		[]interface{}{
			"{{b}}Player{{_b}}",
			"{{b}}Score{{_b}}",
		},
	}
	for playerNum, player := range g.Players {
		playerName := render.PlayerName(playerNum, player)
		if playerNum == g.FirstPlayer {
			playerName += " (started)"
		}
		cells = append(cells, []interface{}{
			playerName,
			strconv.Itoa(g.Scores[playerNum]),
		})
	}
	t = render.Table(cells, 0, 1)
	buf.WriteString(t)
	return buf.String(), nil
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 {
		return errors.New("Farkle requires at least two players")
	}
	g.Log = log.New()
	g.Scores = map[int]int{}
	g.Players = players
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Player = r.Int() % len(g.Players)
	g.FirstPlayer = g.Player
	g.StartTurn()
	return nil
}

func (g *Game) StartTurn() {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("It is now %s's turn",
		render.PlayerName(g.Player, g.Players[g.Player]))))
	g.TurnScore = 0
	g.TakenThisRoll = false
	g.Roll(6)
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	if g.Player != g.FirstPlayer {
		return false
	}
	for _, s := range g.Scores {
		if s >= 10000 {
			return true
		}
	}
	return false
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winners := []string{}
	winningScore := 0
	for playerNum, player := range g.Players {
		if g.Scores[playerNum] > winningScore {
			winners = []string{}
			winningScore = g.Scores[playerNum]
		}
		if g.Scores[playerNum] == winningScore {
			winners = append(winners, player)
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.Player]}
}

func (g *Game) NextPlayer() {
	g.Player = (g.Player + 1) % len(g.Players)
	if !g.IsFinished() {
		g.StartTurn()
	}
}

func (g *Game) Roll(n int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.RemainingDice = make([]int, n)
	for i := 0; i < n; i++ {
		g.RemainingDice[i] = r.Int()%6 + 1
	}
	sort.IntSlice(g.RemainingDice).Sort()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s rolled %s",
		render.PlayerName(g.Player, g.Players[g.Player]),
		RenderDice(g.RemainingDice))))
	if len(AvailableScores(g.RemainingDice)) == 0 {
		// No dice!
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s rolled no scoring dice and lost {{b}}%d{{_b}} points!",
			render.PlayerName(g.Player, g.Players[g.Player]),
			g.TurnScore)))
		g.NextPlayer()
	}
}

func RenderDie(value int) string {
	return render.Markup(die.Render(value), DiceColours[value], true)
}

func RenderDice(values []int) string {
	strs := make([]string, len(values))
	for i, v := range values {
		strs[i] = RenderDie(v)
	}
	return strings.Join(strs, " ")
}

func (g *Game) GetPlayerNum(player string) (int, error) {
	for playerNum, p := range g.Players {
		if p == player {
			return playerNum, nil
		}
	}
	return 0, errors.New("Could not find player " + player)
}
