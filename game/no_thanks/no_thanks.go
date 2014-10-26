package no_thanks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type Game struct {
	Players         []string
	PlayerHands     map[string][]int
	PlayerChips     map[string]int
	CentreChips     int
	RemainingCards  []int
	CurrentlyMoving string
	Log             *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		TakeCommand{},
		PassCommand{},
	}
}

func (g *Game) Name() string {
	return "No Thanks"
}

func (g *Game) Identifier() string {
	return "no_thanks"
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) Encode() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Game) Decode(data []byte) error {
	return json.Unmarshal(data, g)
}

func (g *Game) RenderForPlayer(player string) (string, error) {
	buf := bytes.NewBufferString("")
	if !g.IsFinished() {
		buf.WriteString(`{{b}}Current card:  {{c "blue"}}`)
		buf.WriteString(strconv.Itoa(g.PeekTopCard()))
		buf.WriteString(`{{_c}}{{_b}} (`)
		buf.WriteString(strconv.Itoa(len(g.RemainingCards) - 1))
		buf.WriteString(" cards remaining)\n")
		buf.WriteString(`{{b}}Current chips: {{c "green"}}`)
		buf.WriteString(strconv.Itoa(g.CentreChips))
		buf.WriteString(`{{_c}}{{_b}}`)
		buf.WriteString("\n\n")
		buf.WriteString(`{{b}}Your hand:{{_b}}  `)
		if len(g.PlayerHands[player]) > 0 {
			buf.WriteString(`{{c "blue"}}`)
			buf.WriteString(g.RenderCardsForPlayer(player, g.PeekTopCard()))
			buf.WriteString("{{_c}}")
		} else {
			buf.WriteString(`{{c "gray"}}no cards{{_c}}`)
		}
		buf.WriteString("\n")
		buf.WriteString(`{{b}}Your chips:{{_b}} {{c "green"}}`)
		buf.WriteString(strconv.Itoa(g.PlayerChips[player]))
		buf.WriteString("{{_c}}\n\n")
	}
	header := []string{"{{b}}Players{{_b}}", "{{b}}Cards{{_b}}"}
	if g.IsFinished() {
		header = append(header, "{{b}}Score{{_b}}")
	}
	cells := [][]string{
		header,
	}
	for pNum, p := range g.Players {
		row := []string{
			render.PlayerName(pNum, p),
		}
		if len(g.PlayerHands[p]) > 0 {
			row = append(row, fmt.Sprintf(`{{c "blue"}}%s{{_c}}`,
				g.RenderCardsForPlayer(p, g.PeekTopCard())))
		} else {
			row = append(row, `{{c "gray"}}no cards{{_c}}`)
		}
		if g.IsFinished() {
			row = append(row, fmt.Sprintf(
				`{{b}}{{c "green"}}%d{{_c}}{{_b}} chips, {{b}}{{c "magenta"}}%d{{_c}}{{_b}} points`,
				g.PlayerChips[p], g.FinalPlayerScore(p)))
		}
		cells = append(cells, row)
	}
	table := render.Table(cells, 0, 2)
	buf.WriteString(table)
	return buf.String(), nil
}

func (g *Game) RenderCardsForPlayer(player string, relevant int) string {
	renderGroups := []string{}
	for _, group := range g.PlayerHandGrouped(player) {
		renderGroup := []string{}
		for _, c := range group {
			if c-relevant == 1 || c-relevant == -1 {
				renderGroup = append(renderGroup,
					"{{b}}"+strconv.Itoa(c)+"{{_b}}")
			} else {
				renderGroup = append(renderGroup, strconv.Itoa(c))
			}
		}
		renderGroups = append(renderGroups, strings.Join(renderGroup, " "))
	}
	return strings.Join(renderGroups, "   ")
}

