package sushi_go

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

const Dummy = 2

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players    []string
	AllPlayers []string // Includes dummy if needed

	Round int

	Deck    []int
	Hands   [][]int
	Playing map[int][]int

	Played map[int][]int
	Points map[int]int

	Controller int // For 2 players, who is controlling the dummy this turn

	Log *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
		DummyCommand{},
	}
}

func (g *Game) Name() string {
	return "Sushi Go"
}

func (g *Game) Identifier() string {
	return "sushi_go"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	_, ok := PlayerDrawCounts[len(players)]
	if !ok {
		return errors.New("requires between 2 and 5 players")
	}

	g.Log = log.New()
	g.Players = players
	g.AllPlayers = append([]string{}, players...)
	if len(players) == 2 {
		g.AllPlayers = append(g.AllPlayers, "Dummy")
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"Because there are only two players, you will be joined by %s",
			g.RenderName(Dummy),
		)))
	}

	g.Deck = Shuffle(Deck())
	g.Playing = map[int][]int{}
	g.Played = map[int][]int{}
	g.Points = map[int]int{}
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.Round++
	for p := range g.AllPlayers {
		// Remove anything that's not a pudding
		newPlayed := []int{}
		for _, c := range g.Played[p] {
			if c == CardPudding {
				newPlayed = append(newPlayed, c)
			}
		}
		g.Played[p] = newPlayed
	}
	g.Hands = make([][]int, len(g.AllPlayers))
	drawCount := PlayerDrawCounts[len(g.Players)]
	passDir := "left"
	if g.Round == 2 {
		passDir = "right"
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Starting round {{b}}%d{{_b}}, hands will be passed to the {{b}}%s{{_b}}.  Dealing {{b}}%d{{_b}} cards to each player",
		g.Round,
		passDir,
		drawCount,
	)))
	for p := range g.AllPlayers {
		g.Hands[p], g.Deck = g.Deck[:drawCount], g.Deck[drawCount:]
		g.Hands[p] = Sort(g.Hands[p])
	}
	g.StartHand()
}

func (g *Game) StartHand() {
	if len(g.Players) == 2 {
		// Controller draws a card from the dummy hand.
		i := rnd.Int() % len(g.Hands[Dummy])
		g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
			"You drew %s from %s",
			RenderCard(g.Hands[Dummy][i]),
			g.RenderName(Dummy),
		), []string{g.Players[g.Controller]}))
		g.Hands[g.Controller] = append(g.Hands[g.Controller], g.Hands[Dummy][i])
		g.Hands[g.Controller] = Sort(g.Hands[g.Controller])
		g.Hands[Dummy] = append(g.Hands[Dummy][:i], g.Hands[Dummy][i+1:]...)
	}
}

func (g *Game) EndHand() {
	// Play cards
	for p := range g.AllPlayers {
		g.Hands[p] = TrimPlayed(g.Hands[p])
		g.Played[p] = append(g.Played[p], g.Playing[p]...)
		if len(g.Playing[p]) == 2 {
			// Use chopsticks.
			if i, ok := Contains(CardChopsticks, g.Played[p]); ok {
				g.Hands[p] = append(g.Hands[p], CardChopsticks)
				g.Played[p] = append(g.Played[p][:i], g.Played[p][i+1:]...)
			}
		}
		g.Playing[p] = nil
	}
	if len(g.Players) == 2 {
		// Next player controls the dummy
		g.Controller = (g.Controller + 1) % 2
	}
	// End round if we're out of cards
	if len(g.Hands[0]) == 0 {
		g.EndRound()
		return
	}
	// Pass hands
	if len(g.Players) == 2 {
		g.Log.Add(log.NewPublicMessage("Players are swapping hands"))
		g.Hands[0], g.Hands[1] = g.Hands[1], g.Hands[0]
	} else if g.Round%2 == 1 {
		g.Log.Add(log.NewPublicMessage("Passing hands to the {{b}}left{{_b}}"))
		extra := g.Hands[0]
		g.Hands = append(g.Hands[1:], extra)
	} else {
		g.Log.Add(log.NewPublicMessage("Passing hands to the {{b}}right{{_b}}"))
		l := len(g.Hands)
		extra := g.Hands[l-1]
		g.Hands = append([][]int{extra}, g.Hands[:l-1]...)
	}
	g.StartHand()
}

