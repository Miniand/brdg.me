package render

import (
	"bytes"
	"fmt"
	"github.com/beefsack/brdg.me/game"
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

func PlayerName(playerNum int, name string) string {
	return fmt.Sprintf(`{{b}}{{c "%s"}}%s{{_c}}{{_b}}`, PlayerColour(playerNum),
		name)
}

func Padded(text string, width int, g game.Playable) (string, error) {
	buf := bytes.NewBufferString(text)
	plain, err := RenderPlain(text, g)
	if err != nil {
		return "", err
	}
	extra := width - utf8.RuneCountInString(plain)
	if extra > 0 {
		buf.WriteString(strings.Repeat(" ", extra))
	}
	return buf.String(), nil
}
