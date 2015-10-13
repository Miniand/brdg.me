package seven_wonders_duel

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
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

	Phase         int
	Age           int
	Layout        Layout
	CurrentPlayer int

	ProgressTokens          []int
	DiscardedProgressTokens []int

	RemainingWonders []int
	PlayerWonders    [2][]int

	PlayerCoins [2]int
	PlayerCards [2][]int
}

func (g *Game) Commands(player string) []command.Command {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return []command.Command{}
	}
	commands := []command.Command{}
	if g.CanChoose(pNum) {
		commands = append(commands, ChooseCommand{})
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
	return false
}

func (g *Game) Winners() []string {
	return []string{}
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

func Opponent(player int) int {
	return (player + 1) % 2
}
