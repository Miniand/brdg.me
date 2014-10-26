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

func Table(cells [][]interface{}, rowPadding, colPadding int) string {
	// First calculate widths
	widths := map[int]int{}
	for _, row := range cells {
		for colIndex, cell := range row {
			switch cell.(type) {
			case Unbounded:
				continue
			}
			w := StrLen(String(cell))
			if w > widths[colIndex] {
				widths[colIndex] = w
			}
		}
	}
	// Output cells
	buf := bytes.NewBuffer([]byte{})
	for rowIndex, row := range cells {
		if rowIndex > 0 {
			buf.WriteString(strings.Repeat("\n", rowPadding+1))
		}
		for colIndex, cellRaw := range row {
			var content string
			switch cellRaw.(type) {
			case Centred:
				content = Centre(cellRaw, widths[colIndex])
			case RightAligned:
				content = Right(cellRaw, widths[colIndex])
			default:
				content = String(cellRaw)
			}
			if colIndex != len(row)-1 {
				content = Padded(content, widths[colIndex]+colPadding)
			}
			buf.WriteString(content)
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
	return fmt.Sprintf("{{b}}%v{{_b}}", s)
}

func Colour(s interface{}, colour string) string {
	if !IsValidColour(colour) {
		log.Fatalf("%s is not a valid colour", colour)
	}
	return fmt.Sprintf(`{{c "%s"}}%v{{_c}}`, colour, s)
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
