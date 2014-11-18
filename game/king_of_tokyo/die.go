package king_of_tokyo

const (
	Die1 = iota
	Die2
	Die3
	DieEnergy
	DieAttack
	DieHeal
)

var Dice = []int{
	Die1,
	Die2,
	Die3,
	DieEnergy,
	DieAttack,
	DieHeal,
}

func RollDie() int {
	return r.Int() % 6
}

func RollDice(n int) []int {
	dice := make([]int, n)
	for i := 0; i < n; i++ {
		dice[i] = RollDie()
	}
	return dice
}
