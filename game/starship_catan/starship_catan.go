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
	PhaseChooseModule = iota
	PhaseProduce
	PhaseChooseSector
	PhaseFlight
	PhaseTradeAndBuild
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players          []string
	PlayerBoards     [2]*PlayerBoard
	SectorCards      map[int]card.Deck
	SectorDrawPile   card.Deck
	FlightCards      card.Deck
	CurrentSector    int
	VisitedCards     card.Deck
	RemainingActions int
	RemainingMoves   int
	AdventureCards   card.Deck
	Phase            int
	CurrentPlayer    int
	GainPlayer       int
	GainResources    []int
	Log              *log.Log
	YellowDice       int
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
	g.SectorCards = map[int]card.Deck{}
	for i := 1; i <= 4; i++ {
		g.SectorCards[i], sectorCards = sectorCards.PopN(10)
	}
	g.SectorDrawPile = sectorCards
	g.FlightCards = card.Deck{}
	g.AdventureCards = ShuffledAdventureCards()
	g.Log = log.New()
	return nil
}

func (g *Game) Commands() []command.Command {
	commands := []command.Command{
		ChooseCommand{},
		GainCommand{},
		SectorCommand{},
	}
	if len(g.FlightCards) > 0 {
		c, _ := g.FlightCards.Pop()
		if c, ok := c.(Commander); ok {
			commands = append(commands, c.Commands()...)
		}
	}
	return commands
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
	gob.Register(MedianCard{})
	gob.Register(EmptyCard{})
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

func (g *Game) Choose(player, module int) error {
	if !g.CanChoose(player) {
		return errors.New("you can't choose a module at the moment")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s chose the {{b}}%s module{{_b}}`,
		g.RenderName(player), ModuleNames[module])))
	g.PlayerBoards[player].Modules[module] = 1
	if len(g.WhoseTurn()) == 0 {
		g.NewTurn()
	}
	return nil
}

func (g *Game) CanSector(player int) bool {
	return g.Phase == PhaseChooseSector && g.CurrentPlayer == player
}

func (g *Game) Sector(player, sector int) error {
	if !g.CanSector(player) {
		return errors.New("you can't choose a sectore at the moment")
	}
	if sector < 1 || sector > 4 {
		return errors.New("sector must be between 1 and 4")
	}
	g.Phase = PhaseFlight
	g.CurrentSector = sector
	g.RemainingActions = 2 + g.PlayerBoards[player].Modules[ModuleCommand]
	g.RemainingMoves = g.FlightDistance()
	return g.NextSectorCard()
}

func (g *Game) NextSectorCard() error {
	var nextCard card.Card
	if g.Phase != PhaseFlight {
		return errors.New("it isn't the flight phase at the moment")
	}
	if g.SectorCards[g.CurrentSector] == nil {
		return fmt.Errorf("%d is not a valid sector number", g.CurrentSector)
	}
	if len(g.SectorCards[g.CurrentSector]) == 0 ||
		len(g.FlightCards) >= g.FlightDistance() || g.RemainingActions == 0 {
		return g.EndFlight()
	}
	nextCard, g.SectorCards[g.CurrentSector] =
		g.SectorCards[g.CurrentSector].Pop()
	g.FlightCards = g.FlightCards.Push(nextCard)
	cardText := fmt.Sprintf("%#v", nextCard)
	if nextCard, ok := nextCard.(fmt.Stringer); ok {
		cardText = nextCard.String()
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`You have arrived at %s`, cardText)))
	if _, ok := nextCard.(Commander); !ok {
		g.Log.Add(log.NewPublicMessage(
			"There is nothing you can do here, continuing flight"))
		return g.NextSectorCard()
	}
	return nil
}

func (g *Game) EndFlight() error {
	if g.Phase != PhaseFlight {
		return errors.New("it isn't the flight phase at the moment")
	}
	if g.SectorCards[g.CurrentSector] == nil {
		return fmt.Errorf("%d is not a valid sector number", g.CurrentSector)
	}
	g.Log.Add(log.NewPublicMessage("The flight has ended"))
	g.SectorCards[g.CurrentSector] = g.SectorCards[g.CurrentSector].PushMany(
		g.FlightCards).Shuffle()
	g.FlightCards = card.Deck{}
	g.Phase = PhaseTradeAndBuild
	return nil
}

func (g *Game) CanFound(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight ||
		len(g.FlightCards) == 0 ||
		g.PlayerBoards[player].Resources[ResourceColonyShip] == 0 {
		return false
	}
	c, _ := g.FlightCards.Pop()
	_, ok := c.(ColonyCard)
	return ok
}

func (g *Game) Found(player int) error {
	var c card.Card

	if !g.CanFound(player) {
		return errors.New("you are not able to found a colony")
	}
	c, g.FlightCards = g.FlightCards.Pop()
	colCard := c.(ColonyCard)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s founded a colony on %s`, g.RenderName(player), colCard)))
	g.PlayerBoards[player].Colonies = g.PlayerBoards[player].Colonies.Push(c)
	if len(g.SectorDrawPile) > 0 {
		c, g.SectorDrawPile = g.SectorDrawPile.Pop()
		g.FlightCards = g.FlightCards.Push(c)
		if cStr, ok := c.(fmt.Stringer); ok {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`The replacement card is %s`, cStr)))
		}
	} else {
		g.Log.Add(log.NewPublicMessage("No replacement cards remain"))
	}
	g.RemainingActions -= 1
	g.NextSectorCard()
	return nil
}

