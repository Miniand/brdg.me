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
	Players        []string
	CurrentPlayer  int
	AttackDamage   int
	AttackPlayers  []int
	LeftPlayer     int
	Phase          int
	CurrentRoll    []int
	ExtraRollable  map[int]bool
	RemainingRolls int
	ExtraTurns     []int
	FaceUpCards    []int
	Deck           []int
	Discard        []int
	Boards         []*PlayerBoard
	Tokyo          []int
	Log            *log.Log
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
	diceCount := 6
	for _, t := range g.Boards[g.CurrentPlayer].Things() {
		if dm, ok := t.(DiceCountModifier); ok {
			diceCount = dm.ModifyDiceCount(g, g.CurrentPlayer, diceCount)
		}
	}
	if diceCount > 8 {
		diceCount = 8
	}
	g.RollPhaceNDice(diceCount)
}

func (g *Game) RollPhaceNDice(diceCount int) {
	g.Phase = PhaseRoll
	// 2 VP for being in Tokyo at the start of the turn
	if g.PlayerLocation(g.CurrentPlayer) != LocationOutside {
		g.Boards[g.CurrentPlayer].ModifyVP(2)
	}
	g.CurrentRoll = RollDice(diceCount)
	g.LogRoll(g.CurrentPlayer, g.CurrentRoll, []int{})
	g.RemainingRolls = 2
	for _, t := range g.Boards[g.CurrentPlayer].Things() {
		if mod, ok := t.(RollCountModifier); ok {
			g.RemainingRolls = mod.ModifyRollCount(
				g,
				g.CurrentPlayer,
				g.RemainingRolls,
			)
		}
	}
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
	g.LeftPlayer = TokyoEmpty
	things := g.Boards[g.CurrentPlayer].Things()
	roll := g.CurrentRoll
	for _, t := range things {
		if preResolve, ok := t.(PreResolveDiceHandler); ok {
			roll = preResolve.HandlePreResolveDice(g, g.CurrentPlayer, roll)
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
				g.Boards[g.CurrentPlayer].ModifyVP(d + count - 2)
			}
		case DieEnergy:
			g.ModifyEnergy(g.CurrentPlayer, count)
		case DieAttack:
			if count > 0 {
				isAttacking = true
			}
		case DieHeal:
			if g.PlayerLocation(g.CurrentPlayer) == LocationOutside {
				g.Boards[g.CurrentPlayer].ModifyHealth(count)
			}
		}
	}
	if isAttacking {
		g.AttackPhase(
			attacked,
			diceCounts[DieAttack],
		)
	} else {
		for _, t := range things {
			if postResolve, ok := t.(PostResolveDiceHandler); ok {
				postResolve.HandlePostResolveDice(g, g.CurrentPlayer, roll)
			}
		}
		g.BuyPhase()
	}
}