func (g *Game) Score() ([]int, []string) {
	scores := make([]int, len(g.AllPlayers))
	output := []string{}

	// Score maki
	maki := map[int]int{}
	for p := range g.AllPlayers {
		for _, c := range g.Played[p] {
			switch c {
			case CardMakiRoll1:
				maki[p]++
			case CardMakiRoll2:
				maki[p] += 2
			case CardMakiRoll3:
				maki[p] += 3
			}
		}
	}
	first := 0
	firstPlayers := []int{}
	second := 0
	secondPlayers := []int{}
	for p, m := range maki {
		if m > first {
			second = first
			secondPlayers = firstPlayers
			first = m
			firstPlayers = []int{}
		}
		if m == first {
			firstPlayers = append(firstPlayers, p)
		} else if m == second {
			secondPlayers = append(secondPlayers, p)
		}
	}
	makiRollsStr := render.Colour("maki rolls", CardColours[CardMakiRoll1])
	if first == 0 {
		output = append(output, fmt.Sprintf(
			"Nobody had %s, no points awarded",
			makiRollsStr,
		))
	} else {
		firstPoints := 6 / len(firstPlayers)
		output = append(output, fmt.Sprintf(
			"%s had {{b}}%d{{_b}} %s, awarding {{b}}%d points{{_b}}",
			render.CommaList(g.RenderNames(firstPlayers)),
			first,
			makiRollsStr,
			firstPoints,
		))
		for _, p := range firstPlayers {
			scores[p] += firstPoints
		}
		if len(firstPlayers) == 1 && second > 0 && len(secondPlayers) <= 3 {
			secondPoints := 3 / len(secondPlayers)
			output = append(output, fmt.Sprintf(
				"%s had {{b}}%d{{_b}} %s, awarding {{b}}%d points{{_b}}",
				render.CommaList(g.RenderNames(secondPlayers)),
				second,
				makiRollsStr,
				secondPoints,
			))
			for _, p := range secondPlayers {
				scores[p] += secondPoints
			}
		}
	}

	if g.Round == 3 {
		// Score puddings
		pudding := map[int]int{}
		for p := range g.AllPlayers {
			for _, c := range g.Played[p] {
				if c == CardPudding {
					pudding[p]++
				}
			}
		}
		first := 0
		firstPlayers := []int{}
		last := 0
		lastPlayers := []int{}
		for p := range g.AllPlayers {
			c := pudding[p]
			if c > first {
				first = c
				firstPlayers = []int{}
			}
			if c == first {
				firstPlayers = append(firstPlayers, p)
			}
			if c < last || len(lastPlayers) == 0 {
				last = c
				lastPlayers = []int{}
			}
			if c == last {
				lastPlayers = append(lastPlayers, p)
			}
		}
		puddingsStr := render.Colour("puddings", CardColours[CardPudding])
		if first == last {
			output = append(output, fmt.Sprintf(
				"Everybody had the same number of %s, no points awarded",
				puddingsStr,
			))
		} else {
			firstPoints := 6 / len(firstPlayers)
			output = append(output, fmt.Sprintf(
				"%s had {{b}}%d{{_b}} %s, awarding {{b}}%d points{{_b}}",
				render.CommaList(g.RenderNames(firstPlayers)),
				first,
				puddingsStr,
				firstPoints,
			))
			for _, p := range firstPlayers {
				scores[p] += firstPoints
			}
			if len(g.Players) != 2 {
				lastPoints := -6 / len(lastPlayers)
				output = append(output, fmt.Sprintf(
					"%s had {{b}}%d{{_b}} %s, awarding {{b}}%d points{{_b}}",
					render.CommaList(g.RenderNames(lastPlayers)),
					last,
					puddingsStr,
					lastPoints,
				))
				for _, p := range lastPlayers {
					scores[p] += lastPoints
				}
			}
		}
	}

	// Score normal cards
	for p := range g.AllPlayers {
		output = append(output, fmt.Sprintf(
			render.Bold("Scoring cards for %s"),
			g.RenderName(p),
		))
		cardCounts := map[int]int{}
		for _, c := range g.Played[p] {
			if s, ok := CardBaseScores[c]; ok {
				text := RenderCard(c)
				if cardCounts[CardWasabi] > 0 {
					s *= 3
					cardCounts[CardWasabi]--
					text = fmt.Sprintf("%s + %s", text, RenderCard(CardWasabi))
				}
				output = append(output, fmt.Sprintf(
					"%s, {{b}}%d{{_b}} points",
					text,
					s,
				))
				scores[p] += s
			} else {
				cardCounts[c]++
			}
		}
		if s := cardCounts[CardTempura] / 2 * 5; s > 0 {
			output = append(output, fmt.Sprintf(
				"%d x %s, {{b}}%d{{_b}} points",
				cardCounts[CardTempura],
				RenderCard(CardTempura),
				s,
			))
			scores[p] += s
		}
		if s := cardCounts[CardSashimi] / 3 * 10; s > 0 {
			output = append(output, fmt.Sprintf(
				"%d x %s, {{b}}%d{{_b}} points",
				cardCounts[CardSashimi],
				RenderCard(CardSashimi),
				s,
			))
			scores[p] += s
		}
		if n := cardCounts[CardDumpling]; n > 0 {
			s := (n*n + n) / 2
			if s > 15 {
				s = 15
			}
			output = append(output, fmt.Sprintf(
				"%d x %s, {{b}}%d{{_b}} points",
				cardCounts[CardDumpling],
				RenderCard(CardDumpling),
				s,
			))
			scores[p] += s
		}
	}
	return scores, output
}

func (g *Game) EndRound() {
	scores, output := g.Score()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"{{b}}It is the end of round %d, scoring{{_b}}\n%s",
		g.Round,
		output,
	)))
	for p := range g.AllPlayers {
		g.Points[p] += scores[p]
	}
	if g.Round < 3 {
		g.StartRound()
	}
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return g.Round == 3 && len(g.Hands[0]) == 0 && g.Playing[0] == nil
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	return []string{}
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	whose := []string{}
	commands := g.Commands()
	for _, pName := range g.Players {
		if len(command.AvailableCommands(pName, g, commands)) > 0 {
			whose = append(whose, pName)
		}
	}
	return whose
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for p, pName := range g.Players {
		if pName == player {
			return p, true
		}
	}
	return 0, false
}

func (g *Game) RenderName(player int) string {
	return render.PlayerName(player, g.AllPlayers[player])
}

func (g *Game) RenderNames(players []int) []string {
	playerStrs := make([]string, len(players))
	for i, p := range players {
		playerStrs[i] = g.RenderName(p)
	}
	return playerStrs
}
