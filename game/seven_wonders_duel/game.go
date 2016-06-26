package seven_wonders_duel

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

const (
	PhaseChooseWonder = iota
	PhasePlay
)

type Game struct {
	Players []string
	Log     *log.Log

	Phase            int
	Finished         bool
	WinnerPlayerNums []int
	Age              int
	Layout           Layout
	CurrentPlayer    int
	Military         int

	ProgressTokens          []int
	DiscardedProgressTokens []int

	RemainingWonders []int
	PlayerWonders    [2][]int

	PlayerCoins    [2]int
	PlayerCards    [2][]int
	DiscardedCards []int
}

func (g *Game) Commands(player string) []command.Command {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return []command.Command{}
	}
	commands := []command.Command{}
	if g.CanChooseWonder(pNum) {
		commands = append(commands, ChooseWonderCommand{})
	}
	if g.CanPlay(pNum) {
		commands = append(commands, PlayCommand{})
	}
	if g.CanDiscard(pNum) {
		commands = append(commands, DiscardCommand{})
	}
	if g.CanWonder(pNum) {
		commands = append(commands, WonderCommand{})
	}
	return commands
}

func (g *Game) Name() string {
	return "7 Wonders: Duel"
}

func (g *Game) Identifier() string {
	return "7_wonders_duel"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if len(players) != 2 {
		return errors.New("2 players only")
	}
	g.Players = players
	g.Log = log.New()
	g.PlayerCoins = [2]int{7, 7}
	g.PlayerCards = [2][]int{{}, {}}
	g.DiscardedCards = []int{}
	g.PlayerWonders = [2][]int{{}, {}}
	progressTokens := helper.IntShuffle(ProgressTokens())
	g.ProgressTokens = progressTokens[:5]
	g.DiscardedProgressTokens = progressTokens[5:]
	g.Phase = PhaseChooseWonder
	g.RemainingWonders = helper.IntShuffle(Wonders())[:8]
	g.Age = 1
	g.Layout = Layout{}
	return nil
}

func (g *Game) StartAge() {
	g.Phase = PhasePlay
	g.Layout = AgeLayouts[g.Age]()
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Finished
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winners := []string{}
	for _, w := range g.WinnerPlayerNums {
		winners = append(winners, g.Players[w])
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	whoseTurn := []string{}
	for _, p := range g.Players {
		if len(g.Commands(p)) > 0 {
			whoseTurn = append(whoseTurn, p)
		}
	}
	return whoseTurn
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) ModifyCoins(player, amount int) {
	if amount == 0 {
		return
	}
	verb := "gained"
	logAmount := amount
	if amount < 0 {
		if g.PlayerCoins[player]-amount < 0 {
			amount = g.PlayerCoins[player]
		}
		verb = "lost"
		logAmount = -amount
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s %s %d %s",
		g.PlayerName(player),
		verb,
		logAmount,
		helper.Plural(logAmount, "coin"),
	)))
	g.PlayerCoins[player] += amount
}

func (g *Game) PlayerCardTypeCount(player, cardType int) int {
	num := 0
	for _, c := range g.PlayerCards[player] {
		if Cards[c].Type == cardType {
			num++
		}
	}
	return num
}

func (g *Game) GreatestCardCount(cardTypes ...int) int {
	num := 0
	for p := range g.Players {
		pNum := 0
		for _, ct := range cardTypes {
			pNum += g.PlayerCardTypeCount(p, ct)
		}
		if pNum > num {
			num = pNum
		}
	}
	return num
}

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}

func (g *Game) PlayerVP(player int) int {
	sum := g.PlayerCoins[player] / 3
	for _, c := range g.PlayerCards[player] {
		sum += Cards[c].VP(g, player)
	}
	return sum
}

func (g *Game) PlayerGoodCount(player, good int) (base, extra int) {
	for _, c := range g.PlayerCards[player] {
		card := Cards[c]
		if card.Type == CardTypeRaw || card.Type == CardTypeManufactured {
			for _, p := range card.Provides {
				base += p[good]
			}
		} else {
			for _, p := range card.Provides {
				extra += p[good]
			}
		}
	}
	return
}

func (g *Game) CheapenedGoods(player int) []int {
	cheapened := cost.Cost{}
	for _, c := range g.PlayerCards[player] {
		crd := Cards[c]
		if crd.Cheapens != nil {
			cheapened = cheapened.Add(cost.FromInts(crd.Cheapens))
		}
	}
	return cheapened.Keys()
}

func BaseTradePrices() cost.Cost {
	return cost.FromInts(Goods).Mul(2)
}

func (g *Game) TradePrices(player int) cost.Cost {
	costs := BaseTradePrices()
	for g, num := range g.TradeGoodCount(Opponent(player)) {
		costs[g] += num
	}
	for _, g := range g.CheapenedGoods(player) {
		costs[g] = 1
	}
	return costs
}

func (g *Game) TradeGoodCount(player int) cost.Cost {
	counts := cost.Cost{}
	for _, c := range g.PlayerCards[player] {
		crd := Cards[c]
		switch crd.Type {
		case CardTypeRaw, CardTypeManufactured:
			for _, p := range crd.Provides {
				counts = counts.Add(p)
			}
		}
	}
	return counts
}

