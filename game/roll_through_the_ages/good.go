package roll_through_the_ages

const (
	GoodWood = iota
	GoodStone
	GoodPottery
	GoodCloth
	GoodSpearhead
)

func GoodMaximum(good int) int {
	return 8 - good
}

func GoodValue(good, n int) int {
	return (n * (n + 1) / 2) * (good + 1)
}
