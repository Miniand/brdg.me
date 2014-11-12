package starship_catan

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/helper"
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
	Players             []string
	PlayerBoards        [2]*PlayerBoard
	SectorCards         map[int]card.Deck
	SectorDrawPile      card.Deck
	Peeking             card.Deck
	FlightCards         card.Deck
	FlightActions       map[int]bool
	CurrentSector       int
	VisitedCards        card.Deck
	TradeAmount         int
	PlayerTradeAmount   int
	AdventureCards      card.Deck
	RemoveAdventureCard int
	Phase               int
	CurrentPlayer       int
	GainPlayer          int
	GainResources       []int
	GainQueue           [][]int
	Log                 *log.Log
	YellowDice          int
	CardFinished        bool
	LosingModule        bool
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
	g.Peeking = card.Deck{}
	g.AdventureCards = ShuffledAdventureCards()
	g.GainQueue = [][]int{}
	g.Log = log.New()
	return nil
}

func (g *Game) Commands() []command.Command {
	commands := []command.Command{}
	if g.GainResources == nil && len(g.FlightCards) > 0 && !g.CardFinished {
		c, _ := g.FlightCards.Pop()
		if c, ok := c.(Commander); ok {
			commands = append(commands, c.Commands()...)
		}
	}
	return append(
		commands,
		ChooseCommand{},
		GainCommand{},
		PutCommand{},
		SectorCommand{},
		BuildCommand{},
		UpgradeCommand{},
		TradePhaseBuyCommand{},
		TradePhaseSellCommand{},
		TakeCommand{},
		DoneCommand{},
		NextCommand{},
		EndCommand{},
	)
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

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	for _, b := range g.PlayerBoards {
		if b.VictoryPoints() >= 10 {
			return true
		}
	}
	return false
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	p1 := g.PlayerBoards[0].VictoryPoints()
	p2 := g.PlayerBoards[1].VictoryPoints()
	switch {
	case p1 > p2:
		return []string{g.Players[0]}
	case p2 > p1:
		return []string{g.Players[0]}
	default:
		return g.Players
	}
}

