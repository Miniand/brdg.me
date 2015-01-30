package alhambra

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
