package alhambra

const (
	GridInvalidNoFountain = "there is no fountain"
	GridInvalidWall       = "adjoining tile sides must match, either both walls or both not walls"
	GridInvalidCannotWalk = "you must be able to walk from the fountain to all other tiles"
	GridInvalidGap        = "you are not allowed to create empty gaps"
)

type Grid map[Vect]Tile

func (g Grid) TileAt(v Vect) Tile {
	t := g[v]
	if t.Walls == nil {
		t.Walls = map[int]bool{}
	}
	return t
}

func (g Grid) Clone() Grid {
	ng := Grid{}
	for v, t := range g {
		ng[v] = t
	}
	return ng
}

func (g Grid) FountainLoc() (Vect, bool) {
	for v, t := range g {
		if t.Type == TileTypeFountain {
			return v, true
		}
	}
	return Vect{}, false
}

func (g Grid) Bounds() (min, max Vect) {
	first := true
	for v := range g {
		if first || v.X < min.X {
			min.X = v.X
		}
		if first || v.Y < min.Y {
			min.Y = v.Y
		}
		if first || v.X > max.X {
			max.X = v.X
		}
		if first || v.Y > max.Y {
			max.Y = v.Y
		}
		first = false
	}
	return
}

func (g Grid) IsValid() (bool, string) {
	fv, ok := g.FountainLoc()
	if !ok {
		return false, GridInvalidNoFountain
	}

	// Walk from the fountain to find all connected tiles.
	var next Vect
	walkStack := []Vect{fv}
	inWalkStack := map[Vect]bool{}
	connected := map[Vect]bool{}
	for len(walkStack) > 0 {
		next, walkStack = walkStack[0], walkStack[1:]
		nextTile := g.TileAt(next)
		inWalkStack[next] = false
		connected[next] = true
		for _, dir := range Dirs {
			dv := next.Add(DirVectMap[dir])
			dvTile := g.TileAt(dv)

			if dvTile.Type == TileTypeEmpty {
				continue
			}

			if nextTile.Walls[dir] {
				if !dvTile.Walls[DirInverse[dir]] {
					return false, GridInvalidWall
				}
				continue
			}

			if inWalkStack[dv] || connected[dv] {
				continue
			}

			walkStack = append(walkStack, dv)
			inWalkStack[dv] = true
		}
	}

	// Iterate over all tiles to make sure they are connected to fountain.
	for v, t := range g {
		if t.Type != TileTypeEmpty && !connected[v] {
			return false, GridInvalidCannotWalk
		}
	}

	// Walk all external space and make sure all space in bounds is subset.
	min, max := g.Bounds()
	fv = min.Add(VectUpLeft)
	walkStack = []Vect{fv}
	inWalkStack = map[Vect]bool{}
	connected = map[Vect]bool{}
	for len(walkStack) > 0 {
		next, walkStack = walkStack[0], walkStack[1:]
		inWalkStack[next] = false
		connected[next] = true
		for _, dir := range Dirs {
			dv := next.Add(DirVectMap[dir])
			dvTile := g.TileAt(dv)

			if dvTile.Type != TileTypeEmpty ||
				inWalkStack[dv] || connected[dv] ||
				dv.X < min.X-1 || dv.X > max.X+1 ||
				dv.Y < min.Y-1 || dv.Y > max.Y+1 {
				continue
			}

			walkStack = append(walkStack, dv)
			inWalkStack[dv] = true
		}
	}

	// Iterate over internal empty tiles to make sure they aren't in gaps.
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			v := Vect{x, y}
			if g.TileAt(v).Type == TileTypeEmpty && !connected[v] {
				return false, GridInvalidGap
			}
		}
	}

	return true, ""
}
