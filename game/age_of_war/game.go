package age_of_war

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Game struct {
	Players       []string
	CurrentPlayer int
	Log           *log.Log

	Conquered    map[int]bool
	CastleOwners map[int]int

	CurrentlyAttacking int
	CompletedLines     map[int]bool
	CurrentRoll        []int
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		AttackCommand{},
		LineCommand{},
		RollCommand{},
	}
}

func (g *Game) Name() string {
	return "Age of War"
}

func (g *Game) Identifier() string {
	return "age_of_war"
}

func (g *Game) Encode() ([]byte, error) {
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	if l := len(players); l < 2 || l > 6 {
		return errors.New("only for 2 to 6 players")
	}
	g.Players = players
	g.Log = log.New()

	g.Conquered = map[int]bool{}
	g.CastleOwners = map[int]int{}
	g.CompletedLines = map[int]bool{}

	g.StartTurn()

	return nil
}

func (g *Game) StartTurn() {
	g.CurrentlyAttacking = -1
	g.CompletedLines = map[int]bool{}
	g.Roll(7)
}

func (g *Game) NextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	g.StartTurn()
}

func (g *Game) CheckEndOfTurn() bool {
	if g.CurrentlyAttacking != -1 {
		c := Castles[g.CurrentlyAttacking]
		lines := c.CalcLines(
			g.Conquered[g.CurrentlyAttacking],
		)
		// If the player has completed all lines, they take the card and it is
		// the end of the turn.
		allLines := true
		for l := range lines {
			if !g.CompletedLines[l] {
				allLines = false
				break
			}
		}
		if allLines {
			suffix := ""
			if g.Conquered[g.CurrentlyAttacking] {
				suffix = fmt.Sprintf(
					" from %s",
					g.PlayerName(g.CastleOwners[g.CurrentlyAttacking]),
				)
			}
			g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
				"%s conquered the castle %s%s",
				g.PlayerName(g.CurrentPlayer),
				c.RenderName(),
				suffix,
			)))
			g.Conquered[g.CurrentlyAttacking] = true
			g.CastleOwners[g.CurrentlyAttacking] = g.CurrentPlayer
			if clanConquered, _ := g.ClanConquered(c.Clan); clanConquered {
				g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
					"%s conquered the clan %s",
					g.PlayerName(g.CurrentPlayer),
					RenderClan(c.Clan),
				)))
			}
			g.NextTurn()
			return true
		}

		// If the player doesn't have enough dice to complete the rest of the
		// lines, it is the end of the turn.
		reqDice := 0
		numDice := len(g.CurrentRoll)
		canAffordLine := false
		for i, l := range lines {
			if g.CompletedLines[i] {
				continue
			}
			reqDice += l.MinDice()
			if reqDice > numDice {
				g.FailedAttackMessage()
				g.NextTurn()
				return true
			}
			if can, _ := l.CanAfford(g.CurrentRoll); can {
				canAffordLine = true
			}
		}

		// If the player has the minimum required dice but they can't afford a
		// line, it is the end of the turn.
		if reqDice == numDice && !canAffordLine {
			g.FailedAttackMessage()
			g.NextTurn()
			return true
		}
	} else {
		// If the player doesn't have enough dice for any castle, it is the end
		// of the turn.
		for i, c := range Castles {
			if g.Conquered[i] && g.CastleOwners[i] == g.CurrentPlayer {
				// They already own it.
				continue
			}
			if conquered, _ := g.ClanConquered(c.Clan); conquered {
				// The clan is conquered, can't steal.
				continue
			}
			minDice := c.MinDice()
			if g.Conquered[i] {
				minDice++
			}
			if minDice <= len(g.CurrentRoll) {
				// They can afford this one
				return false
			}
		}
		// They couldn't afford anything, next turn.
		g.FailedAttackMessage()
		g.NextTurn()
		return true
	}
	return false
}

func (g *Game) FailedAttackMessage() {
	target := "anything"
	if g.CurrentlyAttacking != -1 {
		target = Castles[g.CurrentlyAttacking].RenderName()
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s failed to conquer %s",
		g.PlayerName(g.CurrentPlayer),
		target,
	)))
}

func (g *Game) Scores() map[int]int {
	scores := map[int]int{}
	conqueredClans := map[int]bool{}
	for cIndex, c := range Castles {
		if !g.Conquered[cIndex] {
			continue
		}
		clanConquered, ok := conqueredClans[c.Clan]
		if !ok {
			var conqueredBy int
			clanConquered, conqueredBy = g.ClanConquered(c.Clan)
			conqueredClans[c.Clan] = clanConquered
			if clanConquered {
				scores[conqueredBy] += ClanSetPoints[c.Clan]
			}
		}
		if clanConquered {
			continue
		}
		scores[g.CastleOwners[cIndex]] += c.Points
	}
	return scores
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) IsFinished() bool {
	return len(g.Conquered) == len(Castles)
}

func (g *Game) Winners() []string {
	if !g.IsFinished() {
		return []string{}
	}
	// Winner is determined by score, with ties broken by conquered clans.
	playerConqueredClans := map[int]int{}
	for _, clan := range Clans {
		if conquered, by := g.ClanConquered(clan); conquered {
			playerConqueredClans[by]++
		}
	}
	maxScore := 0
	winners := []string{}
	for p, s := range g.Scores() {
		score := s*10 + playerConqueredClans[p]
		if p == 0 || score > maxScore {
			maxScore = score
			winners = []string{}
		}
		if score == maxScore {
			winners = append(winners, g.Players[p])
		}
	}
	return winners
}

func (g *Game) WhoseTurn() []string {
	return []string{g.Players[g.CurrentPlayer]}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) ClanConquered(clan int) (conquered bool, player int) {
	player = -1
	conquered = true
	for i, c := range Castles {
		if c.Clan != clan {
			continue
		}
		if !g.Conquered[i] {
			conquered = false
			return
		}
		if player == -1 {
			player = g.CastleOwners[i]
		} else if player != g.CastleOwners[i] {
			conquered = false
			return
		}
	}
	return
}

func (g *Game) PlayerNum(player string) (int, bool) {
	for p, name := range g.Players {
		if name == player {
			return p, true
		}
	}
	return 0, false
}
