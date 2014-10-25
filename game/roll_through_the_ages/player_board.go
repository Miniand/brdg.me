package roll_through_the_ages

const (
	BaseCitySize = 3
)

var CityLevels = []int{3, 7, 12, 18}

var MaxCityProgress = CityLevels[len(CityLevels)-1]

type PlayerBoard struct {
	CityProgress       int
	Developments       map[int]bool
	Monuments          map[int]int
	MonumentBuiltFirst map[int]bool
	Food               int
	Resources          map[int]int
	Disasters          int
}

func NewPlayerBoard() *PlayerBoard {
	developments := map[int]bool{}
	for _, d := range Developments {
		if r.Int()%2 == 0 {
			developments[d] = true
		}
	}
	return &PlayerBoard{
		Developments:       developments, //map[int]bool{},
		Monuments:          map[int]int{},
		MonumentBuiltFirst: map[int]bool{},
		Food:               3,
		Resources:          map[int]int{},
	}
}

func (b *PlayerBoard) Cities() int {
	size := BaseCitySize
	for _, l := range CityLevels {
		if b.CityProgress < l {
			break
		}
		size += 1
	}
	return size
}

func (b *PlayerBoard) Score() int {
	score := 0
	// Developments
	for d, ok := range b.Developments {
		if !ok {
			continue
		}
		score += DevelopmentValues[d].Points
	}
	// Monuments
	builtMonuments := 0 // Track how many built for bonus score calculation
	for m, num := range b.Monuments {
		mv := MonumentValues[m]
		if num >= mv.Size {
			builtMonuments += 1
			if b.MonumentBuiltFirst[m] {
				score += mv.Points
			} else {
				score += mv.SubsequentPoints()
			}
		}
	}
	// Bonus points
	if b.Developments[DevelopmentArchitecture] {
		score += builtMonuments
	}
	if b.Developments[DevelopmentEmpire] {
		score += b.Cities()
	}
	return score - b.Disasters
}

func (b *PlayerBoard) CoinsDieValue() int {
	if b.Developments[DevelopmentCoinage] {
		return 12
	}
	return 7
}

func (b *PlayerBoard) FoodModifier() int {
	if b.Developments[DevelopmentAgriculture] {
		return 1
	}
	return 0
}

func (b *PlayerBoard) WorkerModifier() int {
	if b.Developments[DevelopmentMasonry] {
		return 1
	}
	return 0
}
