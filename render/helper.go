package render

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type Centred string
type RightAligned string
type Unbounded string

type CellSpan struct {
	Content interface{}
	Cols    int
}

var Rule = strings.Repeat("=", 80)

func PlayerColour(playerNum int) string {
	colours := []string{
		"green",
		"red",
		"blue",
		"yellow",
		"cyan",
		"magenta",
		"gray",
		"black",
	}
	return colours[playerNum%len(colours)]
}

func TrimPlayerName(player string) string {
	return fmt.Sprintf("• %.12s", strings.Split(
		strings.TrimSpace(player), "@")[0])
}

func PlayerName(playerNum int, player string) string {
	return fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`, PlayerColour(playerNum),
		TrimPlayerName(player))
}

func PlayerNameInPlayers(player string, players []string) string {
	for playerNum, p := range players {
		if player == p {
			return PlayerName(playerNum, player)
		}
	}
	return player
}

func PlayerNamesInPlayers(players []string, playerList []string) []string {
	renderedPlayers := make([]string, len(players))
	for i, p := range players {
		renderedPlayers[i] = PlayerNameInPlayers(p, playerList)
	}
	return renderedPlayers
}

func Padded(s interface{}, width int) string {
	text := String(s)
	buf := bytes.NewBufferString(text)
	extra := width - StrLen(text)
	if extra > 0 {
		buf.WriteString(strings.Repeat(" ", extra))
	}
	return buf.String()
}

type spanCellInTable struct {
	start, span, contentWidth int
}

func Table(cells [][]interface{}, rowPadding, colPadding int) string {
	// First calculate widths
	widths := map[int]int{}
	spans := []spanCellInTable{}
	for _, row := range cells {
		for colIndex, cell := range row {
			w := StrWidth(String(cell))
			switch t := cell.(type) {
			case Unbounded:
				continue
			case CellSpan:
				spans = append(spans, spanCellInTable{
					colIndex,
					t.Cols,
					StrWidth(String(t.Content)),
				})
				continue
			}
			if w > widths[colIndex] {
				widths[colIndex] = w
			}
		}
	}
	for _, span := range spans {
		remaining := span.contentWidth
		for i := span.start; i < span.start+span.span; i++ {
			remaining -= widths[i]
			if i > span.start {
				remaining -= colPadding
			}
		}
		if remaining > widths[span.start+span.span-1] {
			widths[span.start+span.span-1] = remaining
		}
	}
	// Output cells
	buf := bytes.NewBuffer([]byte{})
	for rowIndex, row := range cells {
		if rowIndex > 0 {
			buf.WriteString(strings.Repeat("\n", rowPadding+1))
		}
		rowHeight := 0
		cellLines := map[int][]string{}
		for colIndex, cellRaw := range row {
			if span, ok := cellRaw.(CellSpan); ok {
				cellRaw = span.Content
			}
			cellLines[colIndex] = strings.Split(String(cellRaw), "\n")
			if cl := len(cellLines[colIndex]); cl > rowHeight {
				rowHeight = cl
			}
		}

		for lIndex := 0; lIndex < rowHeight; lIndex++ {
			if lIndex != 0 {
				buf.WriteString("\n")
			}
			for colIndex, cellRaw := range row {
				lineContent := ""
				if len(cellLines[colIndex]) > lIndex {
					lineContent = cellLines[colIndex][lIndex]
				}
				width := widths[colIndex]
				if span, ok := cellRaw.(CellSpan); ok {
					cellRaw = span.Content
					width = 0
					for i := colIndex; i < colIndex+span.Cols; i++ {
						width += widths[i]
						if i > colIndex {
							width += colPadding
						}
					}
				}
				var content string
				switch cellRaw.(type) {
				case Centred:
					content = Centre(lineContent, width)
				case RightAligned:
					content = Right(lineContent, width)
				default:
					content = String(lineContent)
				}
				if colIndex != len(row)-1 {
					content = Padded(content, width+colPadding)
				}
				buf.WriteString(content)
			}
		}
	}
	return buf.String()
}

func CommaList(list []string) string {
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0]
	}
	if len(list) == 2 {
		return list[0] + " and " + list[1]
	}
	return list[0] + ", " + CommaList(list[1:])
}

func StrLen(s string) int {
	return utf8.RuneCountInString(RenderPlain(s))
}

func StrWidth(s string) int {
	w := 0
	for _, l := range strings.Split(s, "\n") {
		if lw := StrLen(l); lw > w {
			w = lw
		}
	}
	return w
}

func Centre(s interface{}, width int) string {
	buf := bytes.NewBuffer([]byte{})
	str := String(s)
	extra := width - StrLen(str)
	left := ""
	right := ""
	if extra > 0 {
		left = strings.Repeat(" ", (extra+1)/2)
		right = strings.Repeat(" ", extra-len(left))
	}
	buf.WriteString(left)
	buf.WriteString(str)
	buf.WriteString(right)
	return buf.String()
}

func Right(s interface{}, width int) string {
	buf := bytes.NewBuffer([]byte{})
	str := String(s)
	extra := width - StrLen(str)
	if extra > 0 {
		buf.WriteString(strings.Repeat(" ", extra))
	}
	buf.WriteString(str)
	return buf.String()
}

func String(s interface{}) string {
	return fmt.Sprintf("%v", s)
}

func Bold(s interface{}) string {
	lines := strings.Split(String(s), "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("{{b}}%s{{_b}}", lines[i])
	}
	return strings.Join(lines, "\n")
}

func Colour(s interface{}, colour string) string {
	if !IsValidColour(colour) {
		log.Fatalf("%s is not a valid colour", colour)
	}
	lines := strings.Split(String(s), "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf(`{{c "%s"}}%s{{_c}}`, colour, lines[i])
	}
	return strings.Join(lines, "\n")
}

func Markup(s interface{}, colour string, bold bool) string {
	str := String(s)
	if colour != "" {
		str = Colour(str, colour)
	}
	if bold {
		str = Bold(str)
	}
	return str
}

func StringsToInterfaces(strs []string) []interface{} {
	ints := make([]interface{}, len(strs))
	for i, s := range strs {
		ints[i] = s
	}
	return ints
}
