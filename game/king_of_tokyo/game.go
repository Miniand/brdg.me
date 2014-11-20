package king_of_tokyo

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	PhaseRoll = iota
	PhaseAttack
	PhaseBuy
)

const (
	TokyoEmpty        = -1
	LocationOutside   = -1
	LocationTokyoCity = 0
	LocationTokyoBay  = 1
)

type Game struct {
	Players           []string
	CurrentPlayer     int
	AttackDamage      int
	AttackPlayers     []int
	CurrentFleeingLoc int
	Phase             int
	CurrentRoll       []int
	ExtraRollable     map[int]bool
	RemainingRolls    int
	Buyable           []CardBase
	Deck              []CardBase
	Discard           []CardBase
	Boards            []*PlayerBoard
	Tokyo             []int
	Log               *log.Log
}

func (g *Game) Commands() []command.Command {
	return []command.Command{
		RollCommand{},
		KeepCommand{},
		StayCommand{},
		LeaveCommand{},
		BuyCommand{},
		SweepCommand{},
		DoneCommand{},
	}
}

func (g *Game) NextPhase() {
	switch g.Phase {
	case PhaseRoll:
		g.ResolveDice()
	case PhaseAttack:
		g.BuyPhase()
	case PhaseBuy:
		g.NextTurn()
	}
}

func (g *Game) RollPhase() {
	// 2 VP for being in Tokyo at the start of the turn
	if g.PlayerLocation(g.CurrentPlayer) != LocationOutside {
		g.Boards[g.CurrentPlayer].VP += 2
	}
	g.Phase = PhaseRoll
	g.CurrentRoll = RollDice(6)
	g.LogRoll(g.CurrentPlayer, g.CurrentRoll, []int{})
	g.RemainingRolls = 2
}

func (g *Game) CheckRollComplete() {
	// Check if we have more rerolls
	extra := map[int]bool{}
	for _, t := range g.Boards[g.CurrentPlayer].Things() {
		if extraReroller, ok := t.(ExtraReroller); ok {
			extra = extraReroller.ExtraReroll(g, g.CurrentPlayer, extra)
		}
	}
	if len(extra) > 0 {
		g.ExtraRollable = extra
	} else {
		g.NextPhase()
	}
}

func (g *Game) ResolveDice() {
	things := g.Boards[g.CurrentPlayer].Things()
	roll := g.CurrentRoll
	for _, t := range things {
		if preResolve, ok := t.(PreResolveDiceHandler); ok {
			roll = preResolve.PreResolveDice(g, g.CurrentPlayer, roll)
		}
	}
	// Handle dice
	diceCounts := map[int]int{}
	for _, d := range roll {
		diceCounts[d] += 1
	}
	// Modify attack
	attacked := g.AttackTargetsForPlayer(g.CurrentPlayer)
	for _, t := range things {
		if attackMod, ok := t.(AttackModifier); ok {
			diceCounts[DieAttack], attacked = attackMod.ModifyAttack(
				g,
				g.CurrentPlayer,
				diceCounts[DieAttack],
				attacked,
			)
		}
	}
	isAttacking := false
	for _, d := range Dice {
		count := diceCounts[d]
		if count == 0 {
			continue
		}
		switch d {
		case Die1, Die2, Die3:
			if count >= 3 {
				g.Boards[g.CurrentPlayer].VP += d + count - 2
			}
		case DieEnergy:
			g.Boards[g.CurrentPlayer].Energy += count
		case DieAttack:
			if count > 0 {
				isAttacking = true
			}
		case DieHeal:
			if g.PlayerLocation(g.CurrentPlayer) == LocationOutside {
				g.Boards[g.CurrentPlayer].Health += count
				if g.Boards[g.CurrentPlayer].Health > 10 {
					g.Boards[g.CurrentPlayer].Health = 10
				}
			}
		}
	}
	if isAttacking {
		g.AttackPhase(
			attacked,
			diceCounts[DieAttack],
		)
	} else {
		g.BuyPhase()
	}
}

