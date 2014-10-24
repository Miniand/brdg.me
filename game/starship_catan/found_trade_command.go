package starship_catan

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/log"
)

type FoundTradeCommand struct{}

func (c FoundTradeCommand) Parse(input string) []string {
	return command.ParseNamedCommand("found", input)
}

func (c FoundTradeCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		panic(err)
	}
	return g.CanFoundTradingPost(p)
}

func (c FoundTradeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	p, err := g.ParsePlayer(player)
	if err != nil {
		return "", err
	}
	return "", g.FoundTradingPost(p)
}

func (c FoundTradeCommand) Usage(player string, context interface{}) string {
	return "{{b}}found{{_b}} to found a trading post here"
}

func (g *Game) CanFoundTradingPost(player int) bool {
	if g.CurrentPlayer != player || g.Phase != PhaseFlight ||
		len(g.FlightCards) == 0 || g.TradeAmount != 0 ||
		g.PlayerBoards[player].Resources[ResourceTradeShip] == 0 {
		return false
	}
	c, _ := g.FlightCards.Pop()
	tp, ok := c.(TradingPoster)
	return ok && tp.CanFoundTradingPost()
}

func (g *Game) FoundTradingPost(player int) error {
	var c card.Card

	if !g.CanFoundTradingPost(player) {
		return errors.New("you are not able to found a trading post")
	}
	c, g.FlightCards = g.FlightCards.Pop()
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		`%s founded a trading post on %s`, g.RenderName(player), c)))
	g.PlayerBoards[player].TradingPosts = g.PlayerBoards[player].TradingPosts.Push(c)
	g.PlayerBoards[player].Resources[ResourceTradeShip] -= 1
	g.ReplaceCard()
	g.MarkCardActioned()
	return nil
}
