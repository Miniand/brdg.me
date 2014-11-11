package roll_through_the_ages

const (
	MonumentStepPyramid = iota
	MonumentStoneCircle
	MonumentTemple
	MonumentObelisk
	MonumentHangingGardens
	MonumentGreatWall
	MonumentGreatPyramid
)

type Monument struct {
	Name             string
	Size             int
	Points           int
	SubsequentPoints int
	Effect           string
}

var Monuments = []int{
	MonumentStepPyramid,
	MonumentStoneCircle,
	MonumentTemple,
	MonumentObelisk,
	MonumentHangingGardens,
	MonumentGreatWall,
	MonumentGreatPyramid,
}

var MonumentValues = map[int]Monument{
	MonumentStepPyramid: {
		Name:   "Step Pyramid",
		Size:   3,
		Points: 1,
	},
	MonumentStoneCircle: {
		Name:             "Stone Circle",
		Size:             5,
		Points:           2,
		SubsequentPoints: 1,
	},
	MonumentTemple: {
		Name:             "Temple",
		Size:             7,
		Points:           4,
		SubsequentPoints: 3,
	},
	MonumentObelisk: {
		Name:             "Obelisk",
		Size:             9,
		Points:           6,
		SubsequentPoints: 4,
	},
	MonumentHangingGardens: {
		Name:             "Hanging Gardens",
		Size:             11,
		Points:           8,
		SubsequentPoints: 5,
	},
	MonumentGreatWall: {
		Name:             "Wall", // Changed from Great Wall to help string matching.
		Size:             13,
		Points:           10,
		SubsequentPoints: 6,
		Effect:           "invasion has no effect",
	},
	MonumentGreatPyramid: {
		Name:             "Great Pyramid",
		Size:             15,
		Points:           12,
		SubsequentPoints: 8,
	},
}
