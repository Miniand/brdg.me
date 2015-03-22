package cathedral

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
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
	if g.NoOpenTiles {
		// Both players play simultaneously.
		return g.CanPlaySomething(player)
	}
	return g.CurrentPlayer == player
}

func (g *Game) CanPlayPiece(player, piece int, loc Loc, dir int) (bool, string) {
	if piece < 0 || piece > len(Pieces[player]) {
		return false, "that is not a valid piece number"
	}
	if g.PlayedPieces[player][piece] {
		return false, "you have already played that piece"
	}
	p := Pieces[player][piece]
	// Special case for player 2, if they haven't played the cathedral they
	// need to play it first.
	if player == 1 && piece != 0 && !g.PlayedPieces[1][0] {
		return false, "cathedral piece must be played before any others"
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
			return false, "playing there would go off the board"
		}
		t := g.Board[l]
		if t.Player != NoPlayer {
			return false, "there is already a piece there"
		}
		if t.Owner != NoPlayer &&
			t.Owner != player {
			return false, "the other player owns that area"
		}
	}
	return true, ""
}

func (g *Game) Play(player, piece int, loc Loc, dir int) error {
	if !g.CanPlay(player) {
		return errors.New("can't make plays at the moment")
	}
	if ok, reason := g.CanPlayPiece(player, piece, loc, dir); !ok {
		return errors.New(reason)
	}
	p := Pieces[player][piece]
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
	for _, l := range rotated {
		l = l.Add(loc)
		t := g.Board[l]
		t.Player = p.Player
		t.Type = p.Type
		g.Board[l] = t
	}
	g.PlayedPieces[player][piece] = true
	g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
		"%s played {{b}}%d{{_b}} (size {{b}}%d{{_b}}) {{b}}%s{{_b}} from {{b}}%s{{_b}}",
		g.PlayerName(player),
		p.Type,
		len(p.Positions),
		OrthoDirNames[dir],
		loc,
	)))
	// Do an ownership check.
	if p.Player != PlayerCathedral && g.PlayedPieces[1][1] {
		g.CheckCaptures(loc)
	}
	// Check if there are any open tiles left, otherwise it becomes
	// simultaneous play.
	if !g.NoOpenTiles {
		openTileExists := false
		for _, l := range AllLocs {
			if t := g.Board[l]; t.Player == NoPlayer && t.Owner == NoPlayer {
				openTileExists = true
				break
			}
		}
		if !openTileExists {
			g.NoOpenTiles = true
			g.Log.Add(log.NewPublicMessage(
				"No open tiles remain, players will play the rest of their pieces simultaneously.",
			))
		}
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
	capturedTileCount := 0
	capturedPieceCount := 0
	capturedPieceSize := 0
	Walk(loc, OrthoDirs, func(l Loc) int {
		if visited[l] {
			return WalkBlocked
		}
		if g.Board[l].Owner == player {
			// Player already owns it so we don't need to keep walking here.
			visited[l] = true
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
		Walk(l, Dirs, func(l Loc) int {
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
			capturedTileCount += len(area)
			for pt := range pieces {
				if pt.Player != PlayerCathedral {
					capturedPieceCount++
					g.PlayedPieces[pt.Player][pt.Type-1] = false
				}
			}
			for _, areaLoc := range area {
				if g.Board[areaLoc].Player != NoPlayer &&
					g.Board[areaLoc].Player != PlayerCathedral {
					capturedPieceSize++
				}
				t := EmptyTile
				t.Owner = player
				g.Board[areaLoc] = t
			}
		}
		return WalkContinue
	})
	if capturedTileCount > 0 {
		suffix := ""
		if capturedPieceCount > 0 {
			suffix = fmt.Sprintf(
				" and returned {{b}}%d{{_b}} pieces with a combined size of {{b}}%d{{_b}}",
				capturedPieceCount,
				capturedPieceSize,
			)
		}
		g.Log.Add(log.NewPublicMessage(fmt.Sprintf(
			"%s captured an area of {{b}}%d{{_b}}%s",
			g.PlayerName(player),
			capturedTileCount,
			suffix,
		)))
	}
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
