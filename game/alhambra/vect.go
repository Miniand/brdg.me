package alhambra

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

func (v Vect) Inverse() Vect {
	return Vect{-v.X, -v.Y}
}

func (v Vect) Add(other Vect) Vect {
	return Vect{v.X + other.X, v.Y + other.Y}
}

func (v Vect) Sub(other Vect) Vect {
	return v.Add(other.Inverse())
}
