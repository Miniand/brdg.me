package grid

type Loc struct {
	X, Y int
}

type Gridder interface {
	SetTile(l Loc, tile interface{})
	Tile(l Loc) interface{}
	Find(tile interface{}) (at Loc, ok bool)
	Locs() chan Loc
	Each(loc Loc, tile interface{})
	Neighbours(l Loc) []Loc
	Neighbour(l Loc, dir int) Loc
	Bounds() (lower, upper Loc)
}

type Colourer interface {
	Colour() string
}

type ColourPrioritiser interface {
	ColourPriority() int
}

type Messager interface {
	Message() string
}

func (l Loc) Add(l2 Loc) Loc {
	return Loc{l.X + l2.X, l.Y + l2.Y}
}
