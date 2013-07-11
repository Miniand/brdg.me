package no_thanks

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	Players         []string
	PlayerHands     map[string][]int
	PlayerChips     map[string]int
	CentreChips     int
	RemainingCards  []int
	CurrentlyMoving string
}

func (g *Game) PlayerAction(player, action string, params []string) error {
	var err error
	switch strings.ToLower(action) {
	case "take":
		err = g.Take(player)
	case "pass":
		err = g.Pass(player)
	}
	return err
}

func (g *Game) Name() string {
	return "No Thanks"
}

func (g *Game) Identifier() string {
	return "no_thanks"
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
		if player == g.CurrentlyMoving {
			buf.WriteString(
				"It's your turn, you can {{b}}take{{_b}} or {{b}}pass{{_b}} the card.\n\n")
		}
		buf.WriteString(
			`{{b}}Current card:  {{c "blue"}}{{.PeekTopCard}}{{_c}}{{_b}} (`)
		buf.WriteString(strconv.Itoa(len(g.RemainingCards) - 1))
		buf.WriteString(" remaining)\n")
		buf.WriteString(
			`{{b}}Current chips: {{c "green"}}{{.CentreChips}}{{_c}}{{_b}}`)
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
	longestPlayerName := 0
	for _, p := range g.Players {
		if len(p) > longestPlayerName {
			longestPlayerName = len(p)
		}
	}
	buf.WriteString("{{b}}Players{{_b}}\n\n")
	for _, p := range g.Players {
		buf.WriteString(`{{b}}`)
		buf.WriteString(p)
		buf.WriteString(":{{_b}}")
		buf.WriteString(strings.Repeat(" ", longestPlayerName-len(p)+1))
		if len(g.PlayerHands[p]) > 0 {
			buf.WriteString(`{{c "blue"}}`)
			buf.WriteString(g.RenderCardsForPlayer(p, g.PeekTopCard()))
			buf.WriteString("{{_c}}")
		} else {
			buf.WriteString(`{{c "gray"}}no cards{{_c}}`)
		}
		if g.IsFinished() {
			buf.WriteString(`     ({{c "green"}}`)
			buf.WriteString(strconv.Itoa(g.PlayerChips[p]))
			buf.WriteString(`{{_c}} chips, {{c "magenta"}}`)
			buf.WriteString(strconv.Itoa(g.FinalPlayerScore(p)))
			buf.WriteString("{{_c}} points)")
		}
		buf.WriteString("\n")
	}
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
	if len(players) < 2 || len(players) > 5 {
		return errors.New("No Thanks requires between 2 and 5 players")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Players = players
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
	cardPool := g.AllCards()
	picked := map[int]bool{}
	g.RemainingCards = make([]int, 24)
	for i := 0; i < 24; i++ {
		c := cardPool[r.Int()%24]
		for picked[c] {
			c = cardPool[r.Int()%24]
		}
		picked[c] = true
		g.RemainingCards[i] = c
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
	return g.NextPlayer()
}

func (g *Game) Take(player string) error {
	err := g.AssertTurn(player)
	if err != nil {
		return err
	}
	g.PlayerHands[player] = append(g.PlayerHands[player], g.PopTopCard())
	g.PlayerChips[player] += g.CentreChips
	g.CentreChips = 0
	return g.NextPlayer()
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
