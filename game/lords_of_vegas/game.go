package lords_of_vegas

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	CasinoNone = iota
	CasinoAlbion
	CasinoSphynx
	CasinoVega
	CasinoTivoli
	CasinoPioneer
	CasinoTheStrip

	StartingCards = 2
	EndCardRange  = 8
)

type Game struct {
	Players []string
	Log     *log.Log

	CurrentPlayer   int
	Money           map[int]int
	Board           map[string]BoardSpaceState
	Deck            []string
	DeckPos         int
	EndCardPos      int
	HasGambled      bool
	ReorganisedLocs map[string]bool

	Finished bool
}

func (g *Game) Commands(player string) []command.Command {
	commands := []command.Command{}
	return commands
}

func (g *Game) Name() string {
	return "Lords of Vegas"
}

func (g *Game) Identifier() string {
	return "lords_of_vegas"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if l := len(players); l < 2 || l > 6 {
		return errors.New("only for 2-6 players")
	}
	g.Players = players
	g.Log = log.New()

	bspLen := len(BoardSpaces)
	g.Deck = make([]string, bspLen)
	for i, n := range rnd.Perm(bspLen) {
		g.Deck[i] = BoardSpaces[n].Location
	}

	// End card is placed three quarters through, but with some variance defined
	// by EndCardRange.
	g.EndCardPos = bspLen*3/4 - EndCardRange/2 + rnd.Int()%EndCardRange

	g.Money = map[int]int{}
	g.Board = map[string]BoardSpaceState{}

	// Draw two cards for each player for starting money and locations.
	for p := range g.Players {
		for i := 0; i < StartingCards; i++ {
			card, ok := g.DrawCard()
			if !ok {
				return errors.New("unable to draw starting card for player")
			}
			if len(g.Players) == 2 && !Valid2PLoc(card) {
				// This space isn't used in two player games, try again.
				i--
				continue
			}
			space, ok := BoardSpaceByLocation[card]
			if !ok {
				return fmt.Errorf("unable to find space for location %s", card)
			}
			g.Board[card] = BoardSpaceState{
				Owned: true,
				Owner: p,
			}
			g.Money[p] += space.StartingMoney
		}
	}

	g.CurrentPlayer = rnd.Int() % len(g.Players)

	return g.StartTurn()
}

func (g *Game) StartTurn() error {
	var (
		loc string
		bs  BoardSpace
		ok  bool
	)
	for {
		loc, ok = g.DrawCard()
		if !ok {
			return nil
		}
		bs, ok = BoardSpaceByLocation[loc]
		if !ok {
			return fmt.Errorf("could not find board space for %s", loc)
		}
		g.PayCasino(bs.PayCasino)
		if len(g.Players) > 2 || Valid2PLoc(loc) {
			break
		}
	}
	bsState := g.Board[loc]
	bsState.Owned = true
	bsState.Owner = g.CurrentPlayer
	g.Board[loc] = bsState
	g.HasGambled = false
	g.ReorganisedLocs = map[string]bool{}
	return nil
}

func (g *Game) PayCasino(casino int) {
	// TODO
}

func (g *Game) EndTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
}

// Valid2PLoc determines if a location is used in two player games.
func Valid2PLoc(loc string) bool {
	return !strings.HasPrefix(loc, "F")
}

func (g *Game) DrawCard() (string, bool) {
	if g.DeckPos < g.EndCardPos && g.DeckPos < len(g.Deck) {
		card := g.Deck[g.DeckPos]
		g.DeckPos++
		return card, true
	}
	g.EndGame()
	return "", false
}

func (g *Game) EndGame() {
	g.Finished = true
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Finished
}

func (g *Game) Winners() []string {
	return nil
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

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}
