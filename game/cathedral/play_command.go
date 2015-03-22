package cathedral

import (
	"errors"
	"fmt"
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
	loc, ok := ParseLoc(a[1])
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
	return "", g.Play(pNum, pieceNum, loc, dir)
}

func (c PlayCommand) Usage(player string, context interface{}) string {
	return "{{b}}play # loc (dir){{_b}} to play a tile in a direction, eg. {{b}}play 1 b5 right{{_b}}"
}

func (g *Game) CanPlay(player int) bool {
	return g.CurrentPlayer == player
}

func (g *Game) Play(player, piece int, loc Loc, dir int) error {
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
		l = l.Add(loc)
		if !l.Valid() {
			fmt.Print(l)
			return errors.New("playing there would go off the board")
		}
		t := g.Board[l]
		if t.Player != NoPlayer {
			return errors.New("there is already a piece there")
		}
		if t.Owner != NoPlayer &&
			t.Owner != player {
			return errors.New("the other player owns that area")
		}
	}
	for _, l := range rotated {
		l = l.Add(loc)
		t := g.Board[l]
		t.Player = p.Player
		t.Type = p.Type
		g.Board[l] = t
	}
	g.PlayedPieces[player][piece] = true
	// Do an ownership check.
	if p.Player != PlayerCathedral && g.PlayedPieces[1][1] {
		g.CheckCaptures(loc)
	}
	if player != 1 || piece != 0 {
		// Go to next player if it wasn't the cathedral just played.
		g.NextPlayer()
	}
	return nil
}

func (g *Game) CheckCaptures(loc Loc) {
	player := g.Board[loc].Player
	// Walk to find all adjoining empty regions.
	visited := map[Loc]bool{}
	Walk(loc, OrthoDirs, func(l Loc) int {
		if visited[l] {
			return WalkBlocked
		}
		if g.Board[l].Player == player {
			// Extension of the player pieces, continue.
			visited[l] = true
			return WalkContinue
		}
		// Check for capture.
		area := []Loc{}
		pieces := map[PlayerType]bool{}
		Walk(l, OrthoDirs, func(l Loc) int {
			if visited[l] || g.Board[l].Player == player {
				return WalkBlocked
			}
			visited[l] = true
			area = append(area, l)
			if g.Board[l].Player != NoPlayer {
				pieces[g.Board[l].PlayerType] = true
			}
			return WalkContinue
		})
		if len(pieces) <= 1 {
			// Capture!
			for pt := range pieces {
				if pt.Player != PlayerCathedral {
					g.PlayedPieces[pt.Player][pt.Type] = false
				}
			}
			for _, areaLoc := range area {
				t := EmptyTile
				t.Owner = player
				g.Board[areaLoc] = t
			}
		}
		return WalkContinue
	})
}

var parseLocRegexp = regexp.MustCompile(`(?i)^([a-j])(\d+)$`)

func ParseLoc(input string) (loc Loc, ok bool) {
	matches := parseLocRegexp.FindStringSubmatch(input)
	if matches == nil {
		return
	}
	loc.Y = int(strings.ToUpper(matches[1])[0] - 'A')
	loc.X, _ = strconv.Atoi(matches[2])
	loc.X--
	ok = loc.Valid()
	return
}
