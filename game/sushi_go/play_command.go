package sushi_go

import (
	"errors"
	"strconv"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("play", 1, 2, input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return false
	}
	return g.CanPlay(pNum)
}

func (c PlayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("it is not your turn at the moment")
	}

	a := command.ExtractNamedCommandArgs(args)
	cards := make([]int, len(a))
	for i := range a {
		card, err := strconv.Atoi(a[i])
		if err != nil {
			return "", errors.New("each card must be a number")
		}
		cards[i] = card - 1 // Input is 1 based
	}

	return "", g.Play(pNum, cards)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play # (#){{_b}} to play one or two cards, eg. {{b}}play 3{{_b}}.  To play two cards you must already have played chopsticks"
}

func (g *Game) CanPlay(player int) bool {
	return g.Playing[player] == nil
}

func (g *Game) Play(player int, cards []int) error {
	if !g.CanPlay(player) {
		return errors.New("you can't play at the moment")
	}
	l := len(cards)
	if l == 0 || l > 2 {
		return errors.New("you must specify one or two cards to play")
	}
	if l == 2 {
		if _, ok := Contains(CardChopsticks, g.Played[player]); !ok {
			return errors.New("you can only play a second card if you've previously played chopsticks")
		}
		if player == g.Controller && g.Playing[Dummy] == nil &&
			len(g.Players) == 2 && len(g.Hands[player]) == 2 {
			// Need to keep room for the dummy player.
			return errors.New("you can't play two cards now, you have to save one for the dummy player")
		}
	}

	return g.PlayCards(player, player, cards)
}

func (g *Game) PlayCards(toPlayer, fromPlayer int, cards []int) error {
	cardMap := map[int]bool{}
	for _, c := range cards {
		if c < 0 || c >= len(g.Hands[fromPlayer]) {
			return errors.New("that card number is not valid")
		}
		if cardMap[c] {
			return errors.New("please specify different cards")
		}
		if g.Hands[fromPlayer][c] == CardPlayed {
			return errors.New("that card has already been played")
		}
		cardMap[c] = true
	}

	// Valid, do that thing
	g.Playing[toPlayer] = make([]int, len(cards))
	for i, c := range cards {
		g.Playing[toPlayer][i] = g.Hands[fromPlayer][c]
		g.Hands[fromPlayer][c] = CardPlayed
	}

	// Check if everyone has played cards
	for p := range g.AllPlayers {
		if g.Playing[p] == nil {
			return nil
		}
	}
	g.EndHand()
	return nil
}
