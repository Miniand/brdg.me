package agricola_2p

const (
	Up = 1 << iota
	Right
	Down
	Left
)

type Loc struct {
	X, Y int
}

func (l Loc) Move(dir int) Loc {
	newX, newY := Move(l.X, l.Y, dir)
	return Loc{
		X: newX,
		Y: newY,
	}
}

var OppDir = map[int]int{
	Up:    Down,
	Down:  Up,
	Left:  Right,
	Right: Left,
}

func Move(x, y, dir int) (newX, newY int) {
	newX, newY = x, y
	if dir&Up != 0 {
		newY--
	}
	if dir&Down != 0 {
		newY++
	}
	if dir&Left != 0 {
		newX--
	}
	if dir&Right != 0 {
		newX++
	}
	return
}

type Tile struct {
	Borders  int
	Building Building
	Trough   bool
}

type Tiles map[Loc]*Tile

func (t Tiles) At(l Loc) *Tile {
	if t[l] == nil {
		return &Tile{}
	}
	return t[l]
}

func (t Tiles) Neighbour(l Loc, dir int) *Tile {
	return t.At(l.Move(dir))
}

func (t Tiles) Border(l Loc, dir int) bool {
	til := t.At(l)
	nei := t.Neighbour(l, dir)
	if til.Building != nil || nei.Building != nil {
		return true
	}
	switch dir {
	case Up, Left:
		return til.Borders&dir != 0
	case Down, Right:
		return nei.Borders&OppDir[dir] != 0
	}
	return false
}
