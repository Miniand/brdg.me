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
	Goods              map[int]int
	Disasters          int
	Ships              int
}

func NewPlayerBoard() *PlayerBoard {
	return &PlayerBoard{
		Developments:       map[int]bool{},
		Monuments:          map[int]int{},
		MonumentBuiltFirst: map[int]bool{},
		Food:               3,
		Goods:              map[int]int{},
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
				score += mv.SubsequentPoints
			}
		}
	}
	// Bonus points
	if b.Developments[DevelopmentCommerce] {
		score += b.GoodsNum()
	}
	if b.Developments[DevelopmentArchitecture] {
		score += builtMonuments * 2
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

func (b *PlayerBoard) GainGoods(n int) {
	good := GoodWood
	for i := 0; i < n; i++ {
		b.GainGood(good)
		good = (good + 1) % len(Goods)
	}
}

func (b *PlayerBoard) GainGood(good int) {
	max := GoodMaximum(good)
	if b.Goods[good] < max {
		b.Goods[good] += 1
	}
	// Extra stone if player has quarry
	if good == GoodStone && b.Developments[DevelopmentQuarrying] &&
		b.Goods[good] < max {
		b.Goods[good] += 1
	}
}

func (b *PlayerBoard) GoodsNum() int {
	num := 0
	for _, n := range b.Goods {
		num += n
	}
	return num
}

func (b *PlayerBoard) GoodsValue() int {
	val := 0
	for g, n := range b.Goods {
		val += GoodValue(g, n)
	}
	return val
}

func (b *PlayerBoard) HasBuilt(monument int) bool {
	return b.Monuments[monument] >= MonumentValues[monument].Size
}
