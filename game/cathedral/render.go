package cathedral

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/render"
)

const (
	TileWidth  = 6
	TileHeight = 3
)

var (
	NoTileStr       = ` `
	PieceBackground = `░`
)

var WallStrs = map[int]string{
	DirUp | DirDown | DirLeft | DirRight: render.Bold("╬"),
	DirUp | DirDown | DirLeft:            render.Bold("╣"),
	DirUp | DirDown | DirRight:           render.Bold("╠"),
	DirUp | DirLeft | DirRight:           render.Bold("╩"),
	DirDown | DirLeft | DirRight:         render.Bold("╦"),
	DirUp | DirLeft:                      render.Bold("╝"),
	DirUp | DirRight:                     render.Bold("╚"),
	DirDown | DirLeft:                    render.Bold("╗"),
	DirDown | DirRight:                   render.Bold("╔"),
	DirLeft | DirRight:                   render.Bold("═"),
	DirLeft:                              render.Bold("═"),
	DirRight:                             render.Bold("═"),
	DirUp | DirDown:                      render.Bold("║"),
	DirUp:                                render.Bold("║"),
	DirDown:                              render.Bold("║"),
}

var SideWall = render.Bold(strings.TrimSpace(strings.Repeat(fmt.Sprintf(
	"%s\n",
	WallStrs[DirUp|DirDown],
), TileHeight)))

func (g *Game) RenderForPlayer(player string) (string, error) {
	pNum, ok := g.PlayerNum(player)
	if !ok {
		return "", errors.New("could not find player")
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(g.Board.Render())
	buf.WriteString(`

All pieces are shown in their {{b}}down{{_b}} position and pivot around the number.

{{b}}Your remaining tiles:{{_b}}
`)
	buf.WriteString(g.RenderPlayerRemainingTiles(pNum))
	buf.WriteString(fmt.Sprintf(
		"\n\n{{b}}%s remaining tiles:{{_b}}\n",
		g.PlayerName(Opponent(pNum)),
	))
	buf.WriteString(g.RenderPlayerRemainingTiles(Opponent(pNum)))
	return buf.String(), nil
}

func (g *Game) RenderPlayerRemainingTiles(pNum int) string {
	buf := bytes.NewBuffer([]byte{})
	cells := [][]interface{}{{}}
	curWidth := 0
	hasTiles := false
	for i, p := range Pieces[pNum] {
		if g.PlayedPieces[pNum][i] {
			continue
		}
		hasTiles = true
		pWidth := p.Width()
		if curWidth+pWidth > 10 {
			buf.WriteString("\n")
			buf.WriteString(render.Table(cells, 0, 0))
			cells = [][]interface{}{{}}
			curWidth = 0
		}
		cells[0] = append(cells[0], p.Render())
		curWidth += pWidth
	}
	if !hasTiles {
		return render.Markup("None", render.Gray, true)
	}
	if len(cells) > 0 {
		buf.WriteString("\n")
		buf.WriteString(render.Table(cells, 0, 0))
	}
	return buf.String()
}

var (
	emptyAbove = (TileHeight - 1) / 2
	emptyBelow = TileHeight / 2
)

func RenderTile(src Tiler, loc Loc) (string, bool) {
	t, ok := src.TileAt(loc)
	if !ok || t.Player == NoPlayer {
		return "", false
	}
	return RenderPlayerTile(t, OpenSides(src, loc)), true
}

func RenderPlayerTile(tile Tile, open map[int]bool) string {
	// Top row
	buf := bytes.NewBufferString(RenderCorner(DirUp|DirLeft, open))
	c := WallStrs[DirLeft|DirRight]
	if open[DirUp] {
		c = PieceBackground
	}
	buf.WriteString(strings.Repeat(c, TileWidth-2))
	buf.WriteString(RenderCorner(DirUp|DirRight, open))
	buf.WriteString("\n")

	// Middle rows
	left := WallStrs[DirUp|DirDown]
	if open[DirLeft] {
		left = PieceBackground
	}
	right := WallStrs[DirUp|DirDown]
	if open[DirRight] {
		right = PieceBackground
	}
	remainingWidth := TileWidth - 2 - render.StrLen(tile.Text)
	leftPadding := strings.Repeat(PieceBackground, remainingWidth/2)
	rightPadding := strings.Repeat(PieceBackground, (remainingWidth+1)/2)
	middleRow := fmt.Sprintf(
		"%s%s%s%s%s\n",
		left,
		leftPadding,
		render.Colour(tile.Text, render.Black),
		rightPadding,
		right,
	)
	buf.WriteString(strings.Repeat(middleRow, TileHeight-2))

	// Bottom row
	buf.WriteString(RenderCorner(DirDown|DirLeft, open))
	c = WallStrs[DirLeft|DirRight]
	if open[DirDown] {
		c = PieceBackground
	}
	buf.WriteString(strings.Repeat(c, TileWidth-2))
	buf.WriteString(RenderCorner(DirDown|DirRight, open))

	return render.Markup(buf.String(), render.PlayerColour(tile.Player), true)
}

func RenderCorner(dir int, open map[int]bool) string {
	// If all three tiles in dir are open, then render nothing.
	numOpen := 0
	for _, d := range Dirs {
		if dir&d == d && open[d] {
			numOpen++
			if numOpen == 3 {
				return PieceBackground
			}
		}
	}

	// Map of one corner direction referencing the other.
	cornerMap := map[int]int{}
	first := -1
	for _, d := range Dirs {
		if dir&d != d {
			continue
		}
		if first == -1 {
			first = d
		} else {
			cornerMap[first] = d
			cornerMap[d] = first
			break
		}
	}

	var corner int
	for d, other := range cornerMap {
		if open[d] {
			corner = corner | d
		} else {
			corner = corner | DirInv(other)
		}
	}
	return WallStrs[corner]
}

func RenderEmptyTile(loc Loc, owner int) string {
	colour := render.Gray
	if owner != NoPlayer {
		colour = render.PlayerColour(owner)
	}
	buf := bytes.NewBufferString(strings.Repeat(fmt.Sprintf(
		"%s\n",
		strings.Repeat(NoTileStr, TileWidth),
	), emptyAbove))
	s := loc.String()
	remainingWidth := TileWidth - len(s)
	buf.WriteString(strings.Repeat(NoTileStr, remainingWidth/2))
	buf.WriteString(render.Markup(s, colour, true))
	buf.WriteString(strings.Repeat(NoTileStr, (remainingWidth+1)/2))
	buf.WriteByte('\n')
	buf.WriteString(strings.TrimSpace(strings.Repeat(fmt.Sprintf(
		"%s\n",
		strings.Repeat(NoTileStr, TileWidth),
	), emptyBelow)))
	return render.Colour(buf.String(), render.Gray)
}

func (g *Game) PlayerName(pNum int) string {
	return render.PlayerName(pNum, g.Players[pNum])
}
