package roll_through_the_ages

const (
	DiceFood = iota
	DiceGood
	DiceSkull
	DiceWorkers
	DiceFoodOrWorkers
	DiceCoins
)

var DiceFaces = []int{
	DiceFood,
	DiceGood,
	DiceSkull,
	DiceWorkers,
	DiceFoodOrWorkers,
	DiceCoins,
}

var DiceStrings = map[int]string{
	DiceFood:          "FFF",
	DiceGood:          "G",
	DiceSkull:         "GXG",
	DiceWorkers:       "WWW",
	DiceFoodOrWorkers: "FF/WW",
	DiceCoins:         "C",
}

var DiceValueColours = map[string]string{
	"F": "green",
	"G": "magenta",
	"X": "red",
	"W": "cyan",
	"C": "yellow",
}

func Roll() int {
	return r.Int() % len(DiceFaces)
}

func RollN(n int) []int {
	dice := make([]int, n)
	for i := 0; i < n; i++ {
		dice[i] = Roll()
	}
	return dice
}
