package starship_catan

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const (
	ResourceAny = iota
	ResourceFood
	ResourceFuel
	ResourceCarbon
	ResourceOre
	ResourceScience
	ResourceTrade
	ResourceAstro
	ResourceColonyShip
	ResourceTradeShip
	ResourceBooster
	ResourceCannon
)

const (
	PhaseChooseModule = iota
	PhaseProduce
	PhaseChooseSector
	PhaseTradeAndBuild
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players        []string
	PlayerBoards   [2]*PlayerBoard
	SectorCards    [4]card.Deck
	AdventureCards card.Deck
	Phase          int
	CurrentPlayer  int
	ProducePlayer  int
	Log            *log.Log
	YellowDice     int
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
	return []command.Command{
		ChooseModuleCommand{},
	}
}

func (g *Game) Name() string {
	return "Starship Catan"
}

func (g *Game) Identifier() string {
	return "starship_catan"
}

func RegisterGobTypes() {
	gob.Register(PlayerBoard{})
	gob.Register(ColonyCard{})
	gob.Register(TradeCard{})
	gob.Register(PirateCard{})
	gob.Register(AdventurePlanetCard{})
	for _, c := range ShuffledAdventureCards() {
		gob.Register(c)
	}
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
	players := []string{}
	switch g.Phase {
	case PhaseChooseModule:
		for p, pName := range g.Players {
			if len(g.PlayerBoards[p].Modules) == 0 {
				players = append(players, pName)
			}
		}
	default:
		players = append(players, g.Players[g.CurrentPlayer])
	}
	return players
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) ParsePlayer(player string) (int, error) {
	for pNum, p := range g.Players {
		if p == player {
			return pNum, nil
		}
	}
	return 0, fmt.Errorf(`could not find player with the name %s`, player)
}

func (g *Game) CanChoose(player int) bool {
	return g.Phase == PhaseChooseModule &&
		len(g.PlayerBoards[player].Modules) == 0
}

func (g *Game) Choose(player int, module int) error {
	if !g.CanChoose(player) {
		return errors.New("you can't choose a module at the moment")
	}
	g.PlayerBoards[player].Modules[module] = 1
	if len(g.WhoseTurn()) == 0 {
		g.NewTurn()
	}
	return nil
}

func (g *Game) NextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 2
	g.NewTurn()
}

func (g *Game) NewTurn() {
	g.Phase = PhaseChooseSector
	g.YellowDice = (r.Int() % 3) + 1
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s rolled a {{b}}%d{{_b}}, flight distance will be {{b}}%d{{_b}}`,
		g.RenderName(g.CurrentPlayer), g.YellowDice, g.FlightDistance())))
	g.Produce(g.CurrentPlayer)
}

func (g *Game) FlightDistance() int {
	return g.YellowDice +
		g.PlayerBoards[g.CurrentPlayer].Resources[ResourceBooster]
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.Players[player])
}

func (g *Game) Produce(player int) {
	g.Phase = PhaseProduce
	g.ProducePlayer = player
	if ContainsInt(g.YellowDice, TradeModuleDice(
		g.PlayerBoards[g.ProducePlayer].Modules[ModuleTrade], g.ProducePlayer)) {
		g.ProduceResource(g.ProducePlayer, ResourceTrade)
	}
}

func (g *Game) ProduceResource(player, resource int) {
	limit := 4
	if resource != ResourceScience {
		limit = 2 + g.PlayerBoards[player].Modules[ModuleLogistics]
	}
	if g.PlayerBoards[player].Resources[resource] < limit {
	}
}

func (g *Game) NextProducePlayer() {
	if g.ProducePlayer == g.CurrentPlayer {
		g.Produce((g.ProducePlayer + 1) % 2)
	} else {
		g.Phase = PhaseChooseSector
	}
}

func Itoas(in []int) []string {
	out := make([]string, len(in))
	for k, i := range in {
		out[k] = strconv.Itoa(i)
	}
	return out
}

func ContainsInt(needle int, haystack []int) bool {
	for _, i := range haystack {
		if i == needle {
			return true
		}
	}
	return false
}