func (g *Game) PlayerProvides(player int) [][]cost.Cost {
	provides := [][]cost.Cost{}
	for _, c := range g.PlayerCards[player] {
		crd := Cards[c]
		if crd.Provides != nil {
			provides = append(provides, crd.Provides)
		}
	}
	return provides
}

func (g *Game) TradeCost(player int, crd Card) int {
	tradePrices := g.TradePrices(player)
	remainingCost := crd.Cost.Take(Goods...).Clone()
	remainingProvides := g.PlayerProvides(player)

	removeGoods := 0
	if crd.Type == CardTypeWonder && g.HasCard(player, ProgressArchitecture) ||
		crd.Type == CardTypeCivilian && g.HasCard(player, ProgressMasonry) {
		removeGoods = 2
	}

	// First pass, no decisions required
	changed := true
	for changed && len(remainingProvides) > 0 && !remainingCost.IsZero() {
		changed = false
		deferred := [][]cost.Cost{}
		for _, p := range remainingProvides {
			if len(p) == 1 {
				// There's no branches so just reduce the cost
				remainingCost, _ = remainingCost.Sub(p[0]).PosNeg()
				changed = true
			} else {
				for _, c := range p {
					subProv := []cost.Cost{}
					matchingKeys := helper.IntIntersect(c.Keys(), remainingCost.Keys())
					if len(matchingKeys) > 0 {
						subProv = append(subProv, c.Take(matchingKeys...))
					} else {
						// We've eliminated a provider branch
						changed = true
					}
				}
			}
		}
		remainingProvides = deferred
	}
	if remainingCost.Sum() <= removeGoods {
		// Can afford
		return 0
	}

	excesses := [][]int{}
	if len(remainingProvides) > 0 {
		// Get permutations
		for _, c := range cost.Perm(remainingProvides) {
			left, _ := remainingCost.Sub(c).PosNeg()
			if left.Sum() <= removeGoods {
				// Can afford
				return 0
			}
			excesses = append(excesses, CostToTradePrices(left, tradePrices))
		}
	} else {
		excesses = append(excesses, CostToTradePrices(remainingCost, tradePrices))
	}

	// Find the cheapest excess
	cheapest := 0
	for i, e := range excesses {
		sum := helper.IntSum(helper.IntReverse(helper.IntSort(e))[removeGoods:])
		if i == 0 || sum < cheapest {
			cheapest = sum
		}
	}
	return cheapest
}

func CostToTradePrices(c cost.Cost, tradePrices cost.Cost) []int {
	trades := []int{}
	for good, num := range c {
		trades = append(trades, helper.IntRepeat(tradePrices[good], num)...)
	}
	return trades
}

func (g *Game) HasCard(player, card int) bool {
	_, ok := helper.IntFind(card, g.PlayerCards[player])
	return ok
}

func (g *Game) WinnersByVP() []int {
	p1VP := g.PlayerVP(0)
	p2VP := g.PlayerVP(1)
	winners := []int{}
	if p1VP >= p2VP {
		winners = append(winners, 0)
	}
	if p2VP >= p1VP {
		winners = append(winners, 1)
	}
	return winners
}

func (g *Game) NextAge() {
	if g.Age == 3 {
		g.Finished = true
		g.WinnerPlayerNums = g.WinnersByVP()
		return
	}
	g.Age += 1
	g.StartAge()
}

func (g *Game) Build(player, card int) error {
	c := Cards[card]
	opp := Opponent(player)
	canFreeBuild := false
	for _, pc := range g.PlayerCards[player] {
		if Cards[pc].MakesFree == card {
			canFreeBuild = true
			break
		}
	}
	if !canFreeBuild && c.Type != CardTypeProgress {
		tradeCost := g.TradeCost(player, c)
		totalCost := c.Cost[GoodCoin] + tradeCost
		if totalCost > g.PlayerCoins[player] {
			return errors.New("you don't have enough coins")
		}
		g.ModifyCoins(player, -totalCost)
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			`%s built
%s`,
			g.PlayerName(player),
			c.RenderMultiline(tradeCost),
		)))
		if g.HasCard(opp, ProgressEconomy) {
			g.ModifyCoins(opp, tradeCost)
		}
	}
	g.PlayerCards[player] = append(g.PlayerCards[player], card)
	// Check scientific victory, 6 different science types
	if c.Type == CardTypeScientific || card == ProgressLaw {
		sciences := map[int]bool{}
		for _, pc := range g.PlayerCards[player] {
			pcc := Cards[pc]
			if pcc.Science != 0 {
				sciences[pcc.Science] = true
				if len(sciences) == 6 {
					// Scientific victory
					g.Finished = true
					g.WinnerPlayerNums = []int{player}
					return nil
				}
			}
		}
	}
	// Check military victory
	if !c.ExtraTurn {
		g.CurrentPlayer = opp
	}
	return nil
}

func (g *Game) RemoveFromLayout(loc Loc) {
	g.Layout = g.Layout.Remove(loc)
	if len(g.Layout) == 0 {
		g.NextAge()
	}
}

func (g *Game) PlayerSciences(player int) map[int]int {
	sciences := map[int]int{}
	for _, c := range g.PlayerCards[player] {
		crd := Cards[c]
		if crd.Science != 0 {
			sciences[crd.Science]++
		}
	}
	return sciences
}

func Opponent(player int) int {
	return (player + 1) % 2
}