func (g *Game) AttackPhase(players []int, damage int) {
	g.Phase = PhaseAttack
	g.AttackDamage = damage
	// Make sure outside players are attacked first, then Tokyo players in order
	orderedAttackPlayers := []int{}
	inTokyoMap := map[int]bool{}
	for _, p := range players {
		if loc := g.PlayerLocation(p); loc == LocationOutside {
			orderedAttackPlayers = append(orderedAttackPlayers, p)
		} else {
			inTokyoMap[loc] = true
		}
	}
	for l, p := range g.Tokyo {
		if inTokyoMap[l] {
			orderedAttackPlayers = append(orderedAttackPlayers, p)
		}
	}
	g.AttackPlayers = orderedAttackPlayers
	g.HandleAttackedPlayer()
}

func (g *Game) TakeDamage(player, amount int) {
	g.Boards[player].Health -= amount
	if g.Boards[player].Health <= 0 {
		g.Boards[player].Health = 0
		// Leave Tokyo if they are in it
		if loc := g.PlayerLocation(player); loc != LocationOutside {
			g.Tokyo[loc] = TokyoEmpty
		}
	}
}

func (g *Game) NextAttackedPlayer() {
	if len(g.AttackPlayers) > 1 {
		g.AttackPlayers = g.AttackPlayers[1:]
		g.HandleAttackedPlayer()
	} else {
		g.EndAttackPhase()
	}
}

func (g *Game) EndAttackPhase() {
	// Remove others from Tokyo if their location no longer
	// exists.
	for l, tp := range g.Tokyo[g.TokyoSize():] {
		if tp != TokyoEmpty {
			g.Tokyo[l] = TokyoEmpty
		}
	}
	// Enter tokyo if there's room
	for t, p := range g.TokyoLocs() {
		if p == TokyoEmpty {
			g.Tokyo[t] = g.CurrentPlayer
			g.Boards[g.CurrentPlayer].VP += 1
			break
		}
	}
	for _, t := range g.Boards[g.CurrentPlayer].Things() {
		if postAttack, ok := t.(PostAttackHandler); ok {
			postAttack.PostAttack(g, g.CurrentPlayer, g.AttackDamage)
		}
	}
	g.NextPhase()
}

func (g *Game) HandleAttackedPlayer() {
	if len(g.AttackPlayers) == 0 {
		g.EndAttackPhase()
		return
	}
	p := g.AttackPlayers[0]
	damage := g.AttackDamage
	for _, t := range g.Boards[p].Things() {
		if damageMod, ok := t.(DamageModifier); ok {
			damage = damageMod.ModifyDamage(g, p, g.CurrentPlayer, damage)
		}
	}
	loc := g.PlayerLocation(p)
	if damage <= 0 {
		g.NextAttackedPlayer()
	} else if loc == LocationOutside {
		g.TakeDamage(p, damage)
		g.NextAttackedPlayer()
	}
}

func (g *Game) BuyPhase() {
	// Start buy phase
	g.Phase = PhaseBuy
}

func (g *Game) Name() string {
	return "King of Tokyo"
}

func (g *Game) Identifier() string {
	return "king_of_tokyo"
}

func (g *Game) PlayerLocation(player int) int {
	for loc, p := range g.TokyoLocs() {
		if p == player {
			return loc
		}
	}
	return LocationOutside
}

func (g *Game) PlayersInsideTokyo() []int {
	players := []int{}
	for _, p := range g.TokyoLocs() {
		if p != -1 {
			players = append(players, p)
		}
	}
	return players
}

func (g *Game) PlayersOutsideTokyo() []int {
	players := []int{}
	inside := g.PlayersInsideTokyo()
	for p, _ := range g.Players {
		outside := true
		for _, in := range inside {
			if p == in {
				outside = false
				break
			}
		}
		if outside {
			players = append(players, p)
		}
	}
	return players
}

