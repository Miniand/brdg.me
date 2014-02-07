package liars_dice

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/die"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"strings"
)

type CallCommand struct{}

func (c CallCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("call", 0, input)
}

func (c CallCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return !g.IsFinished() && g.WhoseTurn()[0] == player && g.BidQuantity != 0
}

func (c CallCommand) Call(player string, context interface{}, args []string) (
	output string, err error) {
	var (
		resultText   string
		losingPlayer int
	)
	g := context.(*Game)
	quantity := 0
	for _, pd := range g.PlayerDice {
		for _, d := range pd {
			if d == g.BidValue || d == 1 {
				quantity++
			}
		}
	}
	bidPlayerName := render.PlayerNameInPlayers(g.Players[g.BidPlayer],
		g.Players)
	callPlayerName := render.PlayerNameInPlayers(g.Players[g.CurrentPlayer],
		g.Players)
	if quantity < g.BidQuantity {
		// Caller was correct
		losingPlayer = g.BidPlayer
		resultText = fmt.Sprintf("%s bid too high and lost a die",
			bidPlayerName)
	} else {
		// Bidder was correct
		losingPlayer = g.CurrentPlayer
		resultText = fmt.Sprintf("%s bid correctly and %s lost a die",
			bidPlayerName, callPlayerName)
	}
	cells := [][]string{}
	for _, pNum := range g.ActivePlayers() {
		renderedPlayerDice := []string{}
		for _, d := range g.PlayerDice[pNum] {
			renderedPlayerDie := die.Render(d)
			if d == g.BidValue || d == 1 {
				renderedPlayerDie = fmt.Sprintf(`{{c "red"}}%s{{_c}}`,
					renderedPlayerDie)
			}
			renderedPlayerDice = append(renderedPlayerDice, renderedPlayerDie)
		}
		cells = append(cells, []string{
			render.PlayerNameInPlayers(g.Players[pNum], g.Players),
			fmt.Sprintf(`{{l}}%s{{_l}}`, strings.Join(renderedPlayerDice, " ")),
		})
	}
	g.PlayerDice[losingPlayer] = g.PlayerDice[losingPlayer][1:]
	table, err := render.Table(cells, 0, 1)
	if err != nil {
		return "", err
	}
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(`%s called the bid of %d %s by %s
Everyone revealed the following dice:
%s
%s`, callPlayerName, g.BidQuantity, die.Render(g.BidValue), bidPlayerName,
		table, resultText)))
	if !g.IsFinished() {
		g.StartRound()
		g.CurrentPlayer = g.NextActivePlayer(g.CurrentPlayer)
	}
	return
}

func (c CallCommand) Usage(player string, context interface{}) string {
	return "{{b}}call{{_b}} to call the last bidder if you think their bid is too high."
}
