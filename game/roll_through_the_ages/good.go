package roll_through_the_ages

const (
	GoodWood = iota
	GoodStone
	GoodPottery
	GoodCloth
	GoodSpearhead
)

var Goods = []int{
	GoodWood,
	GoodStone,
	GoodPottery,
	GoodCloth,
	GoodSpearhead,
}

var GoodStrings = map[int]string{
	GoodWood:      "wood",
	GoodStone:     "stone",
	GoodPottery:   "pottery",
	GoodCloth:     "cloth",
	GoodSpearhead: "spearhead",
}

var GoodColours = map[int]string{
	GoodWood:      "magenta",
	GoodStone:     "gray",
	GoodPottery:   "red",
	GoodCloth:     "blue",
	GoodSpearhead: "yellow",
}

func GoodsReversed() []int {
	l := len(Goods)
	rev := make([]int, l)
	for i, _ := range Goods {
		rev[i] = l - i - 1
	}
	return rev
}

func GoodMaximum(good int) int {
	return 8 - good
}

func GoodValue(good, n int) int {
	return (n * (n + 1) / 2) * (good + 1)
}
