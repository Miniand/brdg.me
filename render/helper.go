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

func PlayerName(playerNum int, player string) string {
	return fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`, PlayerColour(playerNum),
		player)
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

func Padded(text string, width int) (string, error) {
	buf := bytes.NewBufferString(text)
	plain, err := RenderPlain(text)
	if err != nil {
		return "", err
	}
	extra := width - utf8.RuneCountInString(plain)
	if extra > 0 {
		buf.WriteString(strings.Repeat(" ", extra))
	}
	return buf.String(), nil
}

func Table(cells [][]string, rowPadding, colPadding int) (string, error) {
	// First calculate widths
	var err error
	widths := map[int]int{}
	for _, row := range cells {
		for colIndex, cell := range row {
			plain, err := RenderPlain(cell)
			if err != nil {
				return "", err
			}
			w := utf8.RuneCountInString(plain)
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
				padded, err = Padded(cell, widths[colIndex]+colPadding)
				if err != nil {
					return "", err
				}
			}
			buf.WriteString(padded)
		}
	}
	return buf.String(), nil
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
