package love_letter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players         []string
	Log             *log.Log
	Round           int
	Deck, Removed   []int
	Hands, Discards [][]int
	Points          []int
	CurrentPlayer   int
	Eliminated      []bool
	Protected       []bool
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		PlayCommand{},
	}
}

func (g *Game) Name() string {
	return "Love Letter"
}

func (g *Game) Identifier() string {
	return "love_letter"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	l := len(players)
	if l < 2 || l > 4 {
		return errors.New("only for 2 to 4 players")
	}
	g.Players = players
	g.Log = log.New()
	g.Points = make([]int, l)
	g.StartRound()
	return nil
}

func (g *Game) StartRound() {
	g.Round++
	l := len(g.Players)
	g.Eliminated = make([]bool, l)
	g.Protected = make([]bool, l)
	deck := helper.IntShuffle(Deck)
	remove := 1
	if l == 2 {
		remove = 4
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"Starting round %d, {{b}}removing %d %s{{_b}}",
		g.Round,
		remove,
		helper.Plural(remove, "card"),
	)))
	g.Deck, g.Removed = deck[remove:], deck[:remove]
	g.Hands = make([][]int, l)
	g.Discards = make([][]int, l)
	for p := range g.Players {
		g.Hands[p] = []int{}
		g.Discards[p] = []int{}
		g.DrawCard(p)
	}
	g.StartTurn()
}

func (g *Game) StartTurn() {
	g.Protected[g.CurrentPlayer] = false
	if len(g.Deck) == 0 {
		g.EndRound()
	} else {
		g.DrawCard(g.CurrentPlayer)
	}
}

func (g *Game) NextPlayer() {
	for {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		if !g.Eliminated[g.CurrentPlayer] {
			break
		}
	}
	g.StartTurn()
}

func (g *Game) DiscardCardLog(player, card int) {
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s discarded %s",
		g.RenderName(player),
		RenderCard(card),
	)))
	g.DiscardCard(player, card)
}

func (g *Game) DiscardCard(player, card int) {
	g.Hands[player] = helper.IntRemove(card, g.Hands[player], 1)
	g.Discards[player] = append(g.Discards[player], card)
	if card == Princess {
		g.Eliminate(player)
	}
}

func (g *Game) Eliminate(player int) {
	if g.Eliminated[player] {
		return
	}

	g.Eliminated[player] = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s has been eliminated from this round",
		g.RenderName(player),
	)))
	for len(g.Hands[player]) > 0 {
		g.DiscardCardLog(player, g.Hands[player][0])
	}

	numRemaining := 0
	for p := range g.Players {
		if !g.Eliminated[p] {
			numRemaining++
		}
	}
	if numRemaining <= 1 {
		g.EndRound()
	}
}

func (g *Game) EndRound() {
	output := []string{render.Bold("It is the end of the round")}
	var highestCard, highestPlayer, discardTotal int
	for p := range g.Players {
		if g.Eliminated[p] {
			continue
		}
		c := g.Hands[p][0]
		discarded := helper.IntSum(g.Discards[p])
		output = append(output, fmt.Sprintf(
			"%s had %s (total {{b}}%d{{_b}} discarded)",
			g.RenderName(p),
			RenderCard(c),
			discarded,
		))
		if c > highestCard {
			highestCard = c
			discardTotal = 0
		}
		if c == highestCard {
			if discarded > discardTotal {
				discardTotal = discarded
				highestPlayer = p
			}
		}
	}

	g.Points[highestPlayer]++
	output = append(output, fmt.Sprintf(
		"%s won the round and moved to {{b}}%d %s{{_b}}",
		g.RenderName(highestPlayer),
		g.Points[highestPlayer],
		helper.Plural(g.Points[highestPlayer], "point"),
	))

	isFinished := g.IsFinished()
	if isFinished {
		output = append(output, render.Bold(fmt.Sprintf(
			"It is the end of the game, the winner is %s",
			g.RenderName(g.Leader()),
		)))
	}
	g.Log.Add(log.NewPublicMessage(strings.Join(output, "\n")))
	if !isFinished {
		g.StartRound()
	}
}

func (g *Game) Leader() int {
	var highest, player int
	for p := range g.Players {
		points := g.Points[p]
		if points > highest {
			player = p
			highest = points
		}
	}
	return player
}

func (g *Game) DrawCard(player int) {
	var card int
	if len(g.Deck) > 0 {
		card, g.Deck = g.Deck[0], g.Deck[1:]
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s drew a card from the draw pile, {{b}}%d{{_b}} remaining",
			g.RenderName(player),
			len(g.Deck),
		)))
	} else {
		card = g.Removed[0]
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s drew a card from the removed cards",
			g.RenderName(player),
		)))
	}
	g.Log.Add(log.NewPrivateMessage(fmt.Sprintf(
		"You drew %s",
		RenderCard(card),
	), []string{g.Players[player]}))
	g.Hands[player] = append(g.Hands[player], card)
}

func (g *Game) PlayerList() []string {
	return g.Players
}

var endScores = map[int]int{
	2: 7,
	3: 5,
	4: 4,
}

func (g *Game) IsFinished() bool {
	return helper.IntMax(g.Points...) >= endScores[len(g.Players)]
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	return []string{g.Players[g.Leader()]}
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) PlayerNum(player string) (int, bool) {
	return helper.StringInStrings(player, g.Players)
}

func (g *Game) AvailableTargets(forPlayer int) []int {
	targets := []int{}
	for p := range g.Players {
		if p != forPlayer && !g.Eliminated[p] && !g.Protected[p] {
			targets = append(targets, p)
		}
	}
	return targets
}

func (g *Game) ParseTarget(player int, incSelf bool, args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("please specify a player name, if everyone else is protected by the Handmaid you must specify yourself")
	}

	target, err := helper.MatchStringInStrings(args[0], g.Players)
	if err != nil {
		return 0, err
	}

	targets := g.AvailableTargets(player)
	if len(targets) == 0 {
		if target == player {
			return target, nil
		}
		return 0, errors.New("all other players are protected by the Handmaid, so you must specify yourself")
	}

	if !incSelf && target == player {
		return 0, errors.New("you cannot specify yourself if there are other players you can target")
	}

	if g.Eliminated[target] {
		return 0, errors.New("that player is eliminated")
	}
	if g.Eliminated[target] {
		return 0, errors.New("that player is protected by the Handmaid")
	}

	return target, nil
}
