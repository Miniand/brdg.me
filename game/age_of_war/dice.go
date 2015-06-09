package age_of_war

const (
	Dice1Infantry = iota
	Dice2Infantry
	Dice3Infantry
	DiceArchery
	DiceCavalry
	DiceDaimyo
)

var DiceInfantry = map[int]int{
	Dice1Infantry: 1,
	Dice2Infantry: 2,
	Dice3Infantry: 3,
}

func Roll() int {
	return rnd.Int() % 6
}

func RollN(n int) []int {
	if n <= 0 {
		return []int{}
	}
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = Roll()
	}
	return ints
}
