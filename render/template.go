package render

import (
	"strings"
	"unicode"
)

const BaseColour = Black

func WalkTemplate(tmpl string, cb func(text, colour string, bold, large bool)) {
	colour := BaseColour
	colourStack := []string{}
	bold := 0
	large := 0

	lex := &lexer{
		input: tmpl,
	}

	for {
		upto := lex.readUntil("{{")
		if len(upto) > 0 {
			cb(upto, colour, bold > 0, large > 0)
		}
		if lex.eof() {
			break
		}
		lex.readN(2)
		contents := lex.readUntil("}}")
		lex.readN(2)
		parts := strings.SplitN(contents, " ", 2)
		action := parts[0]
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}
		switch action {
		case "c":
			if l := len(args); l > 1 {
				colour = args[1 : l-1]
				colourStack = append(colourStack, colour)
			}
		case "_c":
			l := len(colourStack)
			if l > 1 {
				colour = colourStack[l-2]
			} else {
				colour = BaseColour
			}
			if l > 0 {
				colourStack = colourStack[:l-1]
			}
		case "b":
			bold++
		case "_b":
			bold--
		case "l":
			large++
		case "_l":
			large--
		}
	}
}

type lexer struct {
	input string
	pos   int
}

func (l *lexer) readUntil(until string) string {
	index := strings.Index(l.input[l.pos:], until)
	if index == -1 {
		return l.readToEnd()
	}
	return l.readN(index)
}

func (l *lexer) readWord() string {
	for i, b := range []byte(l.input[l.pos:]) {
		if unicode.IsSpace(rune(b)) {
			return l.readN(i)
		}
	}
	return l.readToEnd()
}

func (l *lexer) eof() bool {
	return l.pos >= len(l.input)-1
}

func (l *lexer) readN(n int) string {
	if inLen := len(l.input); l.pos+n > inLen {
		n = inLen - l.pos
	}
	newPos := l.pos + n
	str := l.input[l.pos:newPos]
	l.pos = newPos
	return str
}

func (l *lexer) readToEnd() string {
	return l.readN(len(l.input) - l.pos)
}
