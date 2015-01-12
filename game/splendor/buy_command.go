package splendor

import (
	"errors"

	"github.com/Miniand/brdg.me/command"
)

type BuyCommand struct{}

func (c BuyCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("buy", 1, input)
}

func (c BuyCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	return err != nil && g.CanBuy(pNum)
}

func (c BuyCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	pNum, err := g.PlayerNum(player)
	if err != nil {
		return "", err
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("you must specify which card")
	}
	row, col, err := ParseLoc(a[0])
	if err != nil {
		return "", err
	}
	return "", g.Buy(pNum, row, col)
}

func (c BuyCommand) Usage(player string, context interface{}) string {
	return "{{b}}buy ##{{_b}} to buy a card from the board or your reserve, eg. {{b}}buy 2B{{_b}}"
}

func (g *Game) CanBuy(player int) bool {
	return g.CurrentPlayer == player && g.Phase == PhaseMain
}

func (g *Game) Buy(player, row, col int) error {
	if !g.CanBuy(player) {
		return errors.New("unable to buy right now")
	}
	pb := g.PlayerBoards[player]
	switch row {
	case 0, 1, 2:
		if col < 0 || col > len(g.Board[row]) {
			return errors.New("that is not a valid card")
		}
		if !pb.CanAfford(g.Board[row][col].Cost) {
			return errors.New("you can't afford that card")
		}
		g.PlayerBoards[player].Pay(g.Board[row][col].Cost)
		g.PlayerBoards[player].Cards = append(
			g.PlayerBoards[player].Cards,
			g.Board[row][col],
		)
		if len(g.Decks[row]) > 0 {
			g.Board[row][col] = g.Decks[row][0]
			g.Decks[row] = g.Decks[row][1:]
		} else {
			g.Board[row] = append(
				g.Board[row][:col],
				g.Board[row][col+1:]...,
			)
		}
	case 3:
		if col < 0 || col > len(pb.Reserve) {
			return errors.New("that is not a valid reserve card")
		}
		if !pb.CanAfford(pb.Reserve[col].Cost) {
			return errors.New("you can't afford that card")
		}
		g.PlayerBoards[player].Pay(pb.Reserve[col].Cost)
		g.PlayerBoards[player].Cards = append(
			g.PlayerBoards[player].Cards,
			pb.Reserve[col],
		)
		g.PlayerBoards[player].Reserve = append(
			g.PlayerBoards[player].Reserve[:col],
			g.PlayerBoards[player].Reserve[col+1:]...,
		)
	default:
		return errors.New("that is not a valid row")
	}
	g.NextPhase()
	return nil
}