func (g *Game) WhoseTurn() []string {
	if g.GainResources != nil {
		return []string{g.Players[g.GainPlayer]}
	}
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

func (g *Game) Actions() int {
	return g.PlayerBoards[g.CurrentPlayer].Actions()
}

func (g *Game) RemainingActions() int {
	usedActions := 0
	for _, a := range g.FlightActions {
		if a {
			usedActions += 1
		}
	}
	return g.Actions() - usedActions
}

func (g *Game) RemainingMoves() int {
	return g.FlightDistance() - g.FlightCards.Len()
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
		g.RemainingMoves() <= 0 || g.RemainingActions() <= 0 {
		return g.EndFlight()
	}
	nextCard, g.SectorCards[g.CurrentSector] =
		g.SectorCards[g.CurrentSector].Pop()
	g.FlightCards = g.FlightCards.Push(nextCard)

	g.TradeAmount = 0
	g.CardFinished = false

	cardText := fmt.Sprintf("%#v", nextCard)
	if nextCard, ok := nextCard.(fmt.Stringer); ok {
		cardText = nextCard.String()
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s arrived at %s`, g.RenderName(g.CurrentPlayer), cardText)))
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
	g.PlayerBoards[g.CurrentPlayer].LastSectors = append(
		[]int{g.CurrentSector}, g.PlayerBoards[g.CurrentPlayer].LastSectors...)
	g.TradeAmount = 0
	g.PlayerTradeAmount = 0
	g.Phase = PhaseTradeAndBuild
	return nil
}

func (g *Game) RemainingPlayerTrades() int {
	return g.PlayerBoards[g.CurrentPlayer].Modules[ModuleTrade] -
		g.PlayerTradeAmount
}

func (g *Game) RemainingTrades() int {
	return 2 - g.TradeAmount
}

func (g *Game) ReplaceCard() error {
	var c card.Card
	if len(g.SectorDrawPile) > 0 {
		c, g.SectorDrawPile = g.SectorDrawPile.Pop()
		g.FlightCards = g.FlightCards.Push(c)
		g.CardFinished = true
		str := ""
		switch t := c.(type) {
		case FullStringer:
			str = t.FullString()
		case fmt.Stringer:
			str = t.String()
		}
		if str != "" {
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				`The replacement card is %s`, str)))
		}
	} else {
		g.Log.Add(log.NewPublicMessage("No replacement cards remain"))
		return g.EndFlight()
	}
	return nil
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

func (g *Game) Completed() {
	if g.RemoveAdventureCard > 0 {
		currentLen := g.CurrentAdventureCards().Len()
		acIndex := len(g.AdventureCards) - (currentLen - g.RemoveAdventureCard) - 1
		ac := g.AdventureCards[acIndex]
		g.AdventureCards = append(g.AdventureCards[:acIndex],
			g.AdventureCards[acIndex+1:]...)
		g.PlayerBoards[g.CurrentPlayer].CompletedAdventures =
			g.PlayerBoards[g.CurrentPlayer].CompletedAdventures.Push(ac)
		g.RemoveAdventureCard = 0
		g.RecalculatePeopleCards()
	}
	g.CardFinished = true
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

func (g *Game) GainResource(player, resource int) {
	t := Transaction{resource: 1}
	gained := g.PlayerBoards[player].Transact(t)
	if !gained.IsEmpty() {
		g.LogTransaction(player, gained)
	}
}

func (g *Game) TradableResources() []int {
	switch {
	case g.Phase == PhaseFlight && g.FlightCards.Len() > 0:
		card, _ := g.FlightCards.Pop()
		tradeCard, ok := card.(TradeCard)
		if ok {
			return tradeCard.Resources
		}
	case g.Phase == PhaseTradeAndBuild:
		resources := []int{}
		for r, _ := range g.PlayerBoards[g.CurrentPlayer].TradingPostPrices() {
			resources = append(resources, r)
		}
		return resources
	}
	return []int{}
}

func (g *Game) CanTrade(player, resource, amount int) (ok bool, price int, reason string) {
	if g.CurrentPlayer != player {
		return false, 0, "it's not your turn"
	}
	tradeDir := AmountTradeDir(amount)
	tradableResources := g.TradableResources()
	switch g.Phase {
	case PhaseFlight:
		if g.FlightCards.Len() == 0 {
			return false, 0, "there are no flight cards"
		}
		currentRaw, _ := g.FlightCards.Pop()
		trade, ok := currentRaw.(TradeCard)
		if !ok {
			return false, 0, "the current flight card is not a trade card"
		}
		if resource != ResourceAny &&
			!ContainsInt(resource, tradableResources) {
			return false, 0, fmt.Sprintf(
				"you can only %s %s with this trade card",
				TradeDirStrings[tradeDir],
				strings.Join(ResourceNameArr(tradableResources), ", "),
			)
		}
		if tradeDir != TradeDirBoth &&
			trade.Direction != TradeDirBoth &&
			tradeDir != trade.Direction {
			return false, 0, fmt.Sprintf(
				"you can only %s with this trade card",
				TradeDirStrings[tradeDir],
			)
		}
		targetAmount := amount*tradeDir + g.TradeAmount
		if amount != 0 && trade.Maximum != 0 && targetAmount > trade.Maximum {
			return false, 0, fmt.Sprintf(
				"you can only trade %s with this trade card, you have already traded %d",
				trade.AmountLimitString(),
				g.TradeAmount,
			)
		}
		if tradeDir == TradeDirBuy {
			if amount*trade.Price > g.PlayerBoards[player].Resources[ResourceAstro] {
				return false, 0, fmt.Sprintf(
					"you only have %s",
					RenderMoney(g.PlayerBoards[player].Resources[ResourceAstro]),
				)
			}
			if resource != ResourceAny {
				t := Transaction{resource: amount}
				if !g.PlayerBoards[player].CanFit(t) {
					return false, 0, t.CannotFitError().Error()
				}
			}
		}
		if tradeDir == TradeDirSell && resource != ResourceAny &&
			amount*tradeDir > g.PlayerBoards[player].Resources[resource] {
			return false, 0, fmt.Sprintf(
				"you only have %d %s",
				g.PlayerBoards[player].Resources[resource],
				ResourceNames[resource],
			)
		}
		return true, trade.Price, ""
	case PhaseTradeAndBuild:
		if g.RemainingTrades() == 0 {
			return false, 0, "you have already done two trades this phase"
		}
		if resource == ResourceAny && len(tradableResources) == 0 {
			return false, 0, "you don't have any trading posts"
		}
		if resource != ResourceAny {
			if !ContainsInt(resource, tradableResources) {
				return false, 0, "you don't have any trading posts for that resource"
			}
			prices := g.PlayerBoards[player].TradingPostPrices()
			if tradeDir == TradeDirBuy {
				t := Transaction{resource: amount}
				if !g.PlayerBoards[player].CanFit(t) {
					return false, 0, t.CannotFitError().Error()
				}
				if prices[resource].Buy > 0 {
					return true, prices[resource].Buy, ""
				}
				return false, 0, "you aren't able to buy that resource"
			}
			if tradeDir == TradeDirSell {
				if amount*tradeDir > g.PlayerBoards[player].Resources[resource] {
					return false, 0, fmt.Sprintf(
						"you only have %d %s",
						g.PlayerBoards[player].Resources[resource],
						ResourceNames[resource],
					)
				}
				if prices[resource].Sell > 0 {
					return true, prices[resource].Sell, ""
				}
				return false, 0, "you aren't able to sell that resource"
			}
		}
		return true, 0, ""
	}
	return false, 0, "it is not the correct phase to trade"
}

func (g *Game) Trade(player, resource, amount int) error {
	tradeDir := AmountTradeDir(amount)
	ok, price, reason := g.CanTrade(player, resource, amount)
	if !ok {
		return errors.New(reason)
	}
	if resource == ResourceAny {
		return errors.New("you must specify which resource to trade")
	}
	if tradeDir == TradeDirBoth {
		return errors.New("you must either buy or sell when trading")
	}

	total := amount * price
	g.PlayerBoards[player].Resources[ResourceAstro] -= total
	g.PlayerBoards[player].Resources[resource] += amount
	switch g.Phase {
	case PhaseFlight:
		g.MarkCardActioned()
		g.TradeAmount += amount * tradeDir
	case PhaseTradeAndBuild:
		g.TradeAmount += 1
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s %s %d %s for %s`,
		g.RenderName(player),
		TradeDirPastStrings[tradeDir],
		amount*tradeDir,
		RenderResource(resource),
		RenderMoney(total*tradeDir),
	)))
	return nil
}