func (g *Game) CanNext(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseFlight
}

func (g *Game) Next(player int) error {
	if !g.CanNext(player) {
		return errors.New("you can't advance to the next card")
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s continued their flight without taking an action`,
		g.RenderName(player))))
	return g.NextSectorCard()
}

func (g *Game) NextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 2
	g.NewTurn()
}

func (g *Game) NewTurn() {
	g.Phase = PhaseProduce
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

func (g *Game) CanGain(player int) bool {
	return g.GainPlayer == player && g.GainResources != nil
}

func (g *Game) GainOne(player int, resources []int) {
	if len(resources) == 0 {
		g.Gained(player)
	}
	canProduce := []int{}
	for _, r := range resources {
		if g.PlayerBoards[player].Resources[r] < g.ResourceLimit(player, r) {
			canProduce = append(canProduce, r)
		}
	}
	switch len(canProduce) {
	case 0:
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s did not gain a resource, all full", g.RenderName(player))))
		g.Gained(player)
	case 1:
		g.GainResource(player, canProduce[0])
		g.Gained(player)
	default:
		g.GainPlayer = player
		g.GainResources = canProduce
	}
}

func (g *Game) Gained(player int) {
	g.GainResources = nil
	if g.Phase == PhaseProduce {
		if player == g.CurrentPlayer {
			g.Produce((g.CurrentPlayer + 1) % 2)
		} else {
			g.Phase = PhaseChooseSector
		}
	}
}

func (g *Game) Produce(player int) {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf("%s is producing resources",
		g.RenderName(player))))
	if ContainsInt(g.YellowDice, TradeModuleDice(
		g.PlayerBoards[player].Modules[ModuleTrade], player)) {
		g.GainResource(player, ResourceTrade)
	}
	if ContainsInt(g.YellowDice, ScienceModuleDice(
		g.PlayerBoards[player].Modules[ModuleScience], player)) {
		g.GainResource(player, ResourceScience)
	}
	producing := []int{}
	producingMap := map[int]bool{}
	for _, c := range g.PlayerBoards[player].Colonies {
		col := c.(ColonyCard)
		if col.Dice == g.YellowDice && !producingMap[col.Resource] {
			producing = append(producing, col.Resource)
			producingMap[col.Resource] = true
		}
	}
	if len(producing) == 0 {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s doesn't produce anything with their colonies`,
			g.RenderName(player))))
		g.Gained(player)
	} else {
		g.GainOne(player, producing)
	}
}

func (g *Game) ResourceLimit(player, resource int) int {
	limit := 4
	if resource != ResourceScience {
		limit = 2 + g.PlayerBoards[player].Modules[ModuleLogistics]
	}
	return limit
}

func (g *Game) GainResource(player, resource int) {
	if g.PlayerBoards[player].Resources[resource] <
		g.ResourceLimit(player, resource) {
		g.PlayerBoards[player].Resources[resource] += 1
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained {{b}}1 %s{{_b}}", g.RenderName(player),
			ResourceNames[resource])))
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