func (g *Game) Start(players []string) error {
	if len(players) < 3 || len(players) > 5 {
		return errors.New("No Thanks requires between 3 and 5 players")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Players = players
	g.Log = log.New()
	g.InitCards()
	g.InitPlayerChips()
	g.InitPlayerHands()
	g.CurrentlyMoving = g.Players[r.Int()%len(g.Players)]
	return nil
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return len(g.RemainingCards) == 0
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	winners := []string{}
	winningScore := 0
	for _, p := range g.Players {
		pScore := g.FinalPlayerScore(p)
		if len(winners) == 0 || winningScore > pScore {
			winners = []string{p}
			winningScore = pScore
		} else if pScore == winningScore {
			winners = append(winners, p)
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	return []string{g.CurrentlyMoving}
}

func (g *Game) AllCards() []int {
	cards := make([]int, 33)
	for i := 3; i <= 35; i++ {
		cards[i-3] = i
	}
	return cards
}

func (g *Game) InitCards() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(33)
	cardPool := g.AllCards()
	g.RemainingCards = make([]int, 24)
	for i := 0; i < 24; i++ {
		g.RemainingCards[i] = cardPool[perm[i]]
	}
}

func (g *Game) InitPlayerChips() {
	g.PlayerChips = map[string]int{}
	for _, p := range g.Players {
		g.PlayerChips[p] = 11
	}
}

func (g *Game) InitPlayerHands() {
	g.PlayerHands = map[string][]int{}
	for _, p := range g.Players {
		g.PlayerHands[p] = []int{}
	}
}

func (g *Game) AssertTurn(player string) error {
	if g.IsFinished() {
		return errors.New("The game has already finished")
	}
	if g.CurrentlyMoving != player {
		return errors.New("It's not your turn")
	}
	return nil
}

func (g *Game) Pass(player string) error {
	err := g.AssertTurn(player)
	if err != nil {
		return err
	}
	if g.PlayerChips[player] <= 0 {
		return errors.New("You have no chips left, you must take the card")
	}
	g.PlayerChips[player]--
	g.CentreChips++
	pName := render.PlayerNameInPlayers(player, g.Players)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s passed on the {{b}}{{c "blue"}}%d{{_c}}{{_b}}`,
		pName, g.PeekTopCard())))
	return g.NextPlayer()
}

func (g *Game) Take(player string) error {
	err := g.AssertTurn(player)
	if err != nil {
		return err
	}
	pName := render.PlayerNameInPlayers(player, g.Players)
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s took the {{b}}{{c "blue"}}%d{{_c}}{{_b}} and {{b}}{{c "green"}}%d{{_c}}{{_b}} chips`,
		pName, g.PeekTopCard(), g.CentreChips)))
	g.PlayerHands[player] = append(g.PlayerHands[player], g.PopTopCard())
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s drew {{b}}{{c "blue"}}%d{{_c}}{{_b}} as the new card`,
		pName, g.PeekTopCard())))
	g.PlayerChips[player] += g.CentreChips
	g.CentreChips = 0
	return nil
}

func (g *Game) PeekTopCard() int {
	if len(g.RemainingCards) == 0 {
		return 0
	}
	return g.RemainingCards[len(g.RemainingCards)-1]
}

func (g *Game) PopTopCard() int {
	top := g.PeekTopCard()
	g.RemainingCards = g.RemainingCards[:len(g.RemainingCards)-1]
	return top
}

func (g *Game) NextPlayer() error {
	// Find the index of the current player
	playerIndex := 0
	playerFound := false
	for i, p := range g.Players {
		if p == g.CurrentlyMoving {
			playerIndex = i
			playerFound = true
			break
		}
	}
	if !playerFound {
		return errors.New(
			"Could not find the current player in the player list")
	}
	g.CurrentlyMoving = g.Players[(playerIndex+1)%len(g.Players)]
	return nil
}

func (g *Game) PlayerHandSorted(player string) []int {
	sort.Ints(g.PlayerHands[player])
	return g.PlayerHands[player]
}

func (g *Game) PlayerHandGrouped(player string) [][]int {
	groups := [][]int{}
	curGroup := []int{}
	lastCard := -1
	for _, c := range g.PlayerHandSorted(player) {
		if c == lastCard+1 {
			curGroup = append(curGroup, c)
		} else {
			if len(curGroup) > 0 {
				groups = append(groups, curGroup)
			}
			curGroup = []int{c}
		}
		lastCard = c
	}
	if len(curGroup) > 0 {
		groups = append(groups, curGroup)
	}
	return groups
}

func (g *Game) PlayerHandScore(player string) int {
	score := 0
	for _, g := range g.PlayerHandGrouped(player) {
		score += g[0]
	}
	return score
}

func (g *Game) FinalPlayerScore(player string) int {
	return g.PlayerHandScore(player) - g.PlayerChips[player]
}

func (g *Game) BotPlay(player string) error {
	if g.CurrentlyMoving != player {
		return errors.New("Not their turn")
	}
	// If they're out of money, they gotta take it
	if g.PlayerChips[player] == 0 {
		_, err := command.CallInCommands(player, g, "take", g.Commands())
		return err
	}
	// Decide how much we want the card
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Int()%2 == 0 {
		_, err := command.CallInCommands(player, g, "take", g.Commands())
		return err
	}
	_, err := command.CallInCommands(player, g, "pass", g.Commands())
	return err
}
