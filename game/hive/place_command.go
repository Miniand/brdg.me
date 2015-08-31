package hive

import "github.com/Miniand/brdg.me/command"

type PlaceCommand struct{}

func (pc PlaceCommand) Name() string { return "place" }

func (pc PlaceCommand) Call(
	player string,
	context interface{},
	input *command.Reader,
) (string, error) {
	//a := command.ExtractNamedCommandArgs(args)
	//typeString := a[0]
	//g := context.(*Game)
	return "", nil
}

func (pc PlaceCommand) Usage(player string, context interface{}) string {
	return "{{b}}place TYPE # #{{_b}} to place a new tile on the board, eg. {{b}}place spdr -2 3{{_b}}"
}

func (g *Game) CanPlace(player int) bool {
	return player == g.CurrentPlayer && !g.IsFinished()
}
