package for_sale

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players      []string
	BuildingDeck card.Deck
	ChequeDeck   card.Deck
	Log          *log.Log
}

func (g *Game) Name() string {
	return "For Sale"
}

func (g *Game) Identifier() string {
	return "for_sale"
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Start(players []string) error {
	if len(players) < 3 || len(players) > 5 {
		return errors.New("must have between 3 and 5 players")
	}
	g.Log = log.New()
	g.Players = players
	g.BuildingDeck = BuildingDeck().Shuffle()
	g.ChequeDeck = ChequeDeck().Shuffle()
	if len(players) == 3 {
		g.Log.Add(log.NewPublicMessage(
			"Removing two building and cheque cards for 3 player game"))
		_, g.BuildingDeck = g.BuildingDeck.PopN(2)
		_, g.ChequeDeck = g.ChequeDeck.PopN(2)
	}
	return nil
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	output := bytes.NewBuffer([]byte{})
	return output.String(), nil
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

func RegisterGobTypes() {
	gob.Register(card.SuitRankCard{})
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

func BuildingDeck() card.Deck {
	d := card.Deck{}
	for i := 1; i <= 20; i++ {
		d = d.Push(card.SuitRankCard{
			Rank: i,
		})
	}
	return d
}

func ChequeDeck() card.Deck {
	d := card.Deck{}
	for i := 1; i <= 20; i++ {
		c := card.SuitRankCard{
			Rank: i,
		}
		if i < 3 {
			c.Rank = 0
		}
		d = d.Push(c)
	}
	return d
}