func (g *Game) MarkCardActioned() {
	g.FlightActions[g.FlightCards.Len()] = true
}

func (g *Game) HandleTradeCommand(player string, args []string, tradeDir int) error {
	p, err := g.ParsePlayer(player)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("you must specify an amount")
	}
	amount, err := strconv.Atoi(args[0])
	if err != nil || amount <= 0 {
		return errors.New("the amount must be a positive whole number")
	}

	tradableResources := g.TradableResources()
	var resource int
	if len(args) > 1 {
		resource, err = helper.MatchStringInStringMap(
			args[1],
			ResourceNameMap(tradableResources),
		)
		if err != nil {
			return err
		}
	} else {
		if len(tradableResources) == 1 {
			resource = tradableResources[0]
		}
	}
	if resource == 0 {
		return errors.New("you must specify a resource")
	}
	return g.Trade(p, resource, amount*tradeDir)
}

func (g *Game) CurrentAdventureCards() card.Deck {
	n := 3
	if l := g.AdventureCards.Len(); l < n {
		n = l
	}
	cards, _ := g.AdventureCards.PopN(n)
	return cards
}

func (g *Game) RecalculatePeopleCards() {
	// Diplomat points
	p1D := g.PlayerBoards[0].DiplomatPoints()
	p2D := g.PlayerBoards[1].DiplomatPoints()
	d := -1
	switch {
	case p1D > 3 && p1D > p2D:
		d = 0
	case p2D > 3 && p2D > p1D:
		d = 1
	}
	// Medals
	p1M := g.PlayerBoards[0].Medals()
	p2M := g.PlayerBoards[1].Medals()
	m := -1
	switch {
	case p1M > 3 && p1M > p2M:
		m = 0
	case p2M > 3 && p2M > p1M:
		m = 1
	}
	// Apply
	for p, _ := range g.PlayerBoards {
		g.PlayerBoards[p].FriendOfThePeople = d == p
		g.PlayerBoards[p].HeroOfThePeople = m == p
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
