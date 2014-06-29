package starship_catan

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type Game struct {
	Players        []string
	PlayerBoards   [2]*PlayerBoard
	SectorCards    [4]card.Deck
	AdventureCards card.Deck
	Log            *log.Log
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("this game requires two players")
	}
	g.Players = players
	g.PlayerBoards = [2]*PlayerBoard{
		NewPlayerBoard(0),
		NewPlayerBoard(1),
	}
	sectorCards := ShuffledSectorCards()
	g.SectorCards = [4]card.Deck{}
	for i := 0; i < 4; i++ {
		g.SectorCards[i], sectorCards = sectorCards.PopN(10)
	}
	g.AdventureCards = ShuffledAdventureCards()
	g.Log = log.New()
	return nil
}

func (g *Game) Commands() []command.Command {
	return []command.Command{}
}

func (g *Game) Name() string {
	return "Starship Catan"
}

func (g *Game) Identifier() string {
	return "starship_catan"
}

func RegisterGobTypes() {
	gob.Register(PlayerBoard{})
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

func (g *Game) RenderForPlayer(player string) (string, error) {
	return "", nil
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
