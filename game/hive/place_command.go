package hive

import "github.com/Miniand/brdg.me/command"

type PlaceCommand struct{}

func (pc PlaceCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("place", 3, input)
}

func (pc PlaceCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return player == g.Players[g.CurrentPlayer] && !g.IsFinished()
}

func (pc PlaceCommand) Call(player string, context interface{},
	args []string) (string, error) {
	//a := command.ExtractNamedCommandArgs(args)
	//typeString := a[0]
	//g := context.(*Game)
	return "", nil
}

func (pc PlaceCommand) Usage(player string, context interface{}) string {
	return "{{b}}place TYPE # #{{_b}} to place a new tile on the board, eg. {{b}}place spdr -2 3{{_b}}"
}
