package cathedral

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
)

type PlayCommand struct{}

func (c PlayCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("play", 2, 3, input)
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
	pieceNum-- // Change to zero index
	x, y, ok := ParseLoc(a[1])
	if !ok {
		return "", errors.New("the second argument should be a valid location, such as C7")
	}
	dir := DirDown
	if len(a) > 2 {
		dir, err = helper.MatchStringInStringMap(a[2], OrthoDirNames)
		if err != nil {
			return "", err
		}
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
	if piece < 0 || piece > len(Pieces[player]) {
		return errors.New("that is not a valid piece number")
	}
	if g.PlayedPieces[player][piece] {
		return errors.New("you have already played that piece")
	}
	p := Pieces[player][piece]
	// Special case for player 2, if they haven't played the cathedral they
	// need to play it first.
	if player == 1 && piece != 0 && !g.PlayedPieces[1][0] {
		return errors.New("cathedral piece must be played before any others")
	}
	n := 0
	switch dir {
	case DirUp:
		n = 2
	case DirRight:
		n = -1
	case DirLeft:
		n = 1
	}
	rotated := p.Positions.Rotate(n)
	// First ensure it can actually be played.
	for _, l := range rotated {
		tX := x + l.X
		tY := y + l.Y
		if tX < 0 || tX > 9 || tY < 0 || tY > 9 {
			return errors.New("playing there would go off the board")
		}
		t := g.Board[tY][tX]
		if t.Player != NoPlayer {
			return errors.New("there is already a piece there")
		}
		if t.Owner != NoPlayer &&
			t.Owner != player {
			return errors.New("the other player owns that area")
		}
	}
	for _, l := range rotated {
		g.Board[y+l.Y][x+l.X].Player = p.Player
		g.Board[y+l.Y][x+l.X].Type = p.Type
	}
	g.PlayedPieces[player][piece] = true
	// Do an ownership check.
	if p.Player != PlayerCathedral && g.PlayedPieces[1][0] {
		log.Print("OWNERSHIP CHECK")
	}
	if player != 1 || piece != 0 {
		// Go to next player if it wasn't the cathedral just played.
		g.NextPlayer()
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
	y = int(strings.ToUpper(matches[1])[0] - 'A')
	x, _ = strconv.Atoi(matches[2])
	x--
	if x < 0 || x > 9 {
		ok = false
	}
	return
}
