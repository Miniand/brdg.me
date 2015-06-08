package splendor

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	MaxGold = 5
)

const (
	PhaseMain = iota
	PhaseVisit
)

type Game struct {
	Players []string
	Log     *log.Log

	Decks  [3][]Card
	Board  [3][]Card
	Nobles []Noble
	Tokens Amount

	PlayerBoards []PlayerBoard

	CurrentPlayer int
	Phase         int

	EndTriggered bool
	Ended        bool
}

var LocRegexp = regexp.MustCompile(`^([\dA-Z])([\dA-Z])$`)

func (g *Game) Commands() []command.Command {
	return []command.Command{
		BuyCommand{},
		TakeCommand{},
		ReserveCommand{},
		VisitCommand{},
	}
}

func (g *Game) Name() string {
	return "Splendor"
}

func (g *Game) Identifier() string {
	return "splendor"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) < 2 || len(players) > 4 {
		return errors.New("must be between 2 and 4 players")
	}

	g.Players = players
	g.Log = log.New()

	g.Decks = [3][]Card{}
	g.Board = [3][]Card{}
	for l, cards := range [3][]Card{
		ShuffleCards(Level1Cards()),
		ShuffleCards(Level2Cards()),
		ShuffleCards(Level3Cards()),
	} {
		g.Board[l] = cards[:4]
		g.Decks[l] = cards[4:]
	}

	g.Nobles = ShuffleNobles(NobleCards())[:len(players)+1]

	g.Tokens = Amount{
		Gold: MaxGold,
	}
	maxGems := g.MaxGems()
	for _, r := range Gems {
		g.Tokens[r] = maxGems
	}

	g.PlayerBoards = []PlayerBoard{}
	for range players {
		g.PlayerBoards = append(g.PlayerBoards, NewPlayerBoard())
	}

	return nil
}

func (g *Game) MaxGems() int {
	switch len(g.Players) {
	case 2:
		return 4
	case 3:
		return 5
	default:
		return 7
	}
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Ended
}

func (g *Game) CheckEndTriggered() {
	if g.EndTriggered {
		return
	}
	for p, _ := range g.Players {
		if g.PlayerBoards[p].Prestige() >= 15 {
			g.Log.Add(log.NewPublicMessage(render.Bold(
				"The end of the game has been triggered",
			)))
			g.EndTriggered = true
			return
		}
	}
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winners := []string{}
	prestige := 0
	cards := 0
	for p, name := range g.Players {
		pp := g.PlayerBoards[p].Prestige()
		pc := len(g.PlayerBoards[p].Cards)
		if pp > prestige || (pp == prestige && pc < cards) {
			winners = []string{}
			prestige = pp
			cards = len(g.PlayerBoards[p].Cards)
		}
		if pp == prestige && pc == cards {
			winners = append(winners, name)
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}

func (g *Game) NextPhase() {
	switch g.Phase {
	case PhaseMain:
		g.VisitPhase()
	case PhaseVisit:
		g.NextPlayer()
	}
}

func (g *Game) VisitPhase() {
	g.Phase = PhaseVisit
	pb := g.PlayerBoards[g.CurrentPlayer]
	canVisit := []int{}
	for i, n := range g.Nobles {
		if pb.Bonuses().CanAfford(n.Cost) {
			canVisit = append(canVisit, i)
		}
	}
	switch len(canVisit) {
	case 0:
		g.NextPhase()
	case 1:
		g.Visit(g.CurrentPlayer, canVisit[0])
	}
}

func (g *Game) NextPlayer() {
	g.CheckEndTriggered()
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	if g.EndTriggered && g.CurrentPlayer == 0 {
		g.Ended = true
		return
	}
	g.MainPhase()
}

func (g *Game) MainPhase() {
	g.Phase = PhaseMain
}

func (g *Game) Pay(player int, amount Amount) error {
	if !g.PlayerBoards[player].CanAfford(amount) {
		return errors.New("can't afford that")
	}
	offset := g.PlayerBoards[player].Bonuses().Subtract(amount)
	for _, gem := range Gems {
		if offset[gem] < 0 {
			// Player didn't have enough just with bonuses
			g.PlayerBoards[player].Tokens[gem] += offset[gem]
			g.Tokens[gem] -= offset[gem]
			if g.PlayerBoards[player].Tokens[gem] < 0 {
				// Player didn't have enough normal tokens either, use gold
				g.PlayerBoards[player].Tokens[Gold] +=
					g.PlayerBoards[player].Tokens[gem]
				g.Tokens[gem] += g.PlayerBoards[player].Tokens[gem]
				g.Tokens[Gold] -= g.PlayerBoards[player].Tokens[gem]
				g.PlayerBoards[player].Tokens[gem] = 0
			}
		}
	}
	return nil
}

func ParseLoc(loc string) (row int, col int, err error) {
	matches := LocRegexp.FindStringSubmatch(strings.ToUpper(strings.TrimSpace(loc)))
	if matches == nil {
		return 0, 0, errors.New("invalid location, must be a number and a letter with no spaces")
	}
	if matches[2][0] >= '0' && matches[2][0] <= '9' {
		// Reversed order, swap back
		matches[1], matches[2] = matches[2], matches[1]
	}
	row, err = strconv.Atoi(matches[1])
	row -= 1
	if row < 0 || row > 3 {
		err = errors.New("row must be between 0 and 3")
	}
	col = int(matches[2][0] - 'A')
	return
}
