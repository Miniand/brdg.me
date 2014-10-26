package starship_catan

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type FightCommand struct{}

func (c FightCommand) Parse(input string) []string {
	return command.ParseNamedCommand("fight", input)
}

func (c FightCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanFight(p)
}

func (c FightCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.Fight(p)
}

func (c FightCommand) Usage(player string, context interface{}) string {
	return "{{b}}fight{{_b}} to fight the pirate"
}

func (g *Game) CanFight(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight ||
		g.FlightCards.Len() == 0 || g.LosingModule {
		return false
	}
	card, _ := g.FlightCards.Pop()
	_, ok := card.(PirateCard)
	return ok
}

func (g *Game) Fight(player int) error {
	var c card.Card
	if !g.CanFight(player) {
		return errors.New("you are unable to fight")
	}
	c, _ = g.FlightCards.Pop()
	pirateCard, ok := c.(PirateCard)
	if !ok {
		return errors.New("card isn't a pirate card")
	}

	pirateRoll := (r.Int() % 3) + 1
	pirateAttack := pirateRoll + pirateCard.Strength
	playerRoll := (r.Int() % 3) + 1
	playerCannon := g.PlayerBoards[player].Resources[ResourceCannon]
	playerAttack := playerRoll + playerCannon
	playerWon := playerAttack >= pirateAttack

	cells := [][]interface{}{
		[]interface{}{
			"",
			"{{b}}Str.{{_b}}",
			"{{b}}Roll{{_b}}",
			"{{b}}Attack{{_b}}",
		},
		[]interface{}{
			g.RenderName(player),
			strconv.Itoa(playerCannon),
			strconv.Itoa(playerRoll),
			Bold(strconv.Itoa(playerAttack)),
		},
		[]interface{}{
			`{{c "gray"}}{{b}}pirate{{_b}}{{_c}}`,
			strconv.Itoa(pirateCard.Strength),
			strconv.Itoa(pirateRoll),
			Bold(strconv.Itoa(pirateAttack)),
		},
	}
	table := render.Table(cells, 0, 2)

	var resultStr string
	if playerWon {
		resultStr = fmt.Sprintf(
			`%s has defeated the pirate`,
			g.RenderName(player),
		)
	} else {
		resultStr = fmt.Sprintf(
			`The pirate has defeated %s`,
			g.RenderName(player),
		)
	}

	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s is fighting the pirate\n%s\n%s",
		g.RenderName(player),
		table,
		resultStr,
	)))

	if playerWon {
		c, g.FlightCards = g.FlightCards.Pop()
		g.PlayerBoards[player].DefeatedPirates =
			g.PlayerBoards[player].DefeatedPirates.Push(c)
		g.RecalculatePeopleCards()
		if err := g.ReplaceCard(); err != nil {
			return err
		}
	} else {
		if pirateCard.DestroyCannon {
			if g.PlayerBoards[player].Resources[ResourceCannon] > 0 {
				g.PlayerBoards[player].Resources[ResourceCannon] -= 1
				g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
					`%s had a cannon destroyed by the pirate`,
					g.RenderName(player),
				)))
			}
		}
		if true || pirateCard.DestroyModule &&
			len(g.PlayerBoards[player].ModuleList()) > 0 {
			g.LosingModule = true
		} else {
			return g.EndFlight()
		}
	}
	return nil
}
