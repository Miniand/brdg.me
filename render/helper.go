package render

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

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

func Padded(text string, width int) string {
	buf := bytes.NewBufferString(text)
	extra := width - StrLen(text)
	if extra > 0 {
		buf.WriteString(strings.Repeat(" ", extra))
	}
	return buf.String()
}

func Table(cells [][]string, rowPadding, colPadding int) string {
	// First calculate widths
	widths := map[int]int{}
	for _, row := range cells {
		for colIndex, cell := range row {
			w := StrLen(cell)
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
		for colIndex, cell := range row {
			var padded string
			if colIndex == len(row)-1 {
				// Last col doesn't get right padding
				padded = cell
			} else {
				padded = Padded(cell, widths[colIndex]+colPadding)
			}
			buf.WriteString(padded)
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

func Centre(s string, width int) string {
	buf := bytes.NewBuffer([]byte{})
	extra := width - StrLen(s)
	left := ""
	right := ""
	if extra > 0 {
		left = strings.Repeat(" ", (extra+1)/2)
		right = strings.Repeat(" ", extra-len(left))
	}
	buf.WriteString(left)
	buf.WriteString(s)
	buf.WriteString(right)
	return buf.String()
}

func Right(s string, width int) string {
	buf := bytes.NewBuffer([]byte{})
	extra := width - StrLen(s)
	if extra > 0 {
		buf.WriteString(strings.Repeat(" ", extra))
	}
	buf.WriteString(s)
	return buf.String()
}