func (g *Game) AttackTargetsForPlayer(player int) []int {
	switch g.PlayerLocation(player) {
	case LocationOutside:
		return g.PlayersInsideTokyo()
	default:
		return g.PlayersOutsideTokyo()
	}
}

func RegisterGobTypes() {
	for _, c := range Deck() {
		gob.Register(c)
	}
}

func (g *Game) Encode() ([]byte, error) {
	RegisterGobTypes()
	return helper.Encode(g)
}

func (g *Game) Decode(data []byte) error {
	RegisterGobTypes()
	return helper.Decode(g, data)
}

func (g *Game) Start(players []string) error {
	playerCount := len(players)
	if playerCount < 2 || playerCount > 6 {
		return errors.New("requires between 2 and 6 players")
	}
	g.Players = players
	g.Log = log.New()
	deck := Shuffle(Deck())
	g.Buyable = deck[:3]
	g.Deck = deck[3:]
	g.Discard = []CardBase{}
	g.Tokyo = make([]int, 2)
	for i, _ := range g.Tokyo {
		g.Tokyo[i] = TokyoEmpty
	}
	g.Boards = make([]*PlayerBoard, playerCount)
	for p, _ := range g.Players {
		g.Boards[p] = NewPlayerBoard()
	}
	g.RollPhase()
	return nil
}

func (g *Game) NextTurn() {
	if !g.IsFinished() {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		for g.Boards[g.CurrentPlayer].Health <= 0 {
			g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
		}
		g.RollPhase()
	}
}

func (g *Game) LogRoll(player int, rolled, kept []int) {
	diceStr := []string{}
	for _, d := range rolled {
		diceStr = append(diceStr, render.Bold(RenderDie(d)))
	}
	for _, d := range kept {
		diceStr = append(diceStr, RenderDie(d))
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s rolled  %s",
		g.RenderName(g.CurrentPlayer),
		strings.Join(diceStr, "  "),
	)))
}

func (g *Game) PlayerList() []string {
	return g.Players
}

func (g *Game) CheckGameEnd() (isFinished bool, winners []string) {
	winners = []string{}
	remainingPlayers := []string{}
	for p, b := range g.Boards {
		if b.Health > 0 {
			if b.VP >= 20 {
				return true, []string{g.Players[p]}
			}
			remainingPlayers = append(remainingPlayers, g.Players[p])
		}
	}
	if len(remainingPlayers) < 2 {
		isFinished = true
		winners = remainingPlayers
	}
	return
}

func (g *Game) IsFinished() bool {
	isFinished, _ := g.CheckGameEnd()
	return isFinished
}

func (g *Game) Winners() []string {
	_, winners := g.CheckGameEnd()
	return winners
}

func (g *Game) WhoseTurn() []string {
	if g.IsFinished() {
		return []string{}
	}
	switch g.Phase {
	case PhaseAttack:
		return []string{g.Players[g.AttackPlayers[0]]}
	default:
		return []string{g.Players[g.CurrentPlayer]}
	}
}

func (g *Game) GameLog() *log.Log {
	return g.Log
}

func (g *Game) EliminatedPlayerList() []string {
	eliminated := []string{}
	for p, b := range g.Boards {
		if b.Health <= 0 {
			eliminated = append(eliminated, g.Players[p])
		}
	}
	return eliminated
}

func (g *Game) PlayerNum(player string) (int, error) {
	return helper.StringInStrings(player, g.Players)
}

func (g *Game) TokyoSize() int {
	if len(g.Players)-len(g.EliminatedPlayerList()) > 4 {
		return 2
	}
	return 1
}

func (g *Game) TokyoLocs() []int {
	return g.Tokyo[:g.TokyoSize()]
}

func ContainsInt(i int, s []int) bool {
	for _, si := range s {
		if si == i {
			return true
		}
	}
	return false
}
