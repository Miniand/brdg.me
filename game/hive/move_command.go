package hive

import "github.com/Miniand/brdg.me/command"

type MoveCommand struct{}

func (mc MoveCommand) Parse(input string) []string {
	return command.ParseNamedCommandNArgs("move", 4, input)
}

func (mc MoveCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	return player == g.Players[g.CurrentPlayer] && !g.IsFinished() &&
		!g.IsOpeningPlay()
}

func (mc MoveCommand) Call(player string, context interface{},
	args []string) (string, error) {
	//a := command.ExtractNamedCommandArgs(args)
	//typeString := a[0]
	//g := context.(*Game)
	return "", nil
}

func (mc MoveCommand) Usage(player string, context interface{}) string {
	return "{{b}}move # # # #{{_b}} to move a tile from one location to another, eg. {{b}}move 1 1 -2 3{{_b}} to move the tile at 1 1 to -2 3"
}
