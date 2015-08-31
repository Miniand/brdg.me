package transamerica

import (
	"errors"
	"fmt"
	"strings"
)

const (
	TerrainMountain = iota
	TerrainRiver
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

func (l Loc) IsLess(other Loc) bool {
	return l.X < other.X || l.X == other.X && l.Y < other.Y
}

type Edge struct {
	L1, L2 Loc
}

func NewEdge(l1, l2 Loc) Edge {
	e := Edge{l1, l2}
	// Make sure the edge locs are ordered so it doesn't miss in maps
	if e.L2.IsLess(e.L1) {
		e.L1, e.L2 = e.L2, e.L1
	}
	return e
}

func ParseEdge(input string) (e Edge, err error) {
	parts := strings.Fields(input)
	if len(parts) != 2 {
		err = errors.New("edge text must be two locations next to each other")
		return
	}
	l1, err := ParseLoc(parts[0])
	if err != nil {
		return
	}
	l2, err := ParseLoc(parts[1])
	if err != nil {
		return
	}
	e = NewEdge(l1, l2)
	return
}

func MustParseEdge(input string) Edge {
	e, err := ParseEdge(input)
	if err != nil {
		panic(err)
	}
	return e
}

type Board struct {
	Nodes   map[Loc]*Node
	Terrain map[Edge]int
}
