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
	Name      string
	Size      int
	Points    int
	Blacklist int
	Effect    string
}

func (m Monument) SubsequentPoints() int {
	return m.Points / 2
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

func MonumentsForPlayerCount(count int) []int {
	monuments := []int{}
	for _, m := range Monuments {
		mv := MonumentValues[m]
		if mv.Blacklist != count {
			monuments = append(monuments, m)
		}
	}
	return monuments
}

func (g *Game) Monuments() []int {
	return MonumentsForPlayerCount(len(g.Players))
}

var MonumentValues = map[int]Monument{
	MonumentStepPyramid: {
		Name:   "Step Pyramid",
		Size:   3,
		Points: 1,
	},
	MonumentStoneCircle: {
		Name:   "Stone Circle",
		Size:   5,
		Points: 2,
	},
	MonumentTemple: {
		Name:      "Temple",
		Size:      7,
		Points:    4,
		Blacklist: 2,
	},
	MonumentObelisk: {
		Name:   "Obelisk",
		Size:   9,
		Points: 6,
	},
	MonumentHangingGardens: {
		Name:      "Hanging Gardens",
		Size:      11,
		Points:    8,
		Blacklist: 3,
	},
	MonumentGreatWall: {
		Name:   "Great Wall",
		Size:   13,
		Points: 10,
		Effect: "invasion has no effect",
	},
	MonumentGreatPyramid: {
		Name:      "Great Pyramid",
		Size:      15,
		Points:    12,
		Blacklist: 2,
	},
}