func (g *Game) ModifyEnergy(player, amount int) {
	for _, t := range g.Boards[player].Things() {
		if mod, ok := t.(EnergyModifier); ok {
			amount = mod.ModifyEnergy(g, player, amount)
		}
	}
	g.Boards[player].ModifyEnergy(amount)
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

func (g *Game) DealDamage(attacker, target, damage, defenderAction int) {
	// First attacker modifies damage
	for _, t := range g.Boards[attacker].Things() {
		if damageMod, ok := t.(AttackDamageForPlayerModifier); ok {
			damage = damageMod.ModifyAttackDamageForPlayer(g, attacker, target, damage)
		}
	}
	// Second attacked modifies damage
	for _, t := range g.Boards[target].Things() {
		if damageMod, ok := t.(DamageModifier); ok {
			damage = damageMod.ModifyDamage(g, target, attacker, damage, defenderAction)
		}
	}
	if damage != 0 {
		g.TakeDamage(target, damage)
		for _, t := range g.Boards[attacker].Things() {
			if handler, ok := t.(DamageDealtHandler); ok {
				handler.HandleDamageDealt(g, attacker, target, damage)
			}
		}
	}
}

func (g *Game) TakeDamage(player, damage int) {
	g.Boards[player].ModifyHealth(-damage)
	if g.Boards[player].Health == 0 {
		// Leave Tokyo if they are in it
		if loc := g.PlayerLocation(player); loc != LocationOutside {
			g.Tokyo[loc] = TokyoEmpty
		}
		for p, _ := range g.Players {
			for _, t := range g.Boards[p].Things() {
				if zero, ok := t.(HealthZeroHandler); ok {
					zero.HandleHealthZero(g, p, player)
				}
			}
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
			for _, t := range g.Boards[g.LeftPlayer].Things() {
				if lt, ok := t.(LeaveTokyoHandler); ok {
					lt.HandleLeaveTokyo(g, l, tp, TokyoEmpty)
				}
			}
		}
	}
	// Enter tokyo if there's room
	for _, p := range g.TokyoLocs() {
		if p == TokyoEmpty {
			to := g.TakeControl(g.CurrentPlayer)
			if g.LeftPlayer != TokyoEmpty {
				for _, t := range g.Boards[g.LeftPlayer].Things() {
					if lt, ok := t.(LeaveTokyoHandler); ok {
						lt.HandleLeaveTokyo(g, to, g.LeftPlayer, g.CurrentPlayer)
					}
				}
			}
			break
		}
	}
	for _, t := range g.Boards[g.CurrentPlayer].Things() {
		if postAttack, ok := t.(PostAttackHandler); ok {
			postAttack.HandlePostAttack(g, g.CurrentPlayer, g.AttackDamage)
		}
	}
	g.NextPhase()
}

func (g *Game) TakeControl(player int) int {
	if loc := g.PlayerLocation(player); loc != LocationOutside {
		// Player is already in Tokyo
		return loc
	}
	// Move into the empty one if there, otherwise to Tokyo City (even if
	// another player is there).
	to := LocationTokyoCity
	for l, p := range g.TokyoLocs() {
		if p == TokyoEmpty {
			to = l
			break
		}
	}
	if g.Tokyo[to] == TokyoEmpty {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained control of %s gaining %s",
			g.RenderName(player),
			LocationStrings[to],
			RenderVP(1),
		)))
	} else {
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s gained control of %s from %s, gaining %s",
			g.RenderName(player),
			LocationStrings[to],
			g.RenderName(g.Tokyo[to]),
			RenderVP(1),
		)))
	}
	g.Tokyo[to] = player
	g.Boards[player].VP += 1
	return to
}

func (g *Game) LeaveTokyo(player int) {
	loc := g.PlayerLocation(player)
	if loc != LocationOutside {
		g.Tokyo[loc] = TokyoEmpty
	}
}

func (g *Game) HandleAttackedPlayer() {
	if len(g.AttackPlayers) == 0 {
		g.EndAttackPhase()
		return
	}
	p := g.AttackPlayers[0]
	damage := g.AttackDamage
	if g.PlayerLocation(p) == LocationOutside {
		g.DealDamage(g.CurrentPlayer, p, damage, DefenderActionNone)
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
	var targets []int
	switch g.PlayerLocation(player) {
	case LocationOutside:
		targets = g.PlayersInsideTokyo()
	default:
		targets = g.PlayersOutsideTokyo()
	}
	for _, t := range g.Boards[player].Things() {
		if atm, ok := t.(AttackTargetModifier); ok {
			targets = atm.ModifyAttackTargets(g, player, targets)
		}
	}
	return targets
}

func RegisterGobTypes() {
	for _, c := range Deck {
		gob.Register(Cards[c])
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
	deck := helper.IntShuffle(Deck)
	g.FaceUpCards = deck[:3]
	g.Deck = deck[3:]
	g.ExtraTurns = []int{}
	g.Discard = []int{}
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
	for _, t := range g.Boards[g.CurrentPlayer].Things() {
		if endTurn, ok := t.(EndTurnHandler); ok {
			endTurn.HandleEndTurn(g, g.CurrentPlayer)
		}
	}
	if !g.IsFinished() {
		if len(g.ExtraTurns) > 0 {
			diceCount := g.ExtraTurns[0]
			g.ExtraTurns = g.ExtraTurns[1:]
			g.RollPhaceNDice(diceCount)
		} else {
			g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
			for g.Boards[g.CurrentPlayer].Health <= 0 {
				g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
			}
			g.RollPhase()
		}
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

func (g *Game) PlayerNum(player string) (int, bool) {
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
