package seven_wonders

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players []string
	Log     *log.Log

	Round int
	Hands []card.Deck
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "7 Wonders"
}

func (g *Game) Identifier() string {
	return "7_wonders"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if l := len(players); l < 3 || l > 7 {
		return errors.New("7 Wonders is 3 to 7 player")
	}
	g.Players = players
	g.Log = log.New()

	g.StartRound(1)

	return nil
}

func (g *Game) StartRound(round int) {
	players := len(g.Players)
	switch round {
	case 1:
		g.DealHands(DeckAge1(players).Shuffle())
	case 2:
		g.DealHands(DeckAge2(players).Shuffle())
	case 3:
		g.DealHands(DeckAge3(players).Shuffle())
	}
}

func (g *Game) DealHands(cards card.Deck) {
	players := len(g.Players)
	g.Hands = make([]card.Deck, players)
	per := cards.Len() / players
	for p := range g.Players {
		g.Hands[p], cards = cards.PopN(per)
	}
}

func (g *Game) PlayerList() []string {
	return nil
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

func (g *Game) PlayerNum(player string) (int, bool) {
	for pNum, p := range g.Players {
		if player == p {
			return pNum, true
		}
	}
	return 0, false
}

func (g *Game) PlayerName(player int) string {
	return render.PlayerName(player, g.Players[player])
}
