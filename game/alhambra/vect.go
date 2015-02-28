package alhambra

const (
	DirUp = 1 << iota
	DirRight
	DirDown
	DirLeft
)

var Dirs = []int{
	DirUp,
	DirRight,
	DirDown,
	DirLeft,
}

var DirInverse = map[int]int{
	DirUp:    DirDown,
	DirRight: DirLeft,
	DirDown:  DirUp,
	DirLeft:  DirRight,
}

func RotDir(dir, n int) int {
	if n == 0 {
		return dir
	}
	for n < 0 {
		n += 4
	}
	next := DirUp
	switch dir {
	case DirUp:
		next = DirRight
	case DirRight:
		next = DirDown
	case DirDown:
		next = DirLeft
	}
	return RotDir(next, n-1)
}

type Vect struct {
	X, Y int
}

var VectUp = Vect{0, -1}
var VectDown = Vect{0, 1}
var VectLeft = Vect{-1, 0}
var VectRight = Vect{1, 0}
var VectUpLeft = VectUp.Add(VectLeft)
var VectUpRight = VectUp.Add(VectRight)
var VectDownLeft = VectDown.Add(VectLeft)
var VectDownRight = VectDown.Add(VectRight)

var VectDirsOrth = []Vect{
	VectUp,
	VectRight,
	VectDown,
	VectLeft,
}

var VectDirsAll = []Vect{
	VectUp,
	VectUpRight,
	VectRight,
	VectDownRight,
	VectDown,
	VectDownLeft,
	VectLeft,
	VectUpLeft,
}

func (v Vect) Inverse() Vect {
	return Vect{-v.X, -v.Y}
}

func (v Vect) Add(other Vect) Vect {
	return Vect{v.X + other.X, v.Y + other.Y}
}

func (v Vect) Sub(other Vect) Vect {
	return v.Add(other.Inverse())
}

func (v Vect) RotAll(n int) Vect {
	l := len(VectDirsAll)
	for i, uv := range VectDirsAll {
		if v == uv {
			ni := i + n
			for ni < 0 {
				ni += l
			}
			return VectDirsAll[ni%l]
		}
	}
	panic("Can only call on unit vector")
}

var DirVectMap = map[int]Vect{
	DirUp:    VectUp,
	DirDown:  VectDown,
	DirLeft:  VectLeft,
	DirRight: VectRight,
}
