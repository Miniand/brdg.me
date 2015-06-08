package cathedral

import "fmt"

type Loc struct {
	X, Y int
}

func (l Loc) Add(other Loc) Loc {
	return Loc{l.X + other.X, l.Y + other.Y}
}

func (l Loc) Neg() Loc {
	return Loc{-l.X, -l.Y}
}

func (l Loc) Sub(other Loc) Loc {
	return l.Add(other.Neg())
}

func (l Loc) Neighbour(dir int) Loc {
	return l.Add(UnitLoc(dir))
}

func (l Loc) String() string {
	return fmt.Sprintf("%c%d", 'A'+l.Y, l.X+1)
}

func (l Loc) Rotate(n int) Loc {
	switch {
	case n > 0:
		return (Loc{-l.Y, l.X}).Rotate(n - 1)
	case n < 0:
		return (Loc{l.Y, -l.X}).Rotate(n + 1)
	default:
		return l
	}
}

func (l Loc) Valid() bool {
	return l.X >= 0 && l.X <= 9 && l.Y >= 0 && l.Y <= 9
}

type Locs []Loc

func (ls Locs) Rotate(n int) Locs {
	nls := make(Locs, len(ls))
	for i, l := range ls {
		nls[i] = l.Rotate(n)
	}
	return nls
}

func UnitLoc(dir int) Loc {
	l := Loc{}
	if dir&DirUp == DirUp {
		l.Y--
	}
	if dir&DirRight == DirRight {
		l.X++
	}
	if dir&DirDown == DirDown {
		l.Y++
	}
	if dir&DirLeft == DirLeft {
		l.X--
	}
	return l
}
