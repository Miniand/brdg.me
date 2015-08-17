package command

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Parser struct {
	*bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{bufio.NewReader(r)}
}

func NewParserString(s string) *Parser {
	return NewParser(strings.NewReader(s))
}

func (p *Parser) ReadWhile(fn func(rune) bool) (string, error) {
	var (
		r   rune
		err error
		out = &bytes.Buffer{}
	)
	for {
		r, _, err = p.ReadRune()
		if err != nil {
			if err == io.EOF && out.Len() > 0 {
				err = nil
			}
			break
		}
		if fn(r) {
			out.WriteRune(r)
		} else {
			p.UnreadRune()
			break
		}
	}
	return out.String(), err
}

func (p *Parser) ReadWord() (string, error) {
	return p.ReadWhile(func(r rune) bool {
		return !unicode.IsSpace(r)
	})
}

func (p *Parser) ReadSpace() (string, error) {
	return p.ReadWhile(func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func (p *Parser) ReadToEndOfLine() (string, error) {
	return p.ReadWhile(func(r rune) bool {
		return r != '\r' && r != '\n'
	})
}

func (p *Parser) ReadLineArgs() ([]string, error) {
	line, err := p.ReadToEndOfLine()
	if l := len(line); l > 0 {
		if err == io.EOF {
			err = nil
		}
	}
	if err != nil {
		return nil, err
	}
	return strings.Fields(line), nil
}
