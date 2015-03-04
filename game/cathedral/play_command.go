package cathedral

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("play", 1, 2, input)
}

func (c PlayCommand) CanCall(player string, context interface{}) bool {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	return ok && g.CanPlay(pNum)
}

func (c PlayCommand) Call(
	player string,
	context interface{},
	args []string,
) (output string, err error) {
	g := context.(*Game)
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	a := command.ExtractNamedCommandArgs(args)
	if len(a) < 2 {
		return "", errors.New("the play command requires at least two arguments")
	}
	pieceNum, err := strconv.Atoi(a[0])
	if err != nil {
		return "", errors.New("the first argument should be the piece number to play")
	}
	x, y, ok := ParseLoc(a[1])
	if !ok {
		return "", errors.New("the second argument should be a valid location, such as C7")
	}
	dir := DirDown
	if len(a) > 2 {
	}
	return "", g.Play(pNum, pieceNum, x, y, dir)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play # loc (dir){{_b}} to play a tile in a direction, eg. {{b}}play 1 b5 right{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Play(player, piece, x, y, dir int) error {
	if !g.CanPlay(player) {
		return errors.New("can't make plays at the moment")
	}
	return nil
}

var parseLocRegexp = regexp.MustCompile(`(?i)^([a-j])(\d+)$`)

func ParseLoc(input string) (x, y int, ok bool) {
	matches := parseLocRegexp.FindStringSubmatch(input)
	if matches == nil {
		return
	}
	ok = true
	x = int(strings.ToUpper(matches[1])[0] - 'A')
	y, _ = strconv.Atoi(matches[2])
	if y > 10 {
		ok = false
	}
	return
}
