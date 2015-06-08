package alhambra

func NewTile(tileType, cost int, walls ...int) Tile {
	t := Tile{tileType, cost, map[int]bool{}}
	for _, w := range walls {
		t.Walls[w] = true
	}
	return t
}

var Tiles = []Tile{
	NewTile(TileTypePavillion, 6, DirUp),
	NewTile(TileTypePavillion, 4, DirRight, DirDown),
	NewTile(TileTypePavillion, 8),
	NewTile(TileTypePavillion, 3, DirDown, DirLeft),
	NewTile(TileTypePavillion, 5, DirUp, DirLeft),
	NewTile(TileTypePavillion, 2, DirUp, DirRight, DirLeft),
	NewTile(TileTypePavillion, 7, DirRight),

	NewTile(TileTypeSeraglio, 5, DirDown, DirLeft),
	NewTile(TileTypeSeraglio, 6, DirRight, DirDown),
	NewTile(TileTypeSeraglio, 7, DirLeft),
	NewTile(TileTypeSeraglio, 8, DirDown),
	NewTile(TileTypeSeraglio, 3, DirRight, DirDown, DirLeft),
	NewTile(TileTypeSeraglio, 4, DirUp, DirRight),
	NewTile(TileTypeSeraglio, 9),

	NewTile(TileTypeArcades, 9),
	NewTile(TileTypeArcades, 4, DirUp, DirRight, DirDown),
	NewTile(TileTypeArcades, 10),
	NewTile(TileTypeArcades, 8, DirUp),
	NewTile(TileTypeArcades, 7, DirRight, DirDown),
	NewTile(TileTypeArcades, 6, DirDown, DirLeft),
	NewTile(TileTypeArcades, 6, DirUp, DirRight),
	NewTile(TileTypeArcades, 5, DirUp, DirLeft),
	NewTile(TileTypeArcades, 8, DirRight),

	NewTile(TileTypeChambers, 6, DirRight, DirDown),
	NewTile(TileTypeChambers, 7, DirDown, DirLeft),
	NewTile(TileTypeChambers, 10),
	NewTile(TileTypeChambers, 5, DirUp, DirDown, DirLeft),
	NewTile(TileTypeChambers, 7, DirUp, DirRight),
	NewTile(TileTypeChambers, 9, DirLeft),
	NewTile(TileTypeChambers, 8, DirUp, DirLeft),
	NewTile(TileTypeChambers, 11),
	NewTile(TileTypeChambers, 9, DirDown),

	NewTile(TileTypeGarden, 9, DirRight),
	NewTile(TileTypeGarden, 7, DirUp, DirDown, DirLeft),
	NewTile(TileTypeGarden, 8, DirUp, DirRight),
	NewTile(TileTypeGarden, 11),
	NewTile(TileTypeGarden, 8, DirDown, DirLeft),
	NewTile(TileTypeGarden, 10, DirUp),
	NewTile(TileTypeGarden, 6, DirRight, DirDown, DirLeft),
	NewTile(TileTypeGarden, 8, DirUp, DirLeft),
	NewTile(TileTypeGarden, 10),
	NewTile(TileTypeGarden, 12, DirDown),
	NewTile(TileTypeGarden, 10, DirLeft),

	NewTile(TileTypeTower, 8, DirUp, DirRight, DirDown),
	NewTile(TileTypeTower, 9, DirUp, DirLeft),
	NewTile(TileTypeTower, 13, DirRight),
	NewTile(TileTypeTower, 9, DirRight, DirDown),
	NewTile(TileTypeTower, 7, DirUp, DirRight, DirLeft),
	NewTile(TileTypeTower, 11, DirUp),
	NewTile(TileTypeTower, 9, DirUp, DirRight),
	NewTile(TileTypeTower, 11, DirDown),
	NewTile(TileTypeTower, 12),
	NewTile(TileTypeTower, 11),
	NewTile(TileTypeTower, 10, DirLeft),
}
