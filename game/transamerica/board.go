package transamerica

import (
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	City string
}

type Loc struct {
	X, Y byte
}

func ParseLoc(input string) (l Loc, err error) {
	if len(input) != 2 {
		err = errors.New("location text must be two letters between A and Z")
		return
	}
	upper := strings.ToUpper(input)
	if upper[0] < 'A' || upper[0] > 'Z' || upper[1] < 'A' || upper[1] > 'Z' {
		err = errors.New("location text must be two letters between A and Z")
		return
	}
	l.Y = upper[0]
	l.X = upper[1]
	return
}

func MustParseLoc(input string) Loc {
	l, err := ParseLoc(input)
	if err != nil {
		panic(err)
	}
	return l
}

func (l Loc) String() string {
	return fmt.Sprintf("%c%c", l.Y, l.X)
}

type Board struct {
	Nodes map[Loc]*Node
}
