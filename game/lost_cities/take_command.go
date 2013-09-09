package lost_cities

import (
	"errors"
	"github.com/Miniand/brdg.me/command"
	//"strconv"
)

type TakeCommand struct{}

func (d TakeCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("take", 1, input)
}

func (d TakeCommand) CanCall(player string, context interface{}) bool {
	//g := context.(*Game)
	//return g.Players[g.CurrentlyMoving] == player &&
	//	g.TurnPhase == TURN_PHASE_PLAY_OR_DISCARD && !g.IsFinished()
	return true
}

func (d TakeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g := context.(*Game)
	suitnum := 0
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 1 {
		return "", errors.New("You must specify a type of card to take, such as r")
	}
	playerNum, err := g.PlayerFromString(player)
	if err != nil {
		return "", err
	} 
	suit:=args[1]
	switch suit {
	case " r":
		suitnum = SUIT_RED
	case " y":
		suitnum = SUIT_YELLOW
	case " b":
		suitnum = SUIT_BLUE
	case " w":
		suitnum = SUIT_WHITE
	case " g":
		suitnum = SUIT_GREEN
	default:
		return "", errors.New("Could not parse suit")
	}  
	//i, err := strconv.Atoi(args[1])
    //if err != nil {
	//	return "", err
    //}
    //fmt.Println(s, i)
	return "", g.TakeCard(playerNum, suitnum)
}

func (d TakeCommand) Usage(player string, context interface{}) string {
	return "{{b}}take #{{_b}} to take a card from a discard pile, eg. {{b}}take r{{_b}}"
}
