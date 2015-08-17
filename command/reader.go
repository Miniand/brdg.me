package command

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Reader struct {
	*bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}

func NewReaderString(s string) *Reader {
	return NewReader(strings.NewReader(s))
}

func (p *Reader) ReadWhile(fn func(rune) bool) (string, error) {
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

func (p *Reader) ReadWord() (string, error) {
	return p.ReadWhile(func(r rune) bool {
		return !unicode.IsSpace(r)
	})
}

func (p *Reader) ReadSpace() (string, error) {
	return p.ReadWhile(func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func (p *Reader) ReadToEndOfLine() (string, error) {
	return p.ReadWhile(func(r rune) bool {
		return r != '\r' && r != '\n'
	})
}

func (p *Reader) ReadLineArgs() ([]string, error) {
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
