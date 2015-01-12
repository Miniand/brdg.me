package splendor

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	MaxGold = 5
)

type Game struct {
	Players []string
	Log     *log.Log

	Decks  [3][]Card
	Board  [3][]Card
	Nobles []Noble
	Tokens Amount

	PlayerBoards []PlayerBoard
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
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
	return false
}

func (g *Game) Winners() []string {
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	return []string{}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, error) {
	return helper.StringInStrings(player, g.Players)
}
